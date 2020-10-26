package types

type NetInterface interface {
	Transmit(pkt []byte) error
	HwAddress() HwAddr
	NetAddress() NetAddr
}

type NetNode interface {
	Interfaces() []NetInterface
}
