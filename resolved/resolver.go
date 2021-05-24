package resolved

import (
	"context"
	"errors"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

// resolver is the interface to implements the same methods as the net.Resolver
type resolver interface {
	LookupAddr(ctx context.Context, addr string) (names []string, err error)
	LookupCNAME(ctx context.Context, host string) (cname string, err error)
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupIP(ctx context.Context, network, host string) ([]net.IP, error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupMX(ctx context.Context, name string) ([]*net.MX, error)
	LookupNS(ctx context.Context, name string) ([]*net.NS, error)
	LookupPort(ctx context.Context, network, service string) (port int, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) ([]string, error)
}

var (
	// ensure that types implement resolver interface
	_ resolver = &Resolver{}
	_ resolver = &net.Resolver{}
)

// Resolver represents the systemd-resolved resolver
// throught dbus connection.
type Resolver struct {
	conn   *Conn
	dialer *net.Dialer
}

type resolverOption func(r *Resolver) error

// WithConn allow you to use a custom systemd-resolved dbus connection.
func WithConn(c *Conn) resolverOption {
	return func(r *Resolver) error {
		if c == nil {
			return errors.New("conn is nil")
		}
		r.conn = c
		return nil
	}
}

// WithDialer allow you to use a custom net.Dialer.
func WithDialer(d *net.Dialer) resolverOption {
	return func(r *Resolver) error {
		if d == nil {
			return errors.New("dialer is nil")
		}
		r.dialer = d
		return nil
	}
}

// NewResolver returns a new systemd Resolver with an initialized dbus connection.
// it's up to you to close that connection when you have been done with the Resolver.
func NewResolver(opts ...resolverOption) (*Resolver, error) {
	r := &Resolver{}
	for _, opt := range opts {
		opt(r)
	}
	if r.conn == nil {
		var err error
		r.conn, err = NewConn()
		if err != nil {
			return nil, err
		}
	}
	if r.dialer == nil {
		r.dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	}
	return r, nil
}

// Close closes the current dbus connection.
// You need to close the connection when you've done with it.
func (r *Resolver) Close() error {
	return r.conn.Close()
}

// DialContext resolves address using systemd-network and use internal dialer with the resolved ip address.
// It is useful when it comes to integration with go standard library.
func (r *Resolver) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	addrs, _, _, err := r.conn.ResolveHostname(ctx, 0, host, syscall.AF_UNSPEC, 0)
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if addr.Address.To4() == nil {
			// prefer ipv6
			address = addr.Address.String()
			break
		}
		address = addr.Address.String()
	}
	return r.dialer.DialContext(ctx, network, net.JoinHostPort(address, port))
}

// HTTPClient returns a new http.Client with systemd-resolved as resolver
// and idle connections + keepalives disabled.
func (r *Resolver) HTTPClient() *http.Client {
	transport := r.pooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return &http.Client{
		Transport: transport,
	}
}

// HTTPPooledClient returns a new http.Client with systemd-resolved as resolver
// and similar default values to http.DefaultTransport.
// Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func (r *Resolver) HTTPPooledClient() *http.Client {
	return &http.Client{
		Transport: r.pooledTransport(),
	}
}

func (r *Resolver) pooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           r.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}

// LookupHost looks up the given host using the systemd-resolved resolver.
// It returns a slice of that host's addresses.
func (r *Resolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	if host == "" {
		return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	addresses, _, _, err := r.conn.ResolveHostname(ctx, 0, host, syscall.AF_UNSPEC, 0)
	if err != nil {
		return nil, err
	}
	addrs = make([]string, len(addresses))
	for i, addr := range addresses {
		addrs[i] = addr.Address.String()
	}
	return
}

// LookupAddr performs a reverse lookup for the given address, returning a list
// of names mapping to that address.
func (r *Resolver) LookupAddr(ctx context.Context, addr string) (names []string, err error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, &net.DNSError{Err: "unrecognized address", Name: addr}
	}
	var family int
	if ipv4 := ip.To4(); ipv4 != nil {
		// use 4-byte representation
		ip = ipv4
		family = syscall.AF_INET
	} else {
		family = syscall.AF_INET6
	}
	hostnames, _, err := r.conn.ResolveAddress(ctx, 0, family, ip, 0)
	if err != nil {
		return nil, err
	}
	names = make([]string, len(hostnames))
	for i, name := range hostnames {
		names[i] = fullyQualified(name.Hostname)
	}
	return
}

// LookupIP looks up host for the given network using the systemd-resolved resolver.
// It returns a slice of that host's IP addresses of the type specified by network.
// network must be one of "ip", "ip4" or "ip6".
func (r *Resolver) LookupIP(ctx context.Context, network, host string) ([]net.IP, error) {
	if host == "" {
		return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	var family int
	switch network {
	case "ip":
		family = syscall.AF_UNSPEC
	case "ip4":
		family = syscall.AF_INET
	case "ip6":
		family = syscall.AF_INET6
	default:
		return nil, errors.New("bad network")
	}
	addresses, _, _, err := r.conn.ResolveHostname(ctx, 0, host, family, 0)
	if err != nil {
		return nil, err
	}
	addrs := make([]net.IP, len(addresses))
	for i, addr := range addresses {
		addrs[i] = addr.Address
	}
	return addrs, nil
}

// LookupIPAddr looks up host using the systemd-resolved resolver.
// It returns a slice of that host's IPv4 and IPv6 addresses.
func (r *Resolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	if host == "" {
		return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	addresses, _, _, err := r.conn.ResolveHostname(ctx, 0, host, syscall.AF_UNSPEC, 0)
	if err != nil {
		return nil, err
	}
	addrs := make([]net.IPAddr, len(addresses))
	for i, addr := range addresses {
		addrs[i] = net.IPAddr{
			IP: addr.Address,
		}
	}
	return addrs, nil
}

// LookupCNAME returns the canonical name for the given host.
func (r *Resolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	if host == "" {
		return "", &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	records, _, err := r.conn.ResolveRecord(ctx, 0, host, dns.ClassINET, dns.Type(dns.TypeCNAME), 0)
	if err != nil {
		return "", err
	}
	for _, record := range records {
		recordCNAME, err := record.CNAME()
		if err != nil {
			return "", err
		}
		return recordCNAME.Target, nil
	}
	return "", &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
}

// LookupMX returns the DNS MX records for the given domain name sorted by preference.
func (r *Resolver) LookupMX(ctx context.Context, name string) ([]*net.MX, error) {
	if !isDomainName(name) {
		return nil, &net.DNSError{Err: "no such host", Name: name, IsNotFound: true}
	}
	records, _, err := r.conn.ResolveRecord(ctx, 0, name, dns.ClassINET, dns.Type(dns.TypeMX), 0)
	if err != nil {
		return nil, err
	}
	mxs := make([]*net.MX, len(records))
	for i, record := range records {
		mx, err := record.MX()
		if err != nil {
			return nil, err
		}
		mxs[i] = &net.MX{
			Host: mx.Mx,
			Pref: mx.Preference,
		}
	}
	sort.Slice(mxs, func(i, j int) bool {
		return mxs[i].Pref < mxs[j].Pref
	})
	return mxs, nil
}

// LookupNS returns the DNS NS records for the given domain name.
func (r *Resolver) LookupNS(ctx context.Context, name string) ([]*net.NS, error) {
	if !isDomainName(name) {
		return nil, &net.DNSError{Err: "no such host", Name: name, IsNotFound: true}
	}
	records, _, err := r.conn.ResolveRecord(ctx, 0, name, dns.ClassINET, dns.Type(dns.TypeNS), 0)
	if err != nil {
		return nil, err
	}
	nss := make([]*net.NS, len(records))
	for i, record := range records {
		ns, err := record.NS()
		if err != nil {
			return nil, err
		}
		nss[i] = &net.NS{
			Host: ns.Ns,
		}
	}
	return nss, nil
}

// LookupPort looks up the port for the given network and service.
func (r *Resolver) LookupPort(ctx context.Context, network, service string) (port int, err error) {
	err = errNotSupported
	return
}

// LookupSRV tries to resolve an SRV query of the given service, protocol, and domain name.
// The proto is "tcp" or "udp". The returned records are sorted by priority and randomized by weight within a priority.
func (r *Resolver) LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error) {
	err = errNotSupported
	return
}

// LookupTXT returns the DNS TXT records for the given domain name.
func (r *Resolver) LookupTXT(ctx context.Context, name string) ([]string, error) {
	if !isDomainName(name) {
		return nil, &net.DNSError{Err: "no such host", Name: name, IsNotFound: true}
	}
	records, _, err := r.conn.ResolveRecord(ctx, 0, name, dns.ClassINET, dns.Type(dns.TypeTXT), 0)
	if err != nil {
		return nil, err
	}
	txts := make([]string, 0, len(records))
	for _, record := range records {
		txt, err := record.TXT()
		if err != nil {
			return nil, err
		}
		txts = append(txts, txt.Txt...)
	}
	return txts, nil
}

func fullyQualified(s string) string {
	if !strings.HasSuffix(s, ".") {
		s = s + "."
	}
	return s
}

func isDomainName(s string) bool {
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}
	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}
	return nonNumeric
}
