package types

type NetInterface interface {
	Send(pkt []byte) error
	HwAddress() HwAddr
	NetAddress() NetAddr
}
