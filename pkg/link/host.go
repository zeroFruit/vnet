package link

import (
	"errors"
	"fmt"

	"github.com/zeroFruit/vnet/pkg/link/na"

	"github.com/zeroFruit/vnet/pkg/types"
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

func NewNode() *Host {
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
	if err := n.Interface.Send(frame); err != nil {
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