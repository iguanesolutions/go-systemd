package systemd

import (
	"fmt"
	"net"
)

var notifySocket *net.UnixAddr

// IsNotifyEnabled tells if systemd notify socket has been detected or not.
func IsNotifyEnabled() bool {
	return notifySocket != nil
}

// NotifyReady sends systemd notify READY=1
func NotifyReady() error {
	return NotifyRaw("READY=1")
}

// NotifyReloading sends systemd notify RELOADING=1
func NotifyReloading() error {
	return NotifyRaw("RELOADING=1")
}

// NotifyStopping sends systemd notify STOPPING=1
func NotifyStopping() error {
	return NotifyRaw("STOPPING=1")
}

// NotifyStatus sends systemd notify STATUS=%s{status}
func NotifyStatus(status string) error {
	return NotifyRaw(fmt.Sprintf("STATUS=%s", status))
}

// NotifyErrNo sends systemd notify ERRNO=%d{errno}
func NotifyErrNo(errno int) error {
	return NotifyRaw(fmt.Sprintf("ERRNO=%d", errno))
}

// NotifyBusError sends systemd notify BUSERROR=%s{buserror}
func NotifyBusError(buserror string) error {
	return NotifyRaw(fmt.Sprintf("BUSERROR=%s", buserror))
}

// NotifyMainPID sends systemd notify MAINPID=%d{mainpid}
func NotifyMainPID(mainpid int) error {
	return NotifyRaw(fmt.Sprintf("MAINPID=%d", mainpid))
}

// NotifyWatchDog sends systemd notify WATCHDOG=1
func NotifyWatchDog() error {
	return NotifyRaw("WATCHDOG=1")
}

// NotifyWatchDogUSec sends systemd notify WATCHDOG_USEC=%d{µsec}
func NotifyWatchDogUSec(usec int64) error {
	return NotifyRaw(fmt.Sprintf("WATCHDOG_USEC=%d", usec))
}

// NotifyRaw send state thru the notify socket if any.
// If the notify socket was not detected, it is a noop call.
// Use IsNotifyEnabled() to determine if the notify socket has been detected.
func NotifyRaw(state string) error {
	if notifySocket == nil {
		return nil
	}
	conn, err := net.DialUnix(notifySocket.Net, nil, notifySocket)
	if err != nil {
		return fmt.Errorf("can't open unix socket: %v", err)
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(state)); err != nil {
		return fmt.Errorf("can't write into the unix socket: %v", err)
	}
	return nil
}
