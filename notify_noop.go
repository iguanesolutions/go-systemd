// +build !linux

package systemd

import (
	"net"
)

var socket *net.UnixAddr

// NotifyRaw send state thru the notify socket if any.
// If the notify socket was not detected, it is a noop call.
// Use IsNotifyEnabled() to determine if the notify socket has been detected
func NotifyRaw(state string) error {
	// fallback implementation for OSes other than linux
	return nil
}
