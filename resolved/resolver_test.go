package resolved

import (
	"context"
	"net"
	"net/http"
	"sort"
	"testing"
)

// In order to run the test make sure that systemd-resolved resolver query the same dns server as the go one.

const (
	lookupHost             = "google.com"
	lookupAddr4            = "142.250.178.142"
	lookupAddr6            = "2a00:1450:4007:81a::200e"
	lookupCNAMEHost        = "en.wikipedia.org"
	lookupSRVDomain        = "google.com"
	lookupSRVService       = "xmpp-server"
	lookupSRVProto         = "tcp"
	lookupSRVServiceDomain = "_xmpp-server._tcp.google.com"
	getUrl                 = "https://google.com"
)

func TestLookupHost(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdAddrs, err := sysdResolver.LookupHost(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goAddrs, err := goResolver.LookupHost(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goAddrs) != len(sysdAddrs) {
		t.Fatal("len(goAddrs) != len(sysdAddrs)", len(goAddrs), len(sysdAddrs))
	}
	sort.Strings(sysdAddrs)
	sort.Strings(goAddrs)
	for i, sAddr := range sysdAddrs {
		goAddr := goAddrs[i]
		if goAddr != sAddr {
			t.Error("goAddr != sAddr", goAddr, sAddr)
		}
	}
}

func TestLookupAddr4(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdNames, err := sysdResolver.LookupAddr(ctx, lookupAddr4)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goNames, err := goResolver.LookupAddr(ctx, lookupAddr4)
	if err != nil {
		t.Fatal(err)
	}
	if len(goNames) != len(sysdNames) {
		t.Fatal("len(goNames) != len(sysdNames)", len(goNames), len(sysdNames))
	}
	sort.Strings(goNames)
	sort.Strings(sysdNames)
	for i, sName := range sysdNames {
		goName := goNames[i]
		if goName != sName {
			t.Error("goName != sName", goName, sName)
		}
	}
}

func TestLookupAddr6(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdNames, err := sysdResolver.LookupAddr(ctx, lookupAddr6)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goNames, err := goResolver.LookupAddr(ctx, lookupAddr6)
	if err != nil {
		t.Fatal(err)
	}
	if len(goNames) != len(sysdNames) {
		t.Fatal("len(goNames) != len(sysdNames)", len(goNames), len(sysdNames))
	}
	sort.Strings(goNames)
	sort.Strings(sysdNames)
	for i, sName := range sysdNames {
		goName := goNames[i]
		if goName != sName {
			t.Error("goName != sName", goName, sName)
		}
	}
}

func TestLookupIP(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdAddrs, err := sysdResolver.LookupIP(ctx, "ip", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goAddrs, err := goResolver.LookupIP(ctx, "ip", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goAddrs) != len(sysdAddrs) {
		t.Fatal("len(goAddrs) != len(sysdAddrs)", len(goAddrs), len(sysdAddrs))
	}
	sort.Slice(sysdAddrs, func(i, j int) bool {
		return sysdAddrs[i].String() < sysdAddrs[j].String()
	})
	sort.Slice(goAddrs, func(i, j int) bool {
		return goAddrs[i].String() < goAddrs[j].String()
	})
	for i, sAddr := range sysdAddrs {
		goAddr := goAddrs[i]
		if goAddr.String() != sAddr.String() {
			t.Error("goAddr != sAddr", goAddr, sAddr)
		}
	}
}

func TestLookupIP4(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdAddrs, err := sysdResolver.LookupIP(ctx, "ip4", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goAddrs, err := goResolver.LookupIP(ctx, "ip4", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goAddrs) != len(sysdAddrs) {
		t.Fatal("len(goAddrs) != len(sysdAddrs)", len(goAddrs), len(sysdAddrs))
	}
	sort.Slice(sysdAddrs, func(i, j int) bool {
		return sysdAddrs[i].String() < sysdAddrs[j].String()
	})
	sort.Slice(goAddrs, func(i, j int) bool {
		return goAddrs[i].String() < goAddrs[j].String()
	})
	for i, sAddr := range sysdAddrs {
		goAddr := goAddrs[i]
		if goAddr.String() != sAddr.String() {
			t.Error("goAddr != sAddr", goAddr, sAddr)
		}
	}
}

func TestLookupIP6(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdAddrs, err := sysdResolver.LookupIP(ctx, "ip6", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goAddrs, err := goResolver.LookupIP(ctx, "ip6", lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goAddrs) != len(sysdAddrs) {
		t.Fatal("len(goAddrs) != len(sysdAddrs)", len(goAddrs), len(sysdAddrs))
	}
	sort.Slice(sysdAddrs, func(i, j int) bool {
		return sysdAddrs[i].String() < sysdAddrs[j].String()
	})
	sort.Slice(goAddrs, func(i, j int) bool {
		return goAddrs[i].String() < goAddrs[j].String()
	})
	for i, sAddr := range sysdAddrs {
		goAddr := goAddrs[i]
		if goAddr.String() != sAddr.String() {
			t.Error("goAddr != sAddr", goAddr, sAddr)
		}
	}
}

func TestLookupIPAddr(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdAddrs, err := sysdResolver.LookupIPAddr(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goAddrs, err := goResolver.LookupIPAddr(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goAddrs) != len(sysdAddrs) {
		t.Fatal("len(goAddrs) != len(sysdAddrs)", len(goAddrs), len(sysdAddrs))
	}
	sort.Slice(sysdAddrs, func(i, j int) bool {
		return sysdAddrs[i].String() < sysdAddrs[j].String()
	})
	sort.Slice(goAddrs, func(i, j int) bool {
		return goAddrs[i].String() < goAddrs[j].String()
	})
	for i, sAddr := range sysdAddrs {
		goAddr := goAddrs[i]
		if goAddr.String() != sAddr.String() {
			t.Error("goAddr != sAddr", goAddr, sAddr)
		}
		if goAddr.Zone != sAddr.Zone {
			t.Error("goAddr .Zone!= sAddr.Zone", goAddr.Zone, sAddr.Zone)
		}
	}
}

func TestLookupCNAME(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdCNAME, err := sysdResolver.LookupCNAME(ctx, lookupCNAMEHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goCNAME, err := goResolver.LookupCNAME(ctx, lookupCNAMEHost)
	if err != nil {
		t.Fatal(err)
	}
	if goCNAME != sysdCNAME {
		t.Error("goCNAME != sysdCNAME", goCNAME, sysdCNAME)
	}
}

func TestLookupMX(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdMxs, err := sysdResolver.LookupMX(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goMxs, err := goResolver.LookupMX(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goMxs) != len(sysdMxs) {
		t.Fatal("len(goMxs) != len(sysdMxs)", len(goMxs), len(sysdMxs))
	}
	for i, sMx := range sysdMxs {
		goMx := goMxs[i]
		if goMx.Host != sMx.Host {
			t.Error("goMx.Host != sMx.Host", goMx.Host, sMx.Host)
		}
		if goMx.Pref != sMx.Pref {
			t.Error("goMx.Pref != sMx.Pref", goMx.Pref, sMx.Pref)
		}
	}
}

func TestLookupNS(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdNss, err := sysdResolver.LookupNS(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goNss, err := goResolver.LookupNS(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goNss) != len(sysdNss) {
		t.Fatal("len(goNss) != len(sysdNss)", len(goNss), len(sysdNss))
	}
	sort.Slice(sysdNss, func(i, j int) bool {
		return sysdNss[i].Host < sysdNss[j].Host
	})
	sort.Slice(goNss, func(i, j int) bool {
		return goNss[i].Host < goNss[j].Host
	})
	for i, sNs := range sysdNss {
		goNs := goNss[i]
		if goNs.Host != sNs.Host {
			t.Error("goNs.Host != sNs.Host", goNs.Host, sNs.Host)
		}
	}
}

func TestLookupSRV(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdCNAME, sysdSrvs, err := sysdResolver.LookupSRV(ctx, lookupSRVService, lookupSRVProto, lookupSRVDomain)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goCNAME, goSrvs, err := goResolver.LookupSRV(ctx, lookupSRVService, lookupSRVProto, lookupSRVDomain)
	if err != nil {
		t.Fatal(err)
	}
	if sysdCNAME != goCNAME {
		t.Fatal("sysdCNAME != goCNAME", sysdCNAME, goCNAME)
	}
	if len(goSrvs) != len(sysdSrvs) {
		t.Fatal("len(goSrvs) != len(sysdSrvs)", len(goSrvs), len(sysdSrvs))
	}
	sort.Slice(sysdSrvs, func(i, j int) bool {
		return sysdSrvs[i].Target < sysdSrvs[j].Target
	})
	sort.Slice(goSrvs, func(i, j int) bool {
		return goSrvs[i].Target < goSrvs[j].Target
	})
	for i, sSrv := range sysdSrvs {
		goSrv := goSrvs[i]
		if goSrv.Target != sSrv.Target {
			t.Error("goSrv.Target != sSrv.Target", goSrv.Target, sSrv.Target)
		}
		if goSrv.Port != sSrv.Port {
			t.Error("goSrv.Port != sSrv.Port", goSrv.Port, sSrv.Port)
		}
		if goSrv.Priority != sSrv.Priority {
			t.Error("goSrv.Priority != sSrv.Priority", goSrv.Priority, sSrv.Priority)
		}
		if goSrv.Weight != sSrv.Weight {
			t.Error("goSrv.Weight != sSrv.Weight", goSrv.Weight, sSrv.Weight)
		}
	}
}

func TestLookupTXT(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdTxts, err := sysdResolver.LookupTXT(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goTxts, err := goResolver.LookupTXT(ctx, lookupHost)
	if err != nil {
		t.Fatal(err)
	}
	if len(goTxts) != len(sysdTxts) {
		t.Fatal("len(goTxts) != len(sysdTxts)", len(goTxts), len(sysdTxts))
	}
	sort.Strings(sysdTxts)
	sort.Strings(goTxts)
	for i, sTxt := range sysdTxts {
		goTxt := goTxts[i]
		if goTxt != sTxt {
			t.Error("goTxt != sTxt", goTxt, sTxt)
		}
	}
}

func BenchmarkLookupHostGoResolver(b *testing.B) {
	r := &net.Resolver{}
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_, err := r.LookupHost(ctx, lookupHost)
		if err != nil {
			b.Error(err)
			continue
		}
	}
}

func BenchmarkLookupHostSystemdResolver(b *testing.B) {
	r, err := NewResolver()
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_, err := r.LookupHost(ctx, lookupHost)
		if err != nil {
			b.Error(err)
			continue
		}
	}
}

func BenchmarkLookupAddrGoResolver(b *testing.B) {
	r := &net.Resolver{}
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_, err := r.LookupAddr(ctx, lookupAddr4)
		if err != nil {
			b.Error(err)
			continue
		}
	}
}

func BenchmarkLookupAddrSystemdResolver(b *testing.B) {
	r, err := NewResolver()
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_, err := r.LookupAddr(ctx, lookupAddr4)
		if err != nil {
			b.Error(err)
			continue
		}
	}
}

func BenchmarkHTTPClientGoResolver(b *testing.B) {
	httpCli := &http.Client{}
	for n := 0; n < b.N; n++ {
		resp, err := httpCli.Get(getUrl)
		if err != nil {
			b.Error(err)
			continue
		}
		resp.Body.Close()
	}
}

func BenchmarkHTTPClientSystemdResolver(b *testing.B) {
	r, err := NewResolver()
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()
	httpCli := r.HTTPClient()
	for n := 0; n < b.N; n++ {
		resp, err := httpCli.Get(getUrl)
		if err != nil {
			b.Error(err)
			continue
		}
		resp.Body.Close()
	}
}
