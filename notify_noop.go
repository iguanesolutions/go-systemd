// +build !linux

package systemd

import (
	"net"
)

var socket *net.UnixAddr

func send(state string) error {
	// fallback implementation for OS other than linux
	return nil
}
