package net_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/net"

	"github.com/zeroFruit/vnet/pkg/link"
)

type mockInterface struct {
	Addr link.Addr
}

func (i *mockInterface) GetLink() *link.Link {
	return nil
}

func (i *mockInterface) AttachLink(link *link.Link) error {
	return nil
}

func (i *mockInterface) Send(pkt []byte) error {
	return nil
}

func (i *mockInterface) Address() link.Addr {
	return i.Addr
}

func TestNode_UpdateAddr_WhenInterfaceExist(t *testing.T) {
	li := &mockInterface{
		Addr: "11-11-11-11-11-11",
	}
	ln := link.NewNode()
	ln.AttachInterface(li)

	nn := net.NewNode(ln)
	nn.ItfList = append(nn.ItfList, net.NewInterface(li, net.AddrFromStr("1.1.1.1")))

	err := nn.UpdateAddr(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("2.2.2.2"))
	if err != nil {
		t.Fatalf("failed to update address: %v", err)
	}
	if len(nn.ItfList) != 1 {
		t.Fatalf("expected interface list length is 1 but got %d", len(nn.ItfList))
	}
	if !nn.ItfList[0].Addr.Equal(net.AddrFromStr("2.2.2.2")) {
		t.Fatalf("expected address is '2.2.2.2' but got %s", nn.ItfList[0].Addr)
	}
}

func TestNode_UpdateAddr_WhenInterfaceNotExist(t *testing.T) {
	li := &mockInterface{
		Addr: "11-11-11-11-11-11",
	}
	ln := link.NewNode()
	ln.AttachInterface(li)

	nn := net.NewNode(ln)

	err := nn.UpdateAddr(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("2.2.2.2"))
	if err != nil {
		t.Fatalf("failed to update address: %v", err)
	}
	if len(nn.ItfList) != 1 {
		t.Fatalf("expected interface list length is 1 but got %d", len(nn.ItfList))
	}
	if !nn.ItfList[0].Addr.Equal(net.AddrFromStr("2.2.2.2")) {
		t.Fatalf("expected address is '2.2.2.2' but got %s", nn.ItfList[0].Addr)
	}
}

func TestNode_UpdateAddr_WhenHwInterfaceNotEnough(t *testing.T) {
	// there's no hw interface attached
	ln := link.NewNode()

	nn := net.NewNode(ln)

	err := nn.UpdateAddr(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("2.2.2.2"))
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}
