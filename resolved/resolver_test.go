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
	lookupHost = "google.fr"
	lookupAddr = "142.250.75.227"
	getUrl     = "https://google.fr"
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

func TestLookupAddr(t *testing.T) {
	sysdResolver, err := NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	defer sysdResolver.Close()
	ctx := context.Background()
	sysdNames, err := sysdResolver.LookupAddr(ctx, lookupAddr)
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goNames, err := goResolver.LookupAddr(ctx, lookupAddr)
	if err != nil {
		t.Fatal(err)
	}
	if len(goNames) != len(sysdNames) {
		t.Fatal("len(goNames) != len(sysdNames)", len(goNames), len(sysdNames))
	}
	sort.Strings(sysdNames)
	sort.Strings(goNames)
	for i, sName := range sysdNames {
		goName := goNames[i]
		if goName != sName {
			t.Error("goName != sName", goName, sName)
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
	sysdCNAME, err := sysdResolver.LookupCNAME(ctx, "ig1-sismobox-01.sadm.ig-1.net")
	if err != nil {
		t.Fatal(err)
	}
	goResolver := &net.Resolver{}
	goCNAME, err := goResolver.LookupCNAME(ctx, "ig1-sismobox-01.sadm.ig-1.net")
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
	sort.Slice(goMxs, func(i, j int) bool {
		return goMxs[i].Host > goMxs[j].Host
	})
	sort.Slice(sysdMxs, func(i, j int) bool {
		return sysdMxs[i].Host > sysdMxs[j].Host
	})
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
	sort.Slice(goNss, func(i, j int) bool {
		return goNss[i].Host > goNss[j].Host
	})
	sort.Slice(sysdNss, func(i, j int) bool {
		return sysdNss[i].Host > sysdNss[j].Host
	})
	for i, sNs := range sysdNss {
		goNs := goNss[i]
		if goNs.Host != sNs.Host {
			t.Error("goNs.Host != sNs.Host", goNs.Host, sNs.Host)
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
		_, err := r.LookupAddr(ctx, lookupAddr)
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
		_, err := r.LookupAddr(ctx, lookupAddr)
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
