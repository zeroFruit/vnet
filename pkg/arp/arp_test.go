package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/net"
)

type mockNetInterface struct {
	hwAddr   link.Addr
	netAddr  net.Addr
	sendFunc func(pkt []byte) error
}

func (i *mockNetInterface) Send(pkt []byte) error {
	return i.sendFunc(pkt)
}
func (i *mockNetInterface) HwAddress() types.HwAddr {
	return i.hwAddr
}
func (i *mockNetInterface) NetAddress() types.NetAddr {
	return i.netAddr
}

type mockNetNode struct {
	intfList []types.NetInterface
}

func (n *mockNetNode) Interfaces() []types.NetInterface {
	return n.intfList
}

func TestService_Broadcast(t *testing.T) {
	itf := &mockNetInterface{
		hwAddr:  "11-11-11-11-11-11",
		netAddr: "1.1.1.1",
	}
	itf.sendFunc = func(pkt []byte) error {
		payload, err := arp.DecodePayload(pkt)
		if err != nil {
			t.Fatalf("failed to unmarshal ARP payload: %v", err)
		}
		if !payload.SHwAddr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
			t.Fatalf("expected sender hw address is 11-11-11-11-11-11, but got %s", payload.SHwAddr)
		}
		if !payload.SNetAddr.Equal(net.AddrFromStr("1.1.1.1")) {
			t.Fatalf("expected sender net address is 1.1.1.1, but got %s", payload.SHwAddr)
		}
		if !payload.THwAddr.Equal(link.BroadcastAddr) {
			t.Fatalf("expected target hw address is %s, but got %s", link.BroadcastAddr, payload.SHwAddr)
		}
		if !payload.TNetAddr.Equal(net.AddrFromStr("2.2.2.2")) {
			t.Fatalf("expected target net address is 2.2.2.2, but got %s", payload.TNetAddr)
		}
		return nil
	}
	node := &mockNetNode{
		intfList: []types.NetInterface{itf},
	}
	service := arp.New(arp.AdaptNode(node))
	errs := service.Broadcast(net.AddrFromStr("2.2.2.2"))
	if len(errs) != 0 {
		t.Fatalf("expected errs length is 0 but got %d, %v", len(errs), errs)
	}
}

func TestService_Recv(t *testing.T) {
	node := &mockNetNode{}
	table := arp.NewTable()
	service := arp.NewWithTable(arp.AdaptNode(node), table)
	req := arp.Request(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("1.1.1.1"), net.AddrFromStr("2.2.2.2"))
	service.Recv(req)
}
