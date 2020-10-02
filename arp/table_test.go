package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/arp"
)

func TestTable(t *testing.T) {
	table := arp.NewTable()

	// insert test
	table.Update(arp.KeyValue("1.1.1.1", "11-11-11-11-11-11"))
	entry1, ok := table.Lookup(arp.Key{NetworkAddr: "1.1.1.1"})
	if !ok {
		t.Fatalf("failed to lookup arp table entry")
	}
	if entry1.NetworkAddr != "1.1.1.1" && entry1.HardwareAddr != "11-11-11-11-11-11" {
		t.Fatalf("expected entry1 value is IP: 1.1.1.1, HwAddr: 11-11-11-11-11-11, "+
			"but got: IP: %s, HwAddr: %s", entry1.NetworkAddr, entry1.HardwareAddr)
	}

	// update test
	table.Update(arp.KeyValue("1.1.1.1", "11-11-11-11-11-12"))
	entry2, ok := table.Lookup(arp.Key{NetworkAddr: "1.1.1.1"})
	if !ok {
		t.Fatalf("failed to lookup arp table entry")
	}
	if entry2.NetworkAddr != "1.1.1.1" && entry2.HardwareAddr != "11-11-11-11-11-12" {
		t.Fatalf("expected entry1 value is IP: 1.1.1.1, HwAddr: 11-11-11-11-11-12, "+
			"but got: IP: %s, HwAddr: %s", entry1.NetworkAddr, entry1.HardwareAddr)
	}

	// delete test
	table.Forget(arp.Key{NetworkAddr: "1.1.1.1"})
	_, ok = table.Lookup(arp.Key{NetworkAddr: "1.1.1.1"})
	if ok {
		t.Fatalf("failed to delete arp table entry")
	}
}
