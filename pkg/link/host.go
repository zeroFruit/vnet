package link

import (
	"errors"
	"fmt"

	"github.com/zeroFruit/vnet/pkg/link/na"

	"github.com/zeroFruit/vnet/pkg/types"
)

type Id string

func (i Id) Equal(id Id) bool {
	return i == id
}

func (i Id) Empty() bool {
	return i == ""
}

// EndPoint represents point of link. Link is the channel to pass data to end-point
// either side to the opposite
type EndPoint interface {
	Id() Id
	InternalAddress() Addr
	GetLink() *Link
	AttachLink(link *Link) error
}

type Link struct {
	cost uint
	ep1  EndPoint
	ep2  EndPoint
}

func NewLink(cost uint) *Link {
	return &Link{
		cost: cost,
		ep1:  nil,
		ep2:  nil,
	}
}

func (l *Link) GetCost() uint {
	return l.cost
}

func (l *Link) AttachEndpoint(ep EndPoint) error {
	if l.ep1 == nil {
		l.ep1 = ep
		return nil
	}
	if l.ep2 == nil {
		l.ep2 = ep
		return nil
	}
	return fmt.Errorf("link is full with interfaces between %s, %s", l.ep1.Id(), l.ep2.Id())
}

func (l *Link) Opposite(id Id) (EndPoint, error) {
	if l.ep1.Id().Equal(id) {
		return l.ep2, nil
	}
	if l.ep2.Id().Equal(id) {
		return l.ep1, nil
	}
	return nil, fmt.Errorf("cannot find other interface link attached by %s", id)
}

// NetHandler receives serialized frame payload. With this, doing some high-level protocol
type NetHandler interface {
	Handle(pl []byte)
}

type Host struct {
	quit       chan struct{}
	Interface  Interface
	netHandler NetHandler
	frmEnc     *FrameEncoder
	frmDec     *FrameDecoder
}

func NewHost() *Host {
	n := &Host{
		quit:       make(chan struct{}),
		Interface:  nil,
		netHandler: nil,
		frmEnc:     NewFrameEncoder(),
		frmDec:     NewFrameDecoder(),
	}
	return n
}

func (n *Host) RegisterNetHandler(handler NetHandler) {
	n.netHandler = handler
}

func (n *Host) AttachInterface(itf Interface) {
	n.Interface = itf
}

// Send make frame with payload and transfer to destination
func (n *Host) Send(dest types.HwAddr, pl []byte) error {
	frame, err := n.frmEnc.Encode(na.Frame{
		Src:     n.Interface.Address(),
		Dest:    dest,
		Payload: pl,
	})
	if err != nil {
		return err
	}
	if err := n.Interface.Transmit(frame); err != nil {
		return err
	}
	return nil
}

func (n *Host) handle(fd *na.FrameData) error {
	if n.netHandler == nil {
		return errors.New("net handler is not registered")
	}
	frame, err := n.frmDec.Decode(fd.Buf)
	if err != nil {
		return err
	}
	n.netHandler.Handle(frame.Payload)
	return nil
}

func (n *Host) Shutdown() {
	n.quit <- struct{}{}
}
