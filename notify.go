// +build linux

package systemd

import (
	"fmt"
	"net"
	"os"
)

var socket *net.UnixAddr

func init() {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket != "" {
		return
	}
	socket = &net.UnixAddr{
		Name: notifySocket,
		Net:  "unixgram",
	}
}

func send(state string) error {
	if socket == nil {
		return nil
	}
	conn, err := net.DialUnix(socket.Net, nil, socket)
	if err != nil {
		return fmt.Errorf("can't open unix socket: %v", err)
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(state)); err != nil {
		return fmt.Errorf("can't write into the unix socket: %v", err)
	}
	return nil
}
