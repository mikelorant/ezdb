package app

import (
	"context"
	"fmt"
	"net"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

func getDialFunc(tun *Tunnel) (func(ctx context.Context, address string) (net.Conn, error), error) {
	fn := dialFunc(&net.Dialer{})

	if isTunnel(tun) {
		tunnel, err := makeTunnel(tun)
		if err != nil {
			return nil, fmt.Errorf("unable to make tunnel: %w", err)
		}
		fn = dialFunc(tunnel)
	}

	return fn, nil
}

func dialFunc(dialer Dialer) func(ctx context.Context, address string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return dialer.Dial("tcp", address)
	}
}
