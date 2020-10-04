package arp_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/net"
)

type mockNetInterface struct {
	hwAddr   link.Addr
	netAddr  net.Addr
	sendFunc func(pkt []byte) error
}

func (i *mockNetInterface) Send(pkt []byte) error {
	return i.sendFunc(pkt)
}
func (i *mockNetInterface) HwAddress() types.HwAddr {
	return i.hwAddr
}
func (i *mockNetInterface) NetAddress() types.NetAddr {
	return i.netAddr
}

type mockNetNode struct {
<<<<<<< HEAD
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
=======
	intfList []types.NetInterface
}

func (n *mockNetNode) Interfaces() []types.NetInterface {
	return n.intfList
}

func TestService_Broadcast(t *testing.T) {
>>>>>>> 274bb3e... feat: implement data receving part from network layer to link layer
	itf := &mockNetInterface{
		hwAddr:  "11-11-11-11-11-11",
		netAddr: "1.1.1.1",
	}
<<<<<<< HEAD
	enc.encodeFunc = func(payload arp.Payload) ([]byte, error) {
=======
	itf.sendFunc = func(pkt []byte) error {
		payload, err := arp.DecodePayload(pkt)
		if err != nil {
			t.Fatalf("failed to unmarshal ARP payload: %v", err)
		}
>>>>>>> 274bb3e... feat: implement data receving part from network layer to link layer
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
<<<<<<< HEAD
		return []byte("hello"), nil
	}
	itf.sendFunc = func(pkt []byte) error {
		if string(pkt) != "hello" {
			t.Fatalf("expected pkt message is 'hello', but got %s", string(pkt))
		}
		return nil
	}
	node := &mockNetNode{
		itfList: []types.NetInterface{itf},
	}
	service := arp.New(node, enc)
=======
		return nil
	}
	node := &mockNetNode{
		intfList: []types.NetInterface{itf},
	}
	service := arp.New(arp.AdaptNode(node))
>>>>>>> 274bb3e... feat: implement data receving part from network layer to link layer
	errs := service.Broadcast(net.AddrFromStr("2.2.2.2"))
	if len(errs) != 0 {
		t.Fatalf("expected errs length is 0 but got %d, %v", len(errs), errs)
	}
}

func TestService_Recv(t *testing.T) {
<<<<<<< HEAD
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
				sendFunc: func(pkt []byte) error {
					if string(pkt) != "result" {
						t.Fatalf("expected pkt is 'result', but got '%s'", string(pkt))
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
=======
	node := &mockNetNode{}
	table := arp.NewTable()
	service := arp.NewWithTable(arp.AdaptNode(node), table)
	req := arp.Request(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("1.1.1.1"), net.AddrFromStr("2.2.2.2"))
	service.Recv(req)
>>>>>>> 274bb3e... feat: implement data receving part from network layer to link layer
}
