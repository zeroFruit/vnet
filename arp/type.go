package arp

import (
	"github.com/zeroFruit/vnet/net"
	"github.com/zeroFruit/vnet/phy"
)

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
	SHA   phy.Addr
	SPA   net.Addr
	THA   phy.Addr
	TPA   net.Addr
}

func Request(sha phy.Addr, spa net.Addr, tpa net.Addr) Payload {
	return Payload{
		HType: Ethernet,
		PType: IPV4,
		Op:    Req,
		SHA:   sha,
		SPA:   spa,
		TPA:   tpa,
	}
}

func Response() Payload {
	return Payload{}
}
