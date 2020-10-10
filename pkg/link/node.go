package link

import (
	"errors"
	"fmt"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/link/na"
)

type LinkInterface interface {
	Interface
	InternalAddress() Addr
}

type Link struct {
	cost  uint
	intf1 LinkInterface
	intf2 LinkInterface
}

func NewLink(cost uint) *Link {
	return &Link{
		cost:  cost,
		intf1: nil,
		intf2: nil,
	}
}

func (l *Link) GetCost() uint {
	return l.cost
}

func (l *Link) SetInterface(itf LinkInterface) error {
	if l.intf1 == nil {
		l.intf1 = itf
		return nil
	}
	if l.intf2 == nil {
		l.intf2 = itf
		return nil
	}
	return fmt.Errorf("link is full with interfaces between %s, %s", l.intf1.Address(), l.intf2.Address())
}

func (l *Link) GetOtherInterface(addr types.HwAddr) (LinkInterface, error) {
	if l.intf1.Address().Equal(addr) {
		return l.intf2, nil
	}
	if l.intf2.Address().Equal(addr) {
		return l.intf1, nil
	}
	return nil, fmt.Errorf("cannot find other interface link attached by %s", addr)
}

type NetDatagramHandler interface {
	Handle(data *na.Datagram)
}

type Node struct {
	dataCh     chan *na.Datagram
	quit       chan struct{}
	ItfList    []Interface
	netHandler NetDatagramHandler
}

func NewNode() *Node {
	n := &Node{
		dataCh:     make(chan *na.Datagram), // TODO: set the buffer
		quit:       make(chan struct{}),
		ItfList:    make([]Interface, 0),
		netHandler: nil,
	}
	return n
}

func (n *Node) RegisterNetHandler(handler NetDatagramHandler) {
	n.netHandler = handler
}

func (n *Node) AttachInterface(itf Interface) {
	n.ItfList = append(n.ItfList, itf)
}

func (n *Node) InterfaceOfAddr(addr Addr) (Interface, error) {
	for _, itf := range n.ItfList {
		if itf.Address().Equal(addr) {
			return itf, nil
		}
	}
	return nil, fmt.Errorf("interface of address'%s' not exist", addr.String())
}

func (n *Node) DataSink() chan<- *na.Datagram {
	return n.dataCh
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

func (n *Node) handle(data *na.Datagram) error {
	if n.netHandler == nil {
		return errors.New("net handler is not registered")
	}
	n.netHandler.Handle(data)
	return nil
}

func (n *Node) Shutdown() {
	n.quit <- struct{}{}
}
