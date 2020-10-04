package types

type NetInterface interface {
	Send(pkt []byte) error
	HwAddress() HwAddr
	NetAddress() NetAddr
}
<<<<<<< HEAD

type NetNode interface {
	Interfaces() []NetInterface
}
=======
>>>>>>> 274bb3e... feat: implement data receving part from network layer to link layer
