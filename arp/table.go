package arp

import (
	"github.com/zeroFruit/vnet/net"
	"github.com/zeroFruit/vnet/phy"
)

type Key struct {
	NetworkAddr net.Addr
}

type Entry struct {
	NetworkAddr  net.Addr
	HardwareAddr phy.Addr
}

func KeyValue(na net.Addr, ha phy.Addr) (Key, Entry) {
	return Key{NetworkAddr: na}, Entry{NetworkAddr: na, HardwareAddr: ha}
}

type Table struct {
	entries map[Key]Entry
}

func NewTable() *Table {
	return &Table{
		entries: make(map[Key]Entry),
	}
}

func (t *Table) Update(key Key, entry Entry) {
	t.entries[key] = entry
}

func (t *Table) Lookup(key Key) (Entry, bool) {
	entry, ok := t.entries[key]
	return entry, ok
}

func (t *Table) Forget(key Key) {
	delete(t.entries, key)
}
