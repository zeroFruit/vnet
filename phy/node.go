package phy

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/zeroFruit/vnet/phy/internal"
)

const defaultIP internal.Addr = "127.0.0.1"

type Id struct {
	Name string
}

func (i Id) Get() string {
	return i.Name
}

func (i Id) Equal(id Id) bool {
	return i.Name == id.Get()
}

func IdOf(name string) Id {
	return Id{
		Name: name,
	}
}

type Addr string

type Link struct {
	id    Id
	cost  uint
	intf1 Interface
	intf2 Interface
}

func NewLink(id string, cost uint) *Link {
	return &Link{
		id:    IdOf(id),
		cost:  cost,
		intf1: Interface{},
		intf2: Interface{},
	}
}

func (l *Link) GetId() Id {
	return l.id
}

func (l *Link) GetCost() uint {
	return l.cost
}

func (l *Link) SetInterface(intf Interface) error {
	if l.intf1 == (Interface{}) {
		l.intf1 = intf
		return nil
	}
	if l.intf2 == (Interface{}) {
		l.intf2 = intf
		return nil
	}
	return fmt.Errorf("link is full with interfaces on %s", l.id)
}

func (l *Link) GetOtherInterface(intfId Id) (Interface, error) {
	if l.intf1.id.Equal(intfId) {
		return l.intf2, nil
	}
	if l.intf2.id.Equal(intfId) {
		return l.intf1, nil
	}
	return Interface{}, fmt.Errorf("cannot find interface id %s", intfId.Name)
}

type Interface struct {
	id         Id
	ip         internal.Addr
	internalIP internal.Addr
	port       int
	hwAddr     Addr
	link       *Link
	adapter    NetworkAdapter
	dataSink   chan<- *Datagram
	quit       chan struct{}
}

func NewInterface(id string, ip internal.Addr, port int, hwAddr Addr, dataSink chan<- *Datagram) *Interface {
	intf := &Interface{
		id:         IdOf(id),
		ip:         ip,
		port:       port,
		internalIP: defaultIP,
		hwAddr:     hwAddr,
		dataSink:   dataSink,
		quit:       make(chan struct{}),
	}
	return intf
}

func (i *Interface) GetId() Id {
	return i.id
}

func (i *Interface) GetLink() *Link {
	return i.link
}

func (i *Interface) AttachLink(link *Link) error {
	if i.link != nil {
		return errors.New("link already exist")
	}
	if err := link.SetInterface(*i); err != nil {
		return err
	}
	adapter, err := NewNetworkAdapter(i.internalIP, i.port)
	if err != nil {
		return err
	}
	i.adapter = adapter
	i.link = link
	go i.sink()
	return nil
}

func (i *Interface) Send(buf []byte) error {
	receiver, err := i.link.GetOtherInterface(i.id)
	if err != nil {
		return err
	}
	i.adapter.Send(buf, string(receiver.internalIP)+":"+strconv.Itoa(receiver.port))
	return nil
}

func (i *Interface) sink() {
	for {
		select {
		case data := <-i.adapter.Recv():
			data.From = i.id.Get()
			i.dataSink <- data
		case <-i.quit:
			return
		}
	}
}

func (i *Interface) shutdown() {
	i.quit <- struct{}{}
}

type Node struct {
	id       Id
	intfList []*Interface
	dataCh   chan *Datagram
	quit     chan struct{}
}

func NewNode(id string) *Node {
	n := &Node{
		id:       IdOf(id),
		intfList: make([]*Interface, 0),
		dataCh:   make(chan *Datagram), // TODO: set the buffer
		quit:     make(chan struct{}),
	}
	go n.listen()
	return n
}

func (n *Node) GetId() Id {
	return n.id
}

func (n *Node) DataSink() chan<- *Datagram {
	return n.dataCh
}

func (n *Node) GetInterfaceById(id Id) (*Interface, error) {
	for _, intf := range n.intfList {
		if intf.id.Equal(id) {
			return intf, nil
		}
	}
	return nil, fmt.Errorf("interface '%s' not exist", id.Get())
}

func (n *Node) AttachInterface(intf *Interface) error {
	n.intfList = append(n.intfList, intf)
	return nil
}

func (n *Node) SendTo(intfName string, buf []byte) error {
	intf, err := n.GetInterfaceById(IdOf(intfName))
	if err != nil {
		return err
	}
	if err := intf.Send(buf); err != nil {
		return err
	}
	return nil
}

func (n *Node) listen() {
	for {
		select {
		case data := <-n.dataCh:
			n.handleData(data)
		case <-n.quit:
			return
		}
	}
}

// TODO: implement me
func (n *Node) handleData(data *Datagram) error {
	return nil
}

func (n *Node) Shutdown() {
	for _, intf := range n.intfList {
		intf.shutdown()
	}
	n.quit <- struct{}{}
}
