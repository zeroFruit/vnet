package arp

import (
	"github.com/zeroFruit/vnet/link"
	"github.com/zeroFruit/vnet/net"
)

type Interface interface {
	Send(payload Payload) error
	HwAddr() link.Addr
	NetAddr() net.Addr
}

type AdaptedInterface struct {
	itf *net.Interface
}

func AdaptInterface(itf *net.Interface) Interface {
	return &AdaptedInterface{
		itf: itf,
	}
}

func (i *AdaptedInterface) Send(payload Payload) error {
	// TODO: marshal payload and convert into byte slice
	return nil
}

func (i *AdaptedInterface) HwAddr() link.Addr {
	return i.itf.HwAddress()
}

func (i *AdaptedInterface) NetAddr() net.Addr {
	return i.itf.Addr
}

type Node interface {
	Interfaces() []Interface
}

type AdaptedNode struct {
	node     *net.Node
	intfList []Interface
}

func AdaptNode(node *net.Node) *AdaptedNode {
	intfList := make([]Interface, 0)
	for _, itf := range node.ItfList {
		intfList = append(intfList, AdaptInterface(itf))
	}
	return &AdaptedNode{
		node:     node,
		intfList: intfList,
	}
}

func (n *AdaptedNode) Interfaces() []Interface {
	return n.intfList
}
