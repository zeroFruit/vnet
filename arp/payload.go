package arp

import (
	"encoding/json"
	"fmt"

	"github.com/zeroFruit/vnet/link"
	"github.com/zeroFruit/vnet/net"
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
	HType    HardwareType `json:"hType"`
	PType    ProtocolType `json:"pType"`
	HLen     uint8        `json:"hLen"`
	PLen     uint8        `json:"pLen"`
	Op       Operation    `json:"op"`
	SHwAddr  link.Addr    `json:"sHwAddr"`
	SNetAddr net.Addr     `json:"sNetAddr"`
	THwAddr  link.Addr    `json:"tHwAddr"`
	TNetAddr net.Addr     `json:"tNetAddr"`
}

func Request(sha link.Addr, sna net.Addr, tna net.Addr) Payload {
	return Payload{
		HType:    Ethernet,
		PType:    IPV4,
		Op:       Req,
		SHwAddr:  sha,
		SNetAddr: sna,
		THwAddr:  link.BroadcastAddr,
		TNetAddr: tna,
	}
}

func Response(sha link.Addr, sna net.Addr, tha link.Addr, tna net.Addr) Payload {
	return Payload{
		HType:    Ethernet,
		PType:    IPV4,
		Op:       Reply,
		SHwAddr:  sha,
		SNetAddr: sna,
		THwAddr:  tha,
		TNetAddr: tna,
	}
}

func (p Payload) Marshal() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal ARP payload: %v", err))
	}
	return b
}
