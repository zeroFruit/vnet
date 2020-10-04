package arp

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/zeroFruit/vnet/pkg/net"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/link"
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
	HType    HardwareType  `json:"hType"`
	PType    ProtocolType  `json:"pType"`
	HLen     uint8         `json:"hLen"`
	PLen     uint8         `json:"pLen"`
	Op       Operation     `json:"op"`
	SHwAddr  types.HwAddr  `json:"sHwAddr"`
	SNetAddr types.NetAddr `json:"sNetAddr"`
	THwAddr  types.HwAddr  `json:"tHwAddr"`
	TNetAddr types.NetAddr `json:"tNetAddr"`
}

func Request(sha types.HwAddr, sna types.NetAddr, tna types.NetAddr) Payload {
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

func Response(sha types.HwAddr, sna types.NetAddr, tha types.HwAddr, tna types.NetAddr) Payload {
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

func (p Payload) Encode() ([]byte, error) {
	var buf bytes.Buffer
	gob.Register(link.Addr(""))
	gob.Register(net.Addr(""))
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("failed to encode ARP payload: %v", err)
	}
	return buf.Bytes(), nil
}

func DecodePayload(b []byte) (Payload, error) {
	var payload Payload
	if err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&payload); err != nil {
		return Payload{}, fmt.Errorf("failed to decode ARP payload: %v", err)
	}
	return payload, nil
}
