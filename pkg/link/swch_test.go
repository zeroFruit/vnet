package link_test

import (
	"encoding/gob"
	"testing"

	"github.com/zeroFruit/vnet/pkg/link/na"

	"github.com/zeroFruit/vnet/pkg/errors"
	"github.com/zeroFruit/vnet/pkg/link"
)

type mockSwitchPort struct {
	sendFunc func(pkt []byte) error
	id       link.Id
}

func (si *mockSwitchPort) GetLink() *link.Link {
	return nil
}

func (si *mockSwitchPort) AttachLink(link *link.Link) error {
	return nil
}

func (si *mockSwitchPort) Transmit(pkt []byte) error {
	return si.sendFunc(pkt)
}

func (si *mockSwitchPort) Id() link.Id {
	return si.id
}

func (si *mockSwitchPort) InternalAddress() link.Addr {
	return ""
}

func (si *mockSwitchPort) Register(id link.Id) {}

func (si *mockSwitchPort) Registered() bool {
	return true
}

func decodeFrame(t *testing.T, frm []byte) na.Frame {
	frmDec := link.NewFrameDecoder()
	frame, err := frmDec.Decode(frm)
	if err != nil {
		t.Fatalf("failed to decode frame: %v", err)
	}
	return frame
}

func assertFrame(t *testing.T, frm []byte, src string, dest string, payload string) {
	frame := decodeFrame(t, frm)
	if !frame.Src.Equal(link.AddrFromStr(src)) {
		t.Fatalf("expected frame Src is '%s', but got '%s'", src, frame.Src)
	}
	if !frame.Dest.Equal(link.AddrFromStr(dest)) {
		t.Fatalf("expected frame Src is '%s', but got '%s'", dest, frame.Dest)
	}
	if string(frame.Payload) != "hello" {
		t.Fatalf("expected frame Payload is '%s', but got '%s'", payload, string(frame.Payload))
	}
}

// TestSwitch_Forward_WhenAddressNotExist tests when there's no entry on switch table
// in this case, when frame comes from sender with interface id 'x', then broadcasts
// packets to all interfaces except 'x'
func TestSwitch_Forward_WhenAddressNotExist(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()
	swch := link.NewSwitchWithTable(table)
	sp1 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			t.Fail()
			return nil
		},
		id: link.Id("1"),
	}
	sp2 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		id: link.Id("2"),
	}
	sp3 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		id: link.Id("3"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sp1)).
		Happen(swch.Attach(sp2)).
		Happen(swch.Attach(sp3)).
		Return(); err != nil || len(swch.PortList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", na.Frame{
		Src:     link.AddrFromStr("11-11-11-11-11-11"),
		Dest:    link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward frame: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when frame comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnSameId(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()

	// dest address exists on table with interface addr "1"
	table.Update("1", link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sp1 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		id: link.Id("1"),
	}
	sp2 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		id: link.Id("2"),
	}
	sp3 := &mockSwitchPort{
		sendFunc: func(pkt []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		id: link.Id("3"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sp1)).
		Happen(swch.Attach(sp2)).
		Happen(swch.Attach(sp3)).
		Return(); err != nil || len(swch.PortList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", na.Frame{
		Src:     link.AddrFromStr("11-11-11-11-11-11"),
		Dest:    link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward frame: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when frame comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnDifferentId(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()

	// dest address exists on table with interface id "2"
	table.Update("2", link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sp1 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		id: link.Id("1"),
	}
	sp2 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			// frame need to be forwarded
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		id: link.Id("2"),
	}
	sp3 := &mockSwitchPort{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		id: link.Id("3"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sp1)).
		Happen(swch.Attach(sp2)).
		Happen(swch.Attach(sp3)).
		Return(); err != nil || len(swch.PortList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward("1", na.Frame{
		Src:     link.AddrFromStr("11-11-11-11-11-11"),
		Dest:    link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward frame: %v", err)
	}

	entry1, ok := table.LookupById("1")
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}
