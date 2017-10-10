package systemd

import (
	"errors"
	"fmt"

	"github.com/coreos/go-systemd/daemon"
)

// NotifyReady sends systemd notify READY=1
func NotifyReady() (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, "READY=1"); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyReloading sends systemd notify RELOADING=1
func NotifyReloading() (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, "RELOADING=1"); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyStopping sends systemd notify STOPPING=1
func NotifyStopping() (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, "STOPPING=1"); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyStatus sends systemd notify STATUS=%s{status}
func NotifyStatus(status string) (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, fmt.Sprintf("STATUS=%s", status)); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyErrNo sends systemd notify ERRNO=%d{errno}
func NotifyErrNo(errno int) (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, fmt.Sprintf("ERRNO=%d", errno)); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyBusError sends systemd notify BUSERROR=%s{buserror}
func NotifyBusError(buserror string) (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, fmt.Sprintf("BUSERROR=%s", buserror)); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyMainPID sends systemd notify MAINPID=%d{mainpid}
func NotifyMainPID(mainpid int) (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, fmt.Sprintf("MAINPID=%d", mainpid)); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyWatchDog sends systemd notify WATCHDOG=1
func NotifyWatchDog() (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, "WATCHDOG=1"); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}

// NotifyWatchDogUSec sends systemd notify WATCHDOG_USEC=%d{usec}
func NotifyWatchDogUSec(usec int64) (err error) {
	var sent bool
	if sent, err = daemon.SdNotify(false, fmt.Sprintf("WATCHDOG_USEC=%d", usec)); !sent && err == nil {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return
}
