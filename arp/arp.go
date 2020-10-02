package arp

type Interface interface {
}

type Node interface {
	Interfaces() []Interface
}

type Service struct {
	self  Node
	table Table
}

func (s *Service) Broadcast() error {
	return nil
}

func (s *Service) Reply() error {
	return nil
}
