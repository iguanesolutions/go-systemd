# go-systemd

[![GoDoc](https://godoc.org/github.com/iguanesolutions/go-systemd?status.svg)](https://godoc.org/github.com/iguanesolutions/go-systemd)

Easily communicate with systemd when run as daemon within a service unit.

## Notifier

With notifier you can notify to systemd that your program is starting, stopping, reloading...

For example, if your daemon needs some time for initializing its controllers before really being considered as ready, you can specify to systemd that this is a "notify" service and send it a notification when ready.

```systemdunit
[Service]
Type=notify
```

```go
import (
    "github.com/iguanesolutions/go-systemd"
)

// Init new notifier
// If you os does not support systemd, it will return nil
// So you have to handle sysd != nil
sysd, err := systemd.NewNotifier()
if err != nil {
    log.Printf("can't start systemd notifier: %v\n", err)
}
// Init http server
server := &http.Server{
    Addr:    "host:port",
    Handler: myHTTPHandler,
}

// Do some more inits

// Notify ready to systemd
if sysd != nil {
    if err = sysd.NotifyReady(); err != nil {
        log.Printf("can't notify ready to systemd: %v\n", err)
    }
}

// Start the server
if err = server.ListenAndServe(); err != nil {
    log.Printf("can't start http server: %v\n", err)
}
```

When stopping, you can notify systemd that you have indeed received the SIGTERM and you have launched the stop procedure

```go
import (
    "github.com/iguanesolutions/go-systemd"
)

// Notify to systemd that we are stopping
var err error
if sysd != nil {
    if err = sysd.NotifyStopping(); err != nil {
        log.Printf("can't notify stopping to systemd: %v\n", err)
    }
}

// Stop some more things

// Stop the server (with timeout)
ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelCtx()
if err = server.Shutdown(ctx); err != nil {
    log.Printf("can't shutdown http server: %v\n", err)
}
```

You can also notify status to systemd

```go
import (
    "github.com/iguanesolutions/go-systemd"
)

if sysd != nil {
    if err := sysd.NotifyStatus(fmt.Sprintf("There is currently %d active connections", activeConns)); err != nil {
        log.Printf("can't notify status to systemd: %v\n", err)
    }
}
```

systemctl status output example:

```systemctlstatus
user@host:~$ systemctl status superapp.service
● superapp.service - superapp
   Loaded: loaded (/lib/systemd/system/superapp.service; enabled)
   Active: active (running) since Mon 2018-06-25 08:54:35 UTC; 3 days ago
 Main PID: 2604 (superapp)
   Status: "There is currently 1506 active connections"
   ...
```

## Watchdog

`todo`
