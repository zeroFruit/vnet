package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/arp"
	"github.com/zeroFruit/vnet/link"
	"github.com/zeroFruit/vnet/net"
)

type mockInterface struct {
	hwAddr   link.Addr
	netAddr  net.Addr
	sendFunc func(payload arp.Payload) error
}

func (i *mockInterface) Send(payload arp.Payload) error {
	return i.sendFunc(payload)
}
func (i *mockInterface) HwAddr() link.Addr {
	return i.hwAddr
}
func (i *mockInterface) NetAddr() net.Addr {
	return i.netAddr
}

type mockNode struct {
	intfList []arp.Interface
}

func (n *mockNode) Interfaces() []arp.Interface {
	return n.intfList
}

func TestService_Broadcast(t *testing.T) {
	itf := &mockInterface{
		hwAddr:  "11:11:11:11:11:11",
		netAddr: "1.1.1.1",
	}
	itf.sendFunc = func(payload arp.Payload) error {
		if payload.SHwAddr != "11:11:11:11:11:11" {
			t.Fatalf("expected sender hw address is \"11:11:11:11:11:11\", but got %s", payload.SHwAddr)
		}
		if payload.SNetAddr != "1.1.1.1" {
			t.Fatalf("expected sender net address is \"1.1.1.1\", but got %s", payload.SHwAddr)
		}
		if payload.THwAddr != link.BroadcastAddr {
			t.Fatalf("expected target hw address is \"%s\", but got %s", link.BroadcastAddr, payload.SHwAddr)
		}
		if payload.TNetAddr != "2.2.2.2" {
			t.Fatalf("expected target net address is \"2.2.2.2\", but got %s", payload.TNetAddr)
		}
		return nil
	}
	node := &mockNode{
		intfList: []arp.Interface{itf},
	}
	service := arp.New(node)
	errs := service.Broadcast("2.2.2.2")
	if len(errs) != 0 {
		t.Fatalf("expected errs length is 0 but got %d", len(errs))
	}
}

func TestService_Recv(t *testing.T) {
	node := &mockNode{}
	table := arp.NewTable()
	service := arp.NewWithTable(node, table)
	req := arp.Request("11:11:11:11:11:11", "1.1.1.1", "2.2.2.2")
	service.Recv(req)
}
