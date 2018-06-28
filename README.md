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
notifier, err := systemd.NewNotifier()
if err != nil {
    log.Fatalf("can't start systemd notifier: %v", err)
}
// Init http server
server := &http.Server{
    Addr:    "host:port",
    Handler: myHTTPHandler,
}

// Do some more inits

// Notify ready to systemd
if err = sysd.NotifyReady(); err != nil {
   log.Fatalf("can't notify ready to systemd: %v", err)
}
// Start the server
if err = server.ListenAndServe(); err != nil {
    log.Fatalf("can't start http server: %v", err)
}
```

When stopping, you can notify systemd that you are stopping before shutting down your http server
and stopping your controllers

```go
import (
    "github.com/iguanesolutions/go-systemd"
)

var err error
if err = sysnotifier.NotifyStopping(); err != nil {
    log.Fatalf("can't notify stopping to systemd: %v", err)
}

// Stop some more things

// Stop the server (with timeout)
ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelCtx()
if err = server.Shutdown(ctx); err != nil {
    log.Fatalf("can't shutdown http server: %v", err)
}
```

You can also notify status to systemd

```go
import (
    "github.com/iguanesolutions/go-systemd"
)

if err := sysnotifier.NotifyStatus(fmt.Sprintf("There is currently %d active connections", activeConns)); err != nil {
    log.Fatalf("can't notify status to systemd: %v", err)
}
```

## Watchdog

`todo`
