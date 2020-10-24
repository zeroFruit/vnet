package net

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/zeroFruit/vnet/pkg/link"

	"github.com/zeroFruit/vnet/pkg/arp"
)

<<<<<<< HEAD
<<<<<<< HEAD
type ArpPayloadEncoder struct{}
=======
type ArpPayloadEncoder struct {}
>>>>>>> c752f79... fix: replace Datagram into Frame, fix switch core logic
=======
type ArpPayloadEncoder struct{}
>>>>>>> dbc75fc... fix: update switch logic, add link layer integration tests

func NewArpPayloadEncoder() *ArpPayloadEncoder {
	gob.Register(link.Addr(""))
	gob.Register(Addr(""))
	return &ArpPayloadEncoder{}
}

func (e *ArpPayloadEncoder) Encode(payload arp.Payload) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := gob.NewEncoder(buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode ARP payload: %v", err)
	}
	b := buf.Bytes()
	return b, nil
}

type ArpPayloadDecoder struct{}

func NewArpPayloadDecoder() *ArpPayloadDecoder {
	return &ArpPayloadDecoder{}
}

func (d *ArpPayloadDecoder) Decode(b []byte) (arp.Payload, error) {
	var payload arp.Payload
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	if err := decoder.Decode(&payload); err != nil {
		return arp.Payload{}, fmt.Errorf("failed to decode ARP payload: %v", err)
	}
	return payload, nil
}
