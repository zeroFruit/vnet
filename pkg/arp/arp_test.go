package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/link"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/net"
)

type mockNetInterface struct {
	hwAddr   link.Addr
	netAddr  net.Addr
	sendFunc func(pkt []byte) error
}

func (i *mockNetInterface) Transmit(pkt []byte) error {
	return i.sendFunc(pkt)
}
func (i *mockNetInterface) HwAddress() types.HwAddr {
	return i.hwAddr
}
func (i *mockNetInterface) NetAddress() types.NetAddr {
	return i.netAddr
}

type mockNetNode struct {
	itfList []types.NetInterface
}

func (n *mockNetNode) Interfaces() []types.NetInterface {
	return n.itfList
}

type mockPayloadEncoder struct {
	encodeFunc func(payload arp.Payload) ([]byte, error)
}

func (e *mockPayloadEncoder) Encode(payload arp.Payload) ([]byte, error) {
	return e.encodeFunc(payload)
}

func TestService_Broadcast(t *testing.T) {
	enc := &mockPayloadEncoder{}
	itf := &mockNetInterface{
		hwAddr:  "11-11-11-11-11-11",
		netAddr: "1.1.1.1",
	}
	enc.encodeFunc = func(payload arp.Payload) ([]byte, error) {
		if !payload.SHwAddr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
			t.Fatalf("expected sender hw address is 11-11-11-11-11-11, but got %s", payload.SHwAddr)
		}
		if !payload.SNetAddr.Equal(net.AddrFromStr("1.1.1.1")) {
			t.Fatalf("expected sender net address is 1.1.1.1, but got %s", payload.SHwAddr)
		}
		if !payload.THwAddr.Equal(link.BroadcastAddr) {
			t.Fatalf("expected target hw address is %s, but got %s", link.BroadcastAddr, payload.SHwAddr)
		}
		if !payload.TNetAddr.Equal(net.AddrFromStr("2.2.2.2")) {
			t.Fatalf("expected target net address is 2.2.2.2, but got %s", payload.TNetAddr)
		}
		return []byte("hello"), nil
	}
	itf.sendFunc = func(frame []byte) error {
		f, err := link.NewFrameDecoder().Decode(frame)
		if err != nil {
			t.Fatalf("failed to decode frame: %v", err)
		}
		if string(f.Payload) != "hello" {
			t.Fatalf("expected pkt message is 'hello', but got %s", string(f.Payload))
		}
		return nil
	}
	node := &mockNetNode{
		itfList: []types.NetInterface{itf},
	}
	service := arp.New(node, enc)
	err := service.Broadcast(net.AddrFromStr("2.2.2.2"))
	if err != nil {
		t.Fatalf("failed to broadcast ARP message: %v", err)
	}
}

func TestService_Recv(t *testing.T) {
	for _, tt := range []struct {
		desc                 string
		itf                  *mockNetInterface
		enc                  *mockPayloadEncoder
		createTableFunc      func() *arp.Table
		expectedTableUpdated bool
	}{
		{
			"when entry not exist and net address not matched",
			&mockNetInterface{
				hwAddr:  "33-33-33-33-33-33",
				netAddr: "3.3.3.3",
				sendFunc: func(pkt []byte) error {
					t.FailNow()
					return nil
				},
			},
			&mockPayloadEncoder{
				encodeFunc: func(payload arp.Payload) ([]byte, error) {
					t.FailNow()
					return nil, nil
				},
			},
			func() *arp.Table {
				return arp.NewTable()
			},
			false,
		},
		{
			"when entry exist and net address not matched",
			&mockNetInterface{
				hwAddr:  "33-33-33-33-33-33",
				netAddr: "3.3.3.3",
				sendFunc: func(pkt []byte) error {
					t.FailNow()
					return nil
				},
			},
			&mockPayloadEncoder{
				encodeFunc: func(payload arp.Payload) ([]byte, error) {
					t.FailNow()
					return nil, nil
				},
			},
			func() *arp.Table {
				table := arp.NewTable()
				table.Update(arp.KeyValue(net.AddrFromStr("1.1.1.1"), link.AddrFromStr("99-99-99-99-99-99")))
				return table
			},
			true,
		},
		{
			"when entry not exist and net address matched",
			&mockNetInterface{
				hwAddr:  "22-22-22-22-22-22",
				netAddr: "2.2.2.2",
				sendFunc: func(frame []byte) error {
					f, err := link.NewFrameDecoder().Decode(frame)
					if err != nil {
						t.Fatalf("failed to decode frame: %v", err)
					}
					if string(f.Payload) != "result" {
						t.Fatalf("expected pkt is 'result', but got '%s'", string(f.Payload))
					}
					return nil
				},
			},
			&mockPayloadEncoder{
				encodeFunc: func(payload arp.Payload) ([]byte, error) {
					if payload.Op != arp.Reply {
						t.Fatalf("ARP Op is not Reply")
					}
					if !payload.SHwAddr.Equal(link.AddrFromStr("22-22-22-22-22-22")) {
						t.Fatalf("expected SHwAddr is 22-22-22-22-22-22, but got %s", payload.SHwAddr)
					}
					if !payload.SNetAddr.Equal(net.AddrFromStr("2.2.2.2")) {
						t.Fatalf("expected SHwAddr is 2.2.2.2, but got %s", payload.SNetAddr)
					}
					return []byte("result"), nil
				},
			},
			func() *arp.Table {
				return arp.NewTable()
			},
			true,
		},
	} {
		t.Log(tt.desc)
		node := &mockNetNode{
			itfList: []types.NetInterface{tt.itf},
		}
		table := tt.createTableFunc()
		service := arp.NewWithTable(node, tt.enc, table)
		err := service.Recv(
			arp.Request(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("1.1.1.1"), net.AddrFromStr("2.2.2.2")))
		if err != nil {
			t.Fatalf("failed to receive ARP packet: %v", err)
		}
		if _, ok := table.Lookup(arp.Key{NetAddr: net.AddrFromStr("1.1.1.1")}); ok != tt.expectedTableUpdated {
			t.Fatalf("table entry should not be updated")
		}
	}
}
