package link

import "github.com/zeroFruit/vnet/phy"

const Octet = 8

type EthernetType [2]byte

type EthernetPayload interface{}

type EthernetHeader struct {
	dst phy.Addr
	src phy.Addr
	fcs uint
}

// TODO: implement me
func ShouldHandleFrame(intf phy.Interface, header EthernetHeader) bool {

	return false
}
