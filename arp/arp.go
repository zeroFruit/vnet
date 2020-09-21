package arp

import "github.com/zeroFruit/vnet/physical"

type HardwareType uint16

const Ethernet HardwareType = 1

type ProtocolType uint16

const IPV4 = 0x800

type Operation uint16

const (
	Req   Operation = 1
	Reply           = 2
)

type Payload struct {
	HType HardwareType
	PType ProtocolType
	HLen  uint8
	PLen  uint8
	Op    Operation
	SHA   physical.HardwareAddr
	SPA   physical.IP
	THA   physical.HardwareAddr
	TPA   physical.IP
}
