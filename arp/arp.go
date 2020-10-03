package arp

import (
	"github.com/zeroFruit/vnet/net"
)

type Service interface {
	Broadcast(tna net.Addr) []error
	Reply() error
	Recv(payload Payload)
}

type service struct {
	self  Node
	table *Table
}

func New(node Node) *service {
	return &service{
		self:  node,
		table: NewTable(),
	}
}

func NewWithTable(node Node, table *Table) *service {
	return &service{
		self:  node,
		table: table,
	}
}

func (s *service) Broadcast(tna net.Addr) []error {
	errs := make([]error, 0)

	// FIXME: selectively choose interface depending on the target network address
	for _, itf := range s.self.Interfaces() {
		payload := Request(itf.HwAddr(), itf.NetAddr(), tna)
		if err := itf.Send(payload); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (s *service) Reply() error {
	return nil
}

func (s *service) Recv(payload Payload) {
}
