package sysdwatchdog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	sysdnotify "github.com/iguanesolutions/go-systemd/v5/notify"
)

// WatchDog is an interface to the systemd watchdog mechanism
type WatchDog struct {
	interval time.Duration
	checks   time.Duration
}

// New returns an initialized and ready to use WatchDog
func New() (wd *WatchDog, err error) {
	// Check WatchDog is supported and usable
	interval, err := getWatchDogInterval()
	if err != nil {
		return
	}
	// Return the initialized controller
	wd = &WatchDog{
		interval: interval,
		checks:   interval / 2,
	}
	return
}

// based on https://github.com/coreos/go-systemd/blob/master/daemon/watchdog.go
func getWatchDogInterval() (interval time.Duration, err error) {
	// WATCHDOG_USEC
	wusec := os.Getenv("WATCHDOG_USEC")
	if wusec == "" {
		err = errors.New("watchdog does not seem to be enabled: WATCHDOG_USEC env unset")
		return
	}
	wusecTyped, err := strconv.Atoi(wusec)
	if err != nil {
		err = fmt.Errorf("can't convert WATCHDOG_USEC as int: %s", err)
		return
	}
	if wusecTyped <= 0 {
		err = fmt.Errorf("WATCHDOG_USEC must be a positive number")
		return
	}
	interval = time.Duration(wusecTyped) * time.Microsecond
	// WATCHDOG_PID
	wpid := os.Getenv("WATCHDOG_PID")
	if wpid == "" {
		return // No WATCHDOG_PID: can't check if we are the one, let's go with it
	}
	wpidTyped, err := strconv.Atoi(wpid)
	if err != nil {
		err = fmt.Errorf("can't convert WATCHDOG_PID as int: %s", err)
		return
	}
	if os.Getpid() != wpidTyped {
		err = fmt.Errorf("WATCHDOG_PID is %d and we are %d: we are not the watched PID", wpidTyped, os.Getpid())
	}
	return
}

// SendHeartbeat sends a keepalive notification to systemd watchdog
func (c *WatchDog) SendHeartbeat() error {
	if !sysdnotify.IsEnabled() {
		return errors.New("failed to notify watchdog: systemd notify is diabled")
	}
	return sysdnotify.WatchDog()
}

// GetChecksDuration returns the ideal time for a client to perform (active or passive collect) checks.
// Is is equal at 1/3 of watchdogInterval
func (c *WatchDog) GetChecksDuration() time.Duration {
	return c.checks
}

// GetLimitDuration returns the systemd watchdog limit provided by systemd
func (c *WatchDog) GetLimitDuration() time.Duration {
	return c.interval
}

// NewTicker initializes and returns a ticker set at watchdogChecks (which is set at 1/3 of watchdogInterval).
// It can be used by clients to trigger checks before using SendHeartbeat().
func (c *WatchDog) NewTicker() *time.Ticker {
	return time.NewTicker(c.checks)
}
