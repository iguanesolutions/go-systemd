package systemd

import "fmt"

// IsNotifyEnabled tells if systemd notify is enabled or not.
func IsNotifyEnabled() bool {
	return socket != nil
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
