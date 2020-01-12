package main

import (
	"SSD-Go/ss"
	"errors"
)

// cancel is not nil if and only if an instance is running
var cancel func() error

func runShadowsocks(flags ss.Flags) (err error) {
	if cancel != nil {
		err = errors.New("an instance is already running")
		return
	}
	// cancel will be nil if err is not nil
	cancel, err = ss.Run(flags, ss.Config{
		Verbose:    true,
		UDPTimeout: ss.DefaultConfig.UDPTimeout,
	})
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
