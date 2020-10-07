package arp

import (
	"github.com/zeroFruit/vnet/pkg/types"
)

type Service interface {
	Broadcast(tna types.NetAddr) []error
	Reply() error
	Recv(payload Payload) error
}

type service struct {
	node    types.NetNode
	table   *Table
	encoder PayloadEncoder
}

func New(node types.NetNode, encoder PayloadEncoder) Service {
	return &service{
		node:    node,
		table:   NewTable(),
		encoder: encoder,
	}
}

func NewWithTable(node types.NetNode, encoder PayloadEncoder, table *Table) Service {
	return &service{
		node:    node,
		table:   table,
		encoder: encoder,
	}
}

func (s *service) Broadcast(tna types.NetAddr) []error {
	errs := make([]error, 0)
	// FIXME: selectively choose interface depending on the target network address
	for _, itf := range s.node.Interfaces() {
		pkt, err := s.encoder.Encode(
			Request(itf.HwAddress(), itf.NetAddress(), tna))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := itf.Send(pkt); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (s *service) Reply() error {
	return nil
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
	pkt, err := s.encoder.Encode(
		Response(itf.HwAddress(), itf.NetAddress(), payload.SHwAddr, payload.SNetAddr))
	if err != nil {
		return err
	}
	return itf.Send(pkt)
}
