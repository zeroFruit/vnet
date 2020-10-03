package net

import (
	"fmt"

	"github.com/zeroFruit/vnet/link"
)

type Interface struct {
	Addr Addr
	l    link.Interface
}

func NewInterface(itf link.Interface, addr Addr) *Interface {
	return &Interface{
		Addr: addr,
		l:    itf,
	}
}

func (i *Interface) HwAddress() link.Addr {
	return i.l.Address()
}

type Node struct {
	l       *link.Node
	ItfList []*Interface
}

func NewNode(l *link.Node) *Node {
	return &Node{
		l: l,
	}
}

func (n *Node) InterfaceOfAddr(addr Addr) (*Interface, error) {
	for _, itf := range n.ItfList {
		if itf.Addr.Equal(addr) {
			return itf, nil
		}
	}
	return nil, fmt.Errorf("interface of address'%s' not exist", addr.String())
}

func (n *Node) SendTo(addr Addr, buf []byte) error {
	itf, err := n.InterfaceOfAddr(addr)
	if err != nil {
		return err
	}
	if err := itf.l.Send(buf); err != nil {
		return err
	}
	return nil
}
