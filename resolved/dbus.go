package resolved

import (
	"context"
	"errors"
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

var (
	// errNotSupported is returned when a method is not implemeted.
	errNotSupported = errors.New("method not supported yet")
)

// ResolveHostname, ResolveAddress, ResolveRecord, ResolveService
// The four methods above accept and return a 64-bit flags value.
// In most cases passing 0 is sufficient and recommended.
// However, the following flags are defined to alter the look-up
const (
	SD_RESOLVED_DNS           = uint64(1) << 0
	SD_RESOLVED_LLMNR_IPV4    = uint64(1) << 1
	SD_RESOLVED_LLMNR_IPV6    = uint64(1) << 2
	SD_RESOLVED_MDNS_IPV4     = uint64(1) << 3
	SD_RESOLVED_MDNS_IPV6     = uint64(1) << 4
	SD_RESOLVED_NO_CNAME      = uint64(1) << 5
	SD_RESOLVED_NO_TXT        = uint64(1) << 6
	SD_RESOLVED_NO_ADDRESS    = uint64(1) << 7
	SD_RESOLVED_NO_SEARCH     = uint64(1) << 8
	SD_RESOLVED_AUTHENTICATED = uint64(1) << 9
)

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

// ResolveHostname takes a hostname and resolves it to one or more IP addresses.
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

// ResolveAddress takes a DNS resource record (RR) type, class and name
// and retrieves the full resource record set (RRset), including the RDATA, for it.
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

// TODO
// ResolveRecord
func (c *Conn) ResolveRecord(ctx context.Context, ifindex int, name string, class dnsmessage.Class, rtype dnsmessage.Type, flags uint64) (records []ResourceRecord, outflags uint64, err error) {
	//err = c.Call(ctx, "ResolveRecord", ifindex, name, class, rtype, flags).Store(&records, &outflags)
	err = errNotSupported
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

// TODO
// ResolveService
func (c *Conn) ResolveService(ctx context.Context, ifindex int, name string, stype string, domain string, family int,
	flags uint64) (srvData []interface{}, txtData []interface{}, canonicalName string, canonicalType string, canonicalDomain string, outflags uint64, err error) {
	//err = c.Call(ctx, "ResolveService", ifindex, name, stype, domain, family, flags).Store(&srvData, &txtData, &canonicalName, &canonicalType, &canonicalDomain, &outflags)
	err = errNotSupported
	return
}

// GetLink takes a network interface index and returns the object path
// to the org.freedesktop.resolve1.Link object corresponding to it.
// ctx: Context to use
// ifindex: The network interface index to get link for
func (c *Conn) GetLink(ctx context.Context, ifindex int) (path string, err error) {
	err = c.Call(ctx, "GetLink", ifindex).Store(&path)
	return
}

// LinkDNS represents a DNS server address to use in SetLinkDNS method.
type LinkDNS struct {
	Family  int    // can be either syscall.AF_INET or syscall.AF_INET6
	Address net.IP // binary address
}

// SetLinkDNS sets the DNS servers to use on a specific interface.
// ctx: Context to use
// ifindex: The network interface index to set
// addrs: array of DNS server IP address records.
func (c *Conn) SetLinkDNS(ctx context.Context, ifindex int, addrs []LinkDNS) (err error) {
	err = c.Call(ctx, "SetLinkDNS", ifindex, addrs).Store()
	return
}

type LinkDNSEx struct {
	Family  int    // can be either syscall.AF_INET or syscall.AF_INET6
	Address net.IP // binary address
	Port    uint16 // the port number
	Name    string // the DNS Name
}

// SetLinkDNSEx is similar to SetLinkDNS(), but allows an IP port
// (instead of the default 53) and DNS name to be specified for each DNS server.
// The server name is used for Server Name Indication (SNI), which is useful when DNS-over-TLS is used.
// ctx: Context to use
// ifindex: The network interface index
// addrs: array of DNS server IP address records.
func (c *Conn) SetLinkDNSEx(ctx context.Context, ifindex int, addrs []LinkDNSEx) error {
	return c.Call(ctx, "SetLinkDNSEx", ifindex, addrs).Store()
}

type LinkDomain struct {
	Domain        string // the domain name
	RoutingDomain bool   // whether the specified domain shall be used as a search domain (false), or just as a routing domain (true).
}

// SetLinkDomains sets the search and routing domains to use on a specific network interface for DNS look-ups.
// ctx: Context to use
// ifindex: The network interface index
// domains: array of domains
func (c *Conn) SetLinkDomains(ctx context.Context, ifindex int, domains []LinkDomain) error {
	return c.Call(ctx, "SetLinkDomains", ifindex, domains).Store()
}

// SetLinkDefaultRoute specifies whether the link shall be used as the default route for name queries
// ctx: Context to use
// ifindex: The network interface index
// enable: enable/disable link as default route.
func (c *Conn) SetLinkDefaultRoute(ctx context.Context, ifindex int, enable bool) error {
	return c.Call(ctx, "SetLinkDefaultRoute", ifindex, enable).Store()
}

// SetLinkLLMNR enables or disables LLMNR support on a specific network interface.
// ctx: Context to use
// ifindex: The network interface index
// mode: either empty or one of "yes", "no" or "resolve".
func (c *Conn) SetLinkLLMNR(ctx context.Context, ifindex int, mode string) error {
	return c.Call(ctx, "SetLinkLLMNR", ifindex, mode).Store()
}

// SetLinkMulticastDNS enables or disables MulticastDNS support on a specific interface.
// ctx: Context to use
// ifindex: The network interface index
// mode: either empty or one of "yes", "no" or "resolve".
func (c *Conn) SetLinkMulticastDNS(ctx context.Context, ifindex int, mode string) error {
	return c.Call(ctx, "SetLinkMulticastDNS", ifindex, mode).Store()
}

// SetLinkDNSOverTLS enables or disables enables or disables DNS-over-TLS on a specific network interface.
// ctx: Context to use
// ifindex: The network interface index
// mode: either empty or one of "yes", "no", or "opportunistic"
func (c *Conn) SetLinkDNSOverTLS(ctx context.Context, ifindex int, mode string) error {
	return c.Call(ctx, "SetLinkDNSOverTLS", ifindex, mode).Store()
}

// SetLinkDNSSEC enables or disables DNSSEC validation on a specific network interface.
// ctx: Context to use
// ifindex: The network interface index
// mode: either empty or one of "yes", "no", or "allow-downgrade"
func (c *Conn) SetLinkDNSSEC(ctx context.Context, ifindex int, mode string) error {
	return c.Call(ctx, "SetLinkDNSSEC", ifindex, mode).Store()
}

// SetLinkDNSSECNegativeTrustAnchors configures DNSSEC Negative Trust Anchors (NTAs) for a specific network interface.
// ctx: Context to use
// ifindex: The network interface index
// names: array of domains
func (c *Conn) SetLinkDNSSECNegativeTrustAnchors(ctx context.Context, ifindex int, names []string) error {
	return c.Call(ctx, "SetLinkDNSSECNegativeTrustAnchors", ifindex, names).Store()
}

// RevertLink reverts all per-link settings to the defaults on a specific network interface.
// ctx: Context to use
// ifindex: The network interface index.
func (c *Conn) RevertLink(ctx context.Context, ifindex int) error {
	return c.Call(ctx, "RevertLink", ifindex).Store()
}

// TODO
// RegisterService
func (c *Conn) RegisterService(ctx context.Context, name string, nameTemplate string, stype string,
	svcPort uint16, svcPriority uint16, svcWeight uint16, txtData []interface{}) (svcPath string, err error) {
	//err = c.Call(ctx, "RegisterService", name, nameTemplate, stype, svcPort, svcPriority, svcWeight, txtData).Store(&svcPath)
	err = errNotSupported
	return
}

// TODO
// UnregisterService
func (c *Conn) UnregisterService(ctx context.Context, svcPath string) error {
	//return c.Call(ctx, "UnregisterService", svcPath).Store()
	return errNotSupported
}

// ResetStatistics resets the various statistics counters that systemd-resolved maintains to zero.
func (c *Conn) ResetStatistics(ctx context.Context) error {
	return c.Call(ctx, "ResetStatistics").Store()
}

// FlushCaches
func (c *Conn) FlushCaches(ctx context.Context) error {
	return c.Call(ctx, "FlushCaches").Store()
}

// ResetServerFeatures
func (c *Conn) ResetServerFeatures(ctx context.Context) error {
	return c.Call(ctx, "ResetServerFeatures").Store()
}

type Link struct {
	obj dbus.BusObject
}

func NewLink(c *Conn, path string) Link {
	return Link{
		obj: c.conn.Object(dbusDest, dbus.ObjectPath(path)),
	}
}

// TODO
// 	SetDNS(in  a(iay) addresses);
// 	SetDNSEx(in  a(iayqs) addresses);
// 	SetDomains(in  a(sb) domains);
// 	SetDefaultRoute(in  b enable);
// 	SetLLMNR(in  s mode);
// 	SetMulticastDNS(in  s mode);
// 	SetDNSOverTLS(in  s mode);
// 	SetDNSSEC(in  s mode);
// 	SetDNSSECNegativeTrustAnchors(in  as names);
// 	Revert();
