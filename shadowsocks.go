package main

import (
	"errors"
	"github.com/perqin/go-shadowsocks2"
)

// cancel is not nil if and only if an instance is running
var cancel func() error

func runShadowsocks(flags shadowsocks2.Flags) (err error) {
	if cancel != nil {
		err = errors.New("an instance is already running")
		return
	}
	cancelFunc, err := shadowsocks2.Run(flags)
	if err != nil {
		return
	}
	cancel = func() error {
		cancelFunc()
		return nil
	}
	return
}

func stopShadowsocks() error {
	if cancel == nil {
		return errors.New("no instance is running")
	}
	cancelFunc := cancel
	cancel = nil
	return cancelFunc()
}
