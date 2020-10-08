package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/net"

	"github.com/zeroFruit/vnet/pkg/arp"
)

func TestTable(t *testing.T) {
	table := arp.NewTable()

	// insert test
	table.Update(arp.KeyValue(net.AddrFromStr("1.1.1.1"), link.AddrFromStr("11-11-11-11-11-11")))
	entry1, ok := table.Lookup(arp.Key{NetAddr: net.AddrFromStr("1.1.1.1")})
	if !ok {
		t.Fatalf("failed to lookup arp table entry")
	}
	if !entry1.NetAddr.Equal(net.AddrFromStr("1.1.1.1")) && !entry1.HwAddr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expected entry1 value is IP: 1.1.1.1, Addr: 11-11-11-11-11-11, "+
			"but got: IP: %s, Addr: %s", entry1.NetAddr, entry1.HwAddr)
	}

	// update test
	table.Update(arp.KeyValue(net.AddrFromStr("1.1.1.1"), link.AddrFromStr("11-11-11-11-11-12")))
	entry2, ok := table.Lookup(arp.Key{NetAddr: net.AddrFromStr("1.1.1.1")})
	if !ok {
		t.Fatalf("failed to lookup arp table entry")
	}
	if !entry2.NetAddr.Equal(net.AddrFromStr("1.1.1.1")) && !entry2.HwAddr.Equal(link.AddrFromStr("11-11-11-11-11-12")) {
		t.Fatalf("expected entry1 value is IP: 1.1.1.1, Addr: 11-11-11-11-11-12, "+
			"but got: IP: %s, Addr: %s", entry1.NetAddr, entry1.HwAddr)
	}

	// delete test
	table.Forget(arp.Key{NetAddr: net.AddrFromStr("1.1.1.1")})
	_, ok = table.Lookup(arp.Key{NetAddr: net.AddrFromStr("1.1.1.1")})
	if ok {
		t.Fatalf("failed to delete arp table entry")
	}
}
