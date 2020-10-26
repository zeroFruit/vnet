package arp

import (
	"github.com/zeroFruit/vnet/pkg/errors"
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/link/na"
	"github.com/zeroFruit/vnet/pkg/types"
)

type Service interface {
	Broadcast(tna types.NetAddr) error
	Recv(payload Payload) error
}

type service struct {
	node   types.NetNode
	table  *Table
	plEnc  PayloadEncoder
	frmEnc *link.FrameEncoder
}

func New(node types.NetNode, plEnc PayloadEncoder) Service {
	return &service{
		node:   node,
		table:  NewTable(),
		plEnc:  plEnc,
		frmEnc: link.NewFrameEncoder(),
	}
}

func NewWithTable(node types.NetNode, plEnc PayloadEncoder, table *Table) Service {
	return &service{
		node:   node,
		table:  table,
		plEnc:  plEnc,
		frmEnc: link.NewFrameEncoder(),
	}
}

func (s *service) Broadcast(tna types.NetAddr) error {
	errs := errors.Multiple()
	// FIXME: selectively choose interface depending on the target network address
	for _, itf := range s.node.Interfaces() {
		pl, err := s.plEnc.Encode(
			Request(itf.HwAddress(), itf.NetAddress(), tna))
		if err != nil {
			errs = errs.Happen(err)
			continue
		}
		frame, err := s.frmEnc.Encode(na.Frame{
			Src:     itf.HwAddress(),
			Dest:    link.BroadcastAddr,
			Payload: pl,
		})
		if err != nil {
			errs = errs.Happen(err)
			continue
		}
		if err := itf.Transmit(frame); err != nil {
			errs = errs.Happen(err)
		}
	}
	return errs.Return()
}

func (s *service) Recv(payload Payload) error {
	if _, ok := s.table.Lookup(Key{NetAddr: payload.SNetAddr}); ok {
		s.table.Update(KeyValue(payload.SNetAddr, payload.SHwAddr))
	}
	var itf types.NetInterface
	matched := false
	for _, i := range s.node.Interfaces() {
		if i.NetAddress().Equal(payload.TNetAddr) {
			matched = true
			itf = i
		}
	}
	if !matched {
		return nil
	}
	s.table.Update(KeyValue(payload.SNetAddr, payload.SHwAddr))
	if payload.Op == Reply {
		return nil
	}
	pl, err := s.plEnc.Encode(
		Response(itf.HwAddress(), itf.NetAddress(), payload.SHwAddr, payload.SNetAddr))
	if err != nil {
		return err
	}
	frame, err := s.frmEnc.Encode(na.Frame{
		Src:     itf.HwAddress(),
		Dest:    payload.SHwAddr,
		Payload: pl,
	})
	if err != nil {
		return err
	}
	return itf.Transmit(frame)
}
