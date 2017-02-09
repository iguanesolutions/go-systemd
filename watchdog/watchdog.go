package watchdog

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

// Controller is an interface to the systemd watchdog mechanism
type Controller struct {
	watchdogLimit time.Duration
	sent          bool
	err           error
}

// New returns :
// *Controller, nil : watchdog is enabled and controller ready to be used
// nil, nil         : watchdog is not enabled
// nil, error       : an error occured during initialization
func New() (*Controller, error) {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		return nil, fmt.Errorf("initialization went wrong: %v", err)
	}
	if interval == 0 {
		return nil, nil
	}
	return &Controller{
		watchdogLimit: interval,
	}, nil
}

// SendHeartbeat sends a notification to systemd watchdog
func (c *Controller) SendHeartbeat() error {
	c.sent, c.err = daemon.SdNotify(false, "WATCHDOG=1")
	if c.err != nil {
		return fmt.Errorf("can't send hearbeat: %v", c.err)
	}
	if !c.sent {
		return fmt.Errorf("can't send hearbeat: notifications not supported (NOTIFY_SOCKET is unset)")
	}
	return nil
}

// GetTicker returns a ticker set at 1/3 of the watchdog duration.
// It can be used by clients to trigger checks before using SendHeartbeat().
func (c *Controller) GetTicker() *time.Ticker {
	return time.NewTicker(c.watchdogLimit / 3)
}
