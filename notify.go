package systemd

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// Notifier wraps and facilitate systemd notify communications
type Notifier struct {
	socket *net.UnixAddr
}

// NewNotifier returns an initialized and ready to use Notifier if systemd notify is supported
func NewNotifier() (notifier *Notifier, err error) {
	// Validate Systemd notify is supported
	socketAddr := net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}
	if socketAddr.Name == "" {
		err = errors.New("notifications are not supported (NOTIFY_SOCKET is unset)")
		return
	}
	// All good return the object
	notifier = &Notifier{
		socket: &socketAddr,
	}
	return
}

// NotifyReady sends systemd notify READY=1
func (s *Notifier) NotifyReady() error {
	return s.send("READY=1")
}

// NotifyReloading sends systemd notify RELOADING=1
func (s *Notifier) NotifyReloading() error {
	return s.send("RELOADING=1")
}

// NotifyStopping sends systemd notify STOPPING=1
func (s *Notifier) NotifyStopping() error {
	return s.send("STOPPING=1")
}

// NotifyStatus sends systemd notify STATUS=%s{status}
func (s *Notifier) NotifyStatus(status string) error {
	return s.send(fmt.Sprintf("STATUS=%s", status))
}

// NotifyErrNo sends systemd notify ERRNO=%d{errno}
func (s *Notifier) NotifyErrNo(errno int) error {
	return s.send(fmt.Sprintf("ERRNO=%d", errno))
}

// NotifyBusError sends systemd notify BUSERROR=%s{buserror}
func (s *Notifier) NotifyBusError(buserror string) error {
	return s.send(fmt.Sprintf("BUSERROR=%s", buserror))
}

// NotifyMainPID sends systemd notify MAINPID=%d{mainpid}
func (s *Notifier) NotifyMainPID(mainpid int) error {
	return s.send(fmt.Sprintf("MAINPID=%d", mainpid))
}

// NotifyWatchDog sends systemd notify WATCHDOG=1
func (s *Notifier) NotifyWatchDog() error {
	return s.send("WATCHDOG=1")
}

// NotifyWatchDogUSec sends systemd notify WATCHDOG_USEC=%d{usec}
func (s *Notifier) NotifyWatchDogUSec(usec int64) (err error) {
	return s.send(fmt.Sprintf("WATCHDOG_USEC=%d", usec))
}

func (s *Notifier) send(state string) (err error) {
	// Try to open socket
	conn, err := net.DialUnix(s.socket.Net, nil, s.socket)
	if err != nil {
		err = fmt.Errorf("can't open unix socket: %v", err)
		return
	}
	defer conn.Close()
	// Write data into it
	if _, err = conn.Write([]byte(state)); err != nil {
		err = fmt.Errorf("can't write into the unix socket: %v", err)
	}
	return
}
