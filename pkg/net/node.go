package net

import (
	"errors"
	"fmt"
	"log"

	"github.com/zeroFruit/vnet/pkg/link/na"

	"github.com/zeroFruit/vnet/pkg/arp"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/link"
)

type Interface struct {
	Addr types.NetAddr
	hw   link.AnonymInterface
}

func NewInterface(hw link.AnonymInterface, addr types.NetAddr) *Interface {
	return &Interface{
		Addr: addr,
		hw:   hw,
	}
}

func (i *Interface) HwAddress() types.HwAddr {
	return i.hw.Address()
}

func (i *Interface) NetAddress() types.NetAddr {
	return i.Addr
}

func (i *Interface) Send(pkt []byte) error {
	return i.hw.Send(pkt)
}

type Node struct {
	hw      *link.Node
	ItfList []*Interface
	arp     arp.Service
	decoder arp.PayloadDecoder
}

func NewNode(hw *link.Node) *Node {
	return &Node{
		hw:      hw,
		ItfList: make([]*Interface, 0),
		decoder: NewArpPayloadDecoder(),
	}
}

func (n *Node) UpdateAddr(hwAddr types.HwAddr, addr types.NetAddr) error {
	if ok := n.updateExistAddr(hwAddr, addr); ok {
		return nil
	}
	if ok := n.registerNewAddr(hwAddr, addr); ok {
		return nil
	}
	return fmt.Errorf("failed to register address '%s' on hw address '%s', not enough hw interface", addr, hwAddr)
}

func (n *Node) updateExistAddr(hwAddr types.HwAddr, addr types.NetAddr) (ok bool) {
	ok = false
	for _, itf := range n.ItfList {
		if itf.HwAddress().Equal(hwAddr) {
			itf.Addr = addr
			ok = true
		}
	}
	return
}

func (n *Node) registerNewAddr(hwAddr types.HwAddr, addr types.NetAddr) (ok bool) {
	ok = false
	for _, hwItf := range n.hw.ItfList {
		if hwItf.Address().Equal(hwAddr) {
			itf := NewInterface(hwItf, addr)
			n.ItfList = append(n.ItfList, itf)
			ok = true
		}
	}
	return
}

func (n *Node) InterfaceOfAddr(addr Addr) (*Interface, error) {
	for _, itf := range n.ItfList {
		if itf.Addr.Equal(addr) {
			return itf, nil
		}
	}
	return nil, fmt.Errorf("interface of address'%s' not exist", addr.String())
}

func (n *Node) Send(addr Addr, pkt []byte) error {
	itf, err := n.InterfaceOfAddr(addr)
	if err != nil {
		return err
	}
	if err := itf.Send(pkt); err != nil {
		return err
	}
	return nil
}

func (n *Node) Interfaces() []types.NetInterface {
	r := make([]types.NetInterface, 0)
	for _, itf := range n.ItfList {
		r = append(r, types.NetInterface(itf))
	}
	return r
}

func (n *Node) RegisterArp(arp arp.Service) {
	n.arp = arp
}

func (n *Node) Handle(data *na.Datagram) {
	payload, err := n.decoder.Decode(data.Buf)
	if err == nil {
		if err := n.handleArp(payload); err != nil {
			log.Fatalf("failed to handle ARP packet: %v", err)
		}
	}
	// handle other protocol packet
}

func (n *Node) handleArp(payload arp.Payload) error {
	if n.arp == nil {
		return errors.New("ARP module not registered")
	}
	return n.arp.Recv(payload)
}
