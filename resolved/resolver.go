package resolved

import (
	"context"
	"errors"
	"net"
	"net/http"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
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
func (r *Resolver) Close() error {
	return r.conn.Close()
}

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
			address = "[" + addr.Address.String() + "]"
			break
		}
		address = addr.Address.String()
	}
	return r.dialer.DialContext(ctx, network, address+":"+port)
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
	return mxs, nil
}

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
