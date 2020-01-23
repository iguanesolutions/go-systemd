package systemd

import (
	"fmt"
	"net"
	"os"
)

var socket *net.UnixAddr

func init() {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket == "" {
		return
	}
	socket = &net.UnixAddr{
		Name: notifySocket,
		Net:  "unixgram",
	}
}

// IsNotifyEnabled tells if systemd notify is enabled or not.
func IsNotifyEnabled() bool {
	return socket != nil
}

// NotifyReady sends systemd notify READY=1
func NotifyReady() error {
	return send("READY=1")
}

// NotifyReloading sends systemd notify RELOADING=1
func NotifyReloading() error {
	return send("RELOADING=1")
}

// NotifyStopping sends systemd notify STOPPING=1
func NotifyStopping() error {
	return send("STOPPING=1")
}

// NotifyStatus sends systemd notify STATUS=%s{status}
func NotifyStatus(status string) error {
	return send(fmt.Sprintf("STATUS=%s", status))
}

// NotifyErrNo sends systemd notify ERRNO=%d{errno}
func NotifyErrNo(errno int) error {
	return send(fmt.Sprintf("ERRNO=%d", errno))
}

// NotifyBusError sends systemd notify BUSERROR=%s{buserror}
func NotifyBusError(buserror string) error {
	return send(fmt.Sprintf("BUSERROR=%s", buserror))
}

// NotifyMainPID sends systemd notify MAINPID=%d{mainpid}
func NotifyMainPID(mainpid int) error {
	return send(fmt.Sprintf("MAINPID=%d", mainpid))
}

// NotifyWatchDog sends systemd notify WATCHDOG=1
func NotifyWatchDog() error {
	return send("WATCHDOG=1")
}

// NotifyWatchDogUSec sends systemd notify WATCHDOG_USEC=%d{µsec}
func NotifyWatchDogUSec(usec int64) error {
	return send(fmt.Sprintf("WATCHDOG_USEC=%d", usec))
}

func send(state string) error {
	if socket != nil {
		conn, err := net.DialUnix(socket.Net, nil, socket)
		if err != nil {
			return fmt.Errorf("can't open unix socket: %v", err)
		}
		defer conn.Close()
		if _, err = conn.Write([]byte(state)); err != nil {
			return fmt.Errorf("can't write into the unix socket: %v", err)
		}
	}
	return nil
}
