// +build linux

package systemd

import (
	"net"
	"os"
)

func init() {
	if notifySocketName := os.Getenv("NOTIFY_SOCKET"); notifySocketName != "" {
		notifySocket = &net.UnixAddr{
			Name: notifySocketName,
			Net:  "unixgram",
		}
	}
}
