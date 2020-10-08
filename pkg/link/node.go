package link

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/zeroFruit/vnet/pkg/link/internal"
)

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

func (l *Link) GetOtherInterface(addr Addr) (LinkInterface, error) {
	if l.intf1.Address().Equal(addr) {
		return l.intf2, nil
	}
	if l.intf2.Address().Equal(addr) {
		return l.intf1, nil
	}
	return nil, fmt.Errorf("cannot find other interface link attached by %s", addr)
}

type Interface interface {
	GetLink() *Link
	AttachLink(link *Link) error
	Send(pkt []byte) error
	Address() Addr
}

type LinkInterface interface {
	Interface
	InternalAddress() Addr
}

type DatagramHandler interface {
	handle(datagram *Datagram) error
}

type UDPBasedInterface struct {
	internalIP   internal.Addr
	internalPort int
	Addr         Addr
	link         *Link
	adapter      NetworkAdapter
	handler      DatagramHandler
	quit         chan struct{}
}

func NewInterface(port int, hwAddr Addr, handler DatagramHandler) Interface {
	itf := &UDPBasedInterface{
		internalPort: port,
		internalIP:   internal.DefaultAddr,
		Addr:         hwAddr,
		handler:      handler,
		quit:         make(chan struct{}),
	}
	return itf
}

func (i *UDPBasedInterface) GetLink() *Link {
	return i.link
}

func (i *UDPBasedInterface) AttachLink(link *Link) error {
	if i.link != nil {
		return errors.New("link already exist")
	}
	if err := link.SetInterface(i); err != nil {
		return err
	}
	adapter, err := NewNetworkAdapter(i.internalIP, i.internalPort)
	if err != nil {
		return err
	}
	i.adapter = adapter
	i.link = link
	go i.sink()
	return nil
}

func (i *UDPBasedInterface) Send(buf []byte) error {
	receiver, err := i.link.GetOtherInterface(i.Addr)
	if err != nil {
		return err
	}
	i.adapter.Send(buf, receiver.InternalAddress().String())
	return nil
}

func (i *UDPBasedInterface) sink() {
	for {
		select {
		case data := <-i.adapter.Recv():
			data.From = i.Addr.String()
			if err := i.handler.handle(data); err != nil {
				log.Fatal(err)
			}
		case <-i.quit:
			return
		}
	}
}

func (i *UDPBasedInterface) Address() Addr {
	return i.Addr
}

func (i *UDPBasedInterface) InternalAddress() Addr {
	return Addr(string(i.internalIP) + ":" + strconv.Itoa(i.internalPort))
}

func (i *UDPBasedInterface) shutdown() {
	i.quit <- struct{}{}
}

type NetDatagramHandler interface {
	Handle(data *Datagram)
}

type Node struct {
	dataCh     chan *Datagram
	quit       chan struct{}
	ItfList    []Interface
	netHandler NetDatagramHandler
}

func NewNode() *Node {
	n := &Node{
		dataCh:     make(chan *Datagram), // TODO: set the buffer
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

func (n *Node) DataSink() chan<- *Datagram {
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

func (n *Node) handle(data *Datagram) error {
	if n.netHandler == nil {
		return errors.New("net handler is not registered")
	}
	n.netHandler.Handle(data)
	return nil
}

func (n *Node) Shutdown() {
	n.quit <- struct{}{}
}
