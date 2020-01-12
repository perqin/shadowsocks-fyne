package ss

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type Config struct {
	Verbose    bool
	UDPTimeout time.Duration
}

var DefaultConfig = Config{
	Verbose:    false,
	UDPTimeout: 5 * time.Minute,
}

var config = DefaultConfig

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func logf(f string, v ...interface{}) {
	if config.Verbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

type Flags struct {
	Client    string
	Server    string
	Cipher    string
	Key       string
	Password  string
	Keygen    int
	Socks     string
	RedirTCP  string
	RedirTCP6 string
	TCPTun    string
	UDPTun    string
	UDPSocks  bool
}

func main() {

	var flags Flags

	flag.BoolVar(&config.Verbose, "verbose", false, "verbose mode")
	flag.StringVar(&flags.Cipher, "cipher", "AEAD_CHACHA20_POLY1305", "available ciphers: "+strings.Join(core.ListCipher(), " "))
	flag.StringVar(&flags.Key, "key", "", "base64url-encoded key (derive from password if empty)")
	flag.IntVar(&flags.Keygen, "keygen", 0, "generate a base64url-encoded random key of given length in byte")
	flag.StringVar(&flags.Password, "password", "", "password")
	flag.StringVar(&flags.Server, "s", "", "server listen address or url")
	flag.StringVar(&flags.Client, "c", "", "client connect address or url")
	flag.StringVar(&flags.Socks, "socks", "", "(client-only) SOCKS listen address")
	flag.BoolVar(&flags.UDPSocks, "u", false, "(client-only) Enable UDP support for SOCKS")
	flag.StringVar(&flags.RedirTCP, "redir", "", "(client-only) redirect TCP from this address")
	flag.StringVar(&flags.RedirTCP6, "redir6", "", "(client-only) redirect TCP IPv6 from this address")
	flag.StringVar(&flags.TCPTun, "tcptun", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.StringVar(&flags.UDPTun, "udptun", "", "(client-only) UDP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.DurationVar(&config.UDPTimeout, "udptimeout", 5*time.Minute, "UDP tunnel timeout")
	flag.Parse()

	if flags.Keygen > 0 {
		key := make([]byte, flags.Keygen)
		io.ReadFull(rand.Reader, key)
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		return
	}

	if flags.Client == "" && flags.Server == "" {
		flag.Usage()
		return
	}

	cancel, err := Run(flags, config)
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	_ = cancel()
}

func Run(flags Flags, cfg Config) (func() error, error) {
	config = cfg

	var interrupts []chan struct{}

	var key []byte
	if flags.Key != "" {
		k, err := base64.URLEncoding.DecodeString(flags.Key)
		if err != nil {
			return nil, err
		}
		key = k
	}

	if flags.Client != "" { // client mode
		addr := flags.Client
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				return nil, err
			}
		}

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			return nil, err
		}

		if flags.UDPTun != "" {
			for _, tun := range strings.Split(flags.UDPTun, ",") {
				p := strings.Split(tun, "=")
				interruptUdpTun := make(chan struct{})
				interrupts = append(interrupts, interruptUdpTun)
				go udpLocal(interruptUdpTun, p[0], addr, p[1], ciph.PacketConn)
			}
		}

		if flags.TCPTun != "" {
			for _, tun := range strings.Split(flags.TCPTun, ",") {
				p := strings.Split(tun, "=")
				interruptTcpTun := make(chan struct{})
				interrupts = append(interrupts, interruptTcpTun)
				go tcpTun(interruptTcpTun, p[0], addr, p[1], ciph.StreamConn)
			}
		}

		if flags.Socks != "" {
			socks.UDPEnabled = flags.UDPSocks
			interruptSocksLocal := make(chan struct{})
			interrupts = append(interrupts, interruptSocksLocal)
			go socksLocal(interruptSocksLocal, flags.Socks, addr, ciph.StreamConn)
			if flags.UDPSocks {
				interruptUdpSocksLocal := make(chan struct{})
				interrupts = append(interrupts, interruptUdpSocksLocal)
				go udpSocksLocal(interruptUdpSocksLocal, flags.Socks, addr, ciph.PacketConn)
			}
		}

		if flags.RedirTCP != "" {
			interruptRedirLocal := make(chan struct{})
			interrupts = append(interrupts, interruptRedirLocal)
			go redirLocal(interruptRedirLocal, flags.RedirTCP, addr, ciph.StreamConn)
		}

		if flags.RedirTCP6 != "" {
			interruptRedir6Local := make(chan struct{})
			interrupts = append(interrupts, interruptRedir6Local)
			go redir6Local(interruptRedir6Local, flags.RedirTCP6, addr, ciph.StreamConn)
		}
	}

	// TODO: Support cancellation for server mode
	if flags.Server != "" { // server mode
		addr := flags.Server
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				return nil, err
			}
		}

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			return nil, err
		}

		go udpRemote(addr, ciph.PacketConn)
		go tcpRemote(addr, ciph.StreamConn)
	}
	logf("%d goroutines launched", len(interrupts))
	return func() error {
		for _, ch := range interrupts {
			select {
			case ch <- struct{}{}:
			default:
				logf("Fail to close some channel")
			}
		}
		// TODO: There may be some goroutine not terminated yet when reaching here, because we can know whether they
		//  have finished after receiving the data from channel. Need to refactor
		return nil
	}, nil
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
