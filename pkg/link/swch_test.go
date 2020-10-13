package link_test

import (
	"github.com/zeroFruit/vnet/pkg/errors"
	"testing"

	"github.com/zeroFruit/vnet/pkg/link"
)

type mockAnonymInterface struct {
	sendFunc func(pkt []byte) error
}

func (si *mockAnonymInterface) GetLink() *link.Link {
	return nil
}

func (si *mockAnonymInterface) AttachLink(link *link.Link) error {
	return nil
}

func (si *mockAnonymInterface) Send(pkt []byte) error {
	return si.sendFunc(pkt)
}

// TestSwitch_Forward_WhenAddressNotExist tests when there's no entry on switch table
// in this case, when packet comes from sender with interface id 'x', then broadcasts
// packets to all interfaces except 'x'
func TestSwitch_Forward_WhenAddressNotExist(t *testing.T) {
	table := link.NewSwitchTable()
	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			t.Fail()
			return nil
		},
	}
	sitf2 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			if string(pkt) != "hello" {
				t.Fatalf("expected pkt value is 'hello', but got '%s'", string(pkt))
			}
			return nil
		},
	}
	sitf3 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			if string(pkt) != "hello" {
				t.Fatalf("expected pkt value is 'hello', but got '%s'", string(pkt))
			}
			return nil
		},
	}
	if err := errors.Multiple().
		Happen(swch.Attach("1", sitf1)).
		Happen(swch.Attach("2", sitf2)).
		Happen(swch.Attach("3", sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", link.AddrFromStr("11-11-11-11-11-11"), link.AddrFromStr("33-33-33-33-33-33"), []byte("hello")); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when packet comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnSameId(t *testing.T) {
	table := link.NewSwitchTable()

	// dest address exists on table with interface id "1"
	table.Update("1", link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet must be discard
			t.Fail()
			return nil
		},
	}
	sitf2 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet must be discard
			t.Fail()
			return nil
		},
	}
	sitf3 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet must be discard
			t.Fail()
			return nil
		},
	}
	if err := errors.Multiple().
		Happen(swch.Attach("1", sitf1)).
		Happen(swch.Attach("2", sitf2)).
		Happen(swch.Attach("3", sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", link.AddrFromStr("11-11-11-11-11-11"), link.AddrFromStr("33-33-33-33-33-33"), []byte("hello")); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when packet comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnDifferentId(t *testing.T) {
	table := link.NewSwitchTable()

	// dest address exists on table with interface id "2"
	table.Update("2", link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet must be discard
			t.Fail()
			return nil
		},
	}
	sitf2 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet need to be forwarded
			if string(pkt) != "hello" {
				t.Fatalf("expected pkt value is 'hello', but got '%s'", string(pkt))
			}
			return nil
		},
	}
	sitf3 := &mockAnonymInterface{
		sendFunc: func(pkt []byte) error {
			// packet must be discard
			t.Fail()
			return nil
		},
	}
	if err := errors.Multiple().
		Happen(swch.Attach("1", sitf1)).
		Happen(swch.Attach("2", sitf2)).
		Happen(swch.Attach("3", sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", link.AddrFromStr("11-11-11-11-11-11"), link.AddrFromStr("33-33-33-33-33-33"), []byte("hello")); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}
