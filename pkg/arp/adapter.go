package arp

import "github.com/zeroFruit/vnet/pkg/types"

type AdaptedInterface struct {
	itf types.NetInterface
}

func AdaptInterface(itf types.NetInterface) *AdaptedInterface {
	return &AdaptedInterface{
		itf: itf,
	}
}

func (i *AdaptedInterface) Send(payload Payload) error {
	b, err := payload.Encode()
	if err != nil {
		return err
	}
	return i.itf.Send(b)
}

func (i *AdaptedInterface) HwAddr() types.HwAddr {
	return i.itf.HwAddress()
}

func (i *AdaptedInterface) NetAddr() types.NetAddr {
	return i.itf.NetAddress()
}

type NetNode interface {
	Interfaces() []types.NetInterface
}

type AdaptedNode struct {
	node     NetNode
	intfList []*AdaptedInterface
}

func AdaptNode(node NetNode) *AdaptedNode {
	intfList := make([]*AdaptedInterface, 0)
	for _, itf := range node.Interfaces() {
		intfList = append(intfList, AdaptInterface(itf))
	}
	return &AdaptedNode{
		node:     node,
		intfList: intfList,
	}
}

func (n *AdaptedNode) Interfaces() []*AdaptedInterface {
	return n.intfList
}
