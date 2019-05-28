package activation

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	// https://www.freedesktop.org/software/systemd/man/sd_listen_fds.html
	listenFdsStart = 3
)

var listeners []net.Listener

func init() {
	var err error
	if listeners, err = activationListeners(); err != nil {
		log.Println("warning: systemd socket activation disabled")
	}
}

// Listen returns the net.Listener matching the given address.
func Listen(addr string) (net.Listener, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s addr: %v", addr, err)
	}
	for _, l := range listeners {
		_, p, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s addr: %v", l.Addr().String(), err)
		}
		if port == p {
			return l, nil
		}
	}
	return nil, fmt.Errorf("%s addr not found", addr)
}

func activationListeners() ([]net.Listener, error) {
	files, err := getFiles()
	if err != nil {
		return nil, err
	}
	listeners := make([]net.Listener, len(files))
	for i, f := range files {
		listener, err := net.FileListener(f)
		if err != nil {
			return nil, fmt.Errorf("failed to init new file listener: %v", err)
		}
		if err = f.Close(); err != nil {
			return nil, fmt.Errorf("failed to close %s file: %v", f.Name(), err)
		}
		listeners[i] = listener
	}
	return listeners, nil
}

func getFiles() ([]*os.File, error) {
	err := checkEnv()
	if err != nil {
		return nil, err
	}
	listenPid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse LISTEN_PID as int: %v", err)
	}
	if listenPid != os.Getpid() {
		return nil, fmt.Errorf("bad pid (LISTEN_PID=%d, pid=%d)", listenPid, os.Getpid())
	}
	listenFds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse LISTEN_FDS as int: %v", err)
	}
	var listenFdNames []string
	if _, ok := os.LookupEnv("LISTEN_FDS"); ok {
		listenFdNames = strings.Split(os.Getenv("LISTEN_FDNAMES"), ":")
		if err = os.Unsetenv("LISTEN_FDNAMES"); err != nil {
			return nil, fmt.Errorf("failed to unset env LISTEN_FDNAMES: %v", err)
		}
	}
	if err = os.Unsetenv("LISTEN_FDS"); err != nil {
		return nil, fmt.Errorf("failed to unset env LISTEN_PID: %v", err)
	}
	if err = os.Unsetenv("LISTEN_FDNAMES"); err != nil {
		return nil, fmt.Errorf("failed to unset env LISTEN_PID: %v", err)
	}
	files := make([]*os.File, 0, listenFds)
	for fd := listenFdsStart; fd < listenFdsStart+listenFds; fd++ {
		// http://man7.org/linux/man-pages/man2/fcntl.2.html
		syscall.CloseOnExec(fd)
		name := "LISTEN_FD_" + strconv.Itoa(fd)
		if listenFdNames != nil {
			offset := fd - listenFdsStart
			if offset < len(listenFdNames) && len(listenFdNames[offset]) > 0 {
				name = listenFdNames[offset]
			}
		}
		files = append(files, os.NewFile(uintptr(fd), name))
	}
	return files, nil
}

func checkEnv() error {
	_, ok := os.LookupEnv("LISTEN_PID")
	if !ok {
		return errors.New("LISTEN_PID not set")
	}
	if _, ok = os.LookupEnv("LISTEN_FDS"); !ok {
		return errors.New("LISTEN_FDS not set")
	}
	return nil
}
