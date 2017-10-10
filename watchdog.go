package systemd

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

// WatchDog is an interface to the systemd watchdog mechanism
type WatchDog struct {
	watchdogLimit  time.Duration
	watchdogChecks time.Duration
	sent           bool
	err            error
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
	c.watchdogChecks = c.watchdogLimit / 3
	return &c, nil
}

// SendHeartbeat sends a keepalive notification to systemd watchdog
func (c *WatchDog) SendHeartbeat() error {
	return NotifyWatchDog()
}

// GetChecksDuration returns the ideal time for a client to perform (active or passive collect) checks.
// Is is equal at 1/3 of watchdogLimit
func (c *WatchDog) GetChecksDuration() time.Duration {
	return c.watchdogChecks
}

// GetLimitDuration returns the systemd watchdog limit provided by systemd
func (c *WatchDog) GetLimitDuration() time.Duration {
	return c.watchdogLimit
}

// NewTicker initializes and returns a ticker set at watchdogChecks (which is set at 1/3 of watchdogLimit).
// It can be used by clients to trigger checks before using SendHeartbeat().
func (c *WatchDog) NewTicker() *time.Ticker {
	return time.NewTicker(c.watchdogChecks)
}
