package net_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/link"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/net"
)

func TestArpPayloadEncoderDecoder(t *testing.T) {
	payloads := []arp.Payload{
		arp.Request(link.AddrFromStr("11-11-11-11-11-11"),
			net.AddrFromStr("1.1.1.1"), net.AddrFromStr("2.2.2.2")),
		arp.Request(link.AddrFromStr("11-11-11-11-11-12"),
			net.AddrFromStr("1.1.1.2"), net.AddrFromStr("2.2.2.3")),
	}
	enc := net.NewArpPayloadEncoder()
	dec := net.NewArpPayloadDecoder()

	for i := 0; i < len(payloads); i++ {
		b, err := enc.Encode(payloads[i])
		if err != nil {
			t.Fatalf("failed to encode payload: %v", err)
		}
		result, err := dec.Decode(b)
		if err != nil {
			t.Fatalf("failed to decode payload: %v", err)
		}
		if payloads[i].SHwAddr != result.SHwAddr {
			t.Fatalf("expected sender hw addr %s, but got %s", payloads[i].SHwAddr, result.SHwAddr)
		}
		if payloads[i].SNetAddr != result.SNetAddr {
			t.Fatalf("expected sender net addr %s, but got %s", payloads[i].SNetAddr, result.SNetAddr)
		}
		if payloads[i].THwAddr != result.THwAddr {
			t.Fatalf("expected target hw addr %s, but got %s", payloads[i].THwAddr, result.THwAddr)
		}
		if payloads[i].TNetAddr != result.TNetAddr {
			t.Fatalf("expected target net addr %s, but got %s", payloads[i].TNetAddr, result.TNetAddr)
		}
	}

}
