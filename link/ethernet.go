package link

import "github.com/zeroFruit/vnet/physical"

const Octet = 8

type EthernetType [2]byte

type EthernetPayload interface{}

type EthernetHeader struct {
	dst physical.HardwareAddr
	src physical.HardwareAddr
	fcs uint
}

// TODO: implement me
func ShouldHandleFrame(intf physical.Interface, header EthernetHeader) bool {

	return false
}
