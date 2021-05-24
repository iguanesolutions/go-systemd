package resolved

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/godbus/dbus/v5"
	"golang.org/x/net/dns/dnsmessage"
)

const (
	dbusDest      = "org.freedesktop.resolve1"
	dbusInterface = "org.freedesktop.resolve1.Manager"
	dbusPath      = "/org/freedesktop/resolve1"
)

// Conn represents a systemd-resolved dbus connection.
type Conn struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

// NewConn returns a new and ready to use dbus connection.
// You must close that connection when you have been done with it.
func NewConn() (*Conn, error) {
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return nil, fmt.Errorf("failed to init private conn to system bus: %v", err)
	}
	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}
	err = conn.Auth(methods)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to auth with external method: %v", err)
	}
	err = conn.Hello()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to make hello call: %v", err)
	}
	return &Conn{
		conn: conn,
		obj:  conn.Object(dbusDest, dbus.ObjectPath(dbusPath)),
	}, nil
}

// Call wraps obj.CallWithContext by using 0 as flags and format the method with the dbus manager interface.
func (c *Conn) Call(ctx context.Context, method string, args ...interface{}) *dbus.Call {
	return c.obj.CallWithContext(ctx, fmt.Sprintf("%s.%s", dbusInterface, method), 0, args...)
}

// Close closes the current dbus connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Address represents an address returned by ResolveHostname.
type Address struct {
	IfIndex int    // network interface index
	Family  int    // can be either syscall.AF_INET or syscall.AF_INET6
	Address net.IP // binary address
}

func (a Address) String() string {
	return fmt.Sprintf(`{
	IfIndex: %d,
	Family:  %d,
	IP:      %s,
}`, a.IfIndex, a.Family, a.Address.String())
}

// ResolveHostname
// ctx: Context to use
// ifindex: Network interface index where to look (0 means any)
// name: Hostname
// family: Which address family to look for (syscall.AF_UNSPEC, syscall.AF_INET, syscall.AF_INET6)
// flags: Input flags parameter
func (c *Conn) ResolveHostname(ctx context.Context, ifindex int, name string, family int, flags uint64) (addresses []Address, canonical string, outflags uint64, err error) {
	err = c.Call(ctx, "ResolveHostname", ifindex, name, family, flags).Store(&addresses, &canonical, &outflags)
	return
}

// Name represents a hostname returned by ResolveAddress.
type Name struct {
	IfIndex  int    // network interface index
	Hostname string // hostname
}

func (n Name) String() string {
	return fmt.Sprintf(`{
	IfIndex: %d,
	Name:    %s,
}`, n.IfIndex, n.Hostname)
}

// ResolveAddress
// ctx: Context to use
// ifindex: Network interface index where to look (0 means any)
// family: Address family (syscall.AF_INET, syscall.AF_INET6)
// address: the binary address (4 or 16 bytes)
// flags: Input flags parameter
func (c *Conn) ResolveAddress(ctx context.Context, ifindex int, family int, address net.IP, flags uint64) (names []Name, outflags uint64, err error) {
	err = c.Call(ctx, "ResolveAddress", ifindex, family, address, flags).Store(&names, &outflags)
	return
}

type ResourceRecord struct {
	IfIndex int // network interface index
	Type    dnsmessage.Type
	Class   dnsmessage.Class

	// TODO: Parse raw RR data
	// The raw RR data starts with the RR's domain name, in the original casing, followed by the RR type, class,
	// TTL and RDATA, in the binary format documented in RFC 1035. For RRs that support name compression in the payload
	// (such as MX or PTR), the compression is expanded in the returned data.
	Data []byte
}

func (c *Conn) ResolveRecord(ctx context.Context, ifindex int, name string, class dnsmessage.Class, rtype dnsmessage.Type, flags uint64) (records []ResourceRecord, outflags uint64, err error) {
	err = c.Call(ctx, "ResolveRecord", ifindex, name, class, rtype, flags).Store(&records, &outflags)
	return
}

func (r ResourceRecord) String() string {
	return fmt.Sprintf(`{
	IfIndex: %d,
	Type:    %s,
	Class:   %s,
	RData:   %s,
}`, r.IfIndex, r.Type.String(), r.Class.String(), string(r.Data))
}
