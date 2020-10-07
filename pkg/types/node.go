package types

type NetInterface interface {
	Send(pkt []byte) error
	HwAddress() HwAddr
	NetAddress() NetAddr
}

type NetNode interface {
	Interfaces() []NetInterface
}
