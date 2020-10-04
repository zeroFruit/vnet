package arp

import (
	"github.com/zeroFruit/vnet/pkg/types"
)

type Key struct {
	NetAddr types.NetAddr
}

type Entry struct {
	NetAddr types.NetAddr
	HwAddr  types.HwAddr
}

func KeyValue(na types.NetAddr, ha types.HwAddr) (Key, Entry) {
	return Key{NetAddr: na}, Entry{NetAddr: na, HwAddr: ha}
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
