package systemd

import (
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

// WatchDog is an interface to the systemd watchdog mechanism
type WatchDog struct {
	watchdogLimit time.Duration
	sent          bool
	err           error
}

// NewWatchdog returns :
// *Controller, nil : watchdog is enabled and controller ready to be used
// nil, nil         : watchdog is not enabled
// nil, error       : an error occured during initialization
func NewWatchdog() (*WatchDog, error) {
	c := WatchDog{}
	c.watchdogLimit, c.err = daemon.SdWatchdogEnabled(false)
	if c.err != nil {
		return nil, fmt.Errorf("initialization went wrong: %v", c.err)
	}
	if c.watchdogLimit == 0 {
		return nil, nil
	}
	return &c, nil
}

// SendHeartbeat sends a keepalive notification to systemd watchdog
func (c *WatchDog) SendHeartbeat() error {
	c.sent, c.err = daemon.SdNotify(false, "WATCHDOG=1")
	if c.err != nil {
		return fmt.Errorf("can't send hearbeat: %v", c.err)
	}
	if !c.sent {
		return errors.New("can't send hearbeat: notifications are not supported (NOTIFY_SOCKET is unset)")
	}
	return nil
}

// NewTicker initializes and returns a ticker set at 1/3 of the watchdog duration.
// It can be used by clients to trigger checks before using SendHeartbeat().
func (c *WatchDog) NewTicker() *time.Ticker {
	return time.NewTicker(c.watchdogLimit / 3)
}
