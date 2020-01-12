package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/shadowsocks/go-shadowsocks2/core"
)

//var config struct {
//	Verbose    bool
//	UDPTimeout time.Duration
//}

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func logf(f string, v ...interface{}) {
	//if config.Verbose {
	logger.Output(2, fmt.Sprintf(f, v...))
	//}
}

type shadowsocksConfig struct {
	Client    string
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

var cancel func() error

func RunShadowsocks(flags shadowsocksConfig) (err error) {
	if cancel != nil {
		err = errors.New("an instance is already running")
		return
	}
	// cancel will be nil if err is not nil
	cancel, err = someRunSs(flags)
	return
}

func StopShadowsocks() error {
	if cancel == nil {
		return errors.New("no instance is running")
	}
	return cancel()
}

// someRunSs start a shadowsocks instance with given flags.
// TODO: Should be in shadowsocks package
func someRunSs(flags shadowsocksConfig) (cancel func() error, err error) {
	runningGroup := run.Group{}
	// For user interruption
	//runningGroup.Add(func() error {
	//	return <-interruptChan
	//}, func(err error) {
	//	select {
	//	case interruptChan<-err:
	//	default:
	//	}
	//})
	addr := flags.Client
	cipher := flags.Cipher
	password := flags.Password
	if strings.HasPrefix(addr, "ss://") {
		addr, cipher, password, err = parseURL(addr)
		if err != nil {
			return
		}
	}
	ciph, err := core.PickCipher(cipher, nil, password)
	if err != nil {
		return
	}
	ctx := context.Background()
	ctx, cncl := context.WithCancel(ctx)
	if flags.Socks != "" {
		runningGroup.Add(func() error {
			socksLocal(flags.Socks, addr, ciph.StreamConn)
			return nil
		}, func(err error) {
			// TODO
		})
	}
	runningErr := make(chan error)
	go func() {
		err := runningGroup.Run()
		log.Printf("Stop running deal to %v\n", err)
		runningErr <- err
	}()
	cancel = func() error {
		cncl()
		return <-runningErr
	}
	return
}

func runShadowsocks(flags shadowsocksConfig) (err error) {
	//flag.BoolVar(&config.Verbose, "verbose", false, "verbose mode")
	//flag.StringVar(&flags.Cipher, "cipher", "AEAD_CHACHA20_POLY1305", "available ciphers: "+strings.Join(core.ListCipher(), " "))
	//flag.StringVar(&flags.Key, "key", "", "base64url-encoded key (derive from password if empty)")
	//flag.IntVar(&flags.Keygen, "keygen", 0, "generate a base64url-encoded random key of given length in byte")
	//flag.StringVar(&flags.Password, "password", "", "password")
	//flag.StringVar(&flags.Client, "c", "", "client connect address or url")
	//flag.StringVar(&flags.Socks, "socks", "", "(client-only) SOCKS listen address")
	//flag.BoolVar(&flags.UDPSocks, "u", false, "(client-only) Enable UDP support for SOCKS")
	//flag.StringVar(&flags.RedirTCP, "redir", "", "(client-only) redirect TCP from this address")
	//flag.StringVar(&flags.RedirTCP6, "redir6", "", "(client-only) redirect TCP IPv6 from this address")
	//flag.StringVar(&flags.TCPTun, "tcptun", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	//flag.StringVar(&flags.UDPTun, "udptun", "", "(client-only) UDP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	//flag.DurationVar(&config.UDPTimeout, "udptimeout", 5*time.Minute, "UDP tunnel timeout")
	//flag.Parse()

	addr := flags.Client
	cipher := flags.Cipher
	password := flags.Password
	//var err error

	if strings.HasPrefix(addr, "ss://") {
		addr, cipher, password, err = parseURL(addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	ciph, err := core.PickCipher(cipher, nil, password)
	if err != nil {
		log.Fatal(err)
	}

	//if flags.UDPTun != "" {
	//	for _, tun := range strings.Split(flags.UDPTun, ",") {
	//		p := strings.Split(tun, "=")
	//		go udpLocal(p[0], addr, p[1], ciph.PacketConn)
	//	}
	//}

	//if flags.TCPTun != "" {
	//	for _, tun := range strings.Split(flags.TCPTun, ",") {
	//		p := strings.Split(tun, "=")
	//		go tcpTun(p[0], addr, p[1], ciph.StreamConn)
	//	}
	//}

	if flags.Socks != "" {
		//socks.UDPEnabled = flags.UDPSocks
		go socksLocal(flags.Socks, addr, ciph.StreamConn)
		//if flags.UDPSocks {
		//	go udpSocksLocal(flags.Socks, addr, ciph.PacketConn)
		//}
	}

	//if flags.RedirTCP != "" {
	//	go redirLocal(flags.RedirTCP, addr, ciph.StreamConn)
	//}
	//
	//if flags.RedirTCP6 != "" {
	//	go redir6Local(flags.RedirTCP6, addr, ciph.StreamConn)
	//}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	return nil
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

////////
// Copied from go-shadowsocks2/tcp.go

// Create a SOCKS server listening on addr and proxy to server.
func socksLocal(addr, server string, shadow func(net.Conn) net.Conn) {
	logf("SOCKS proxy %s <-> %s", addr, server)
	tcpLocal(addr, server, shadow, func(c net.Conn) (socks.Addr, error) { return socks.Handshake(c) })
}

// Listen on addr and proxy to server to reach target from getAddr.
func tcpLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(net.Conn) (socks.Addr, error)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logf("failed to listen on %s: %v", addr, err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			logf("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()
			c.(*net.TCPConn).SetKeepAlive(true)
			tgt, err := getAddr(c)
			if err != nil {

				// UDP: keep the connection until disconnect then free the UDP socket
				if err == socks.InfoUDPAssociate {
					buf := []byte{}
					// block here
					for {
						_, err := c.Read(buf)
						if err, ok := err.(net.Error); ok && err.Timeout() {
							continue
						}
						logf("UDP Associate End.")
						return
					}
				}

				logf("failed to get target address: %v", err)
				return
			}

			rc, err := net.Dial("tcp", server)
			if err != nil {
				logf("failed to connect to server %v: %v", server, err)
				return
			}
			defer rc.Close()
			rc.(*net.TCPConn).SetKeepAlive(true)
			rc = shadow(rc)

			if _, err = rc.Write(tgt); err != nil {
				logf("failed to send target address: %v", err)
				return
			}

			logf("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
			_, _, err = relay(rc, c)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return // ignore i/o timeout
				}
				logf("relay error: %v", err)
			}
		}()
	}
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func relay(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}
