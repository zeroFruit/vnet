package link_test

import (
	"encoding/gob"
	"github.com/zeroFruit/vnet/pkg/link/na"
	"testing"

	"github.com/zeroFruit/vnet/pkg/errors"
	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/link"
)

type mockInterface struct {
	sendFunc func(pkt []byte) error
	addr     types.HwAddr
}

func (si *mockInterface) GetLink() *link.Link {
	return nil
}

func (si *mockInterface) AttachLink(link *link.Link) error {
	return nil
}

func (si *mockInterface) Send(pkt []byte) error {
	return si.sendFunc(pkt)
}

func (si *mockInterface) Address() types.HwAddr {
	return si.addr
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
// in this case, when packet comes from sender with interface id 'x', then broadcasts
// packets to all interfaces except 'x'
func TestSwitch_Forward_WhenAddressNotExist(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()
	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockInterface{
		sendFunc: func(frm []byte) error {
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-01"),
	}
	sitf2 := &mockInterface{
		sendFunc: func(frm []byte) error {
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-02"),
	}
	sitf3 := &mockInterface{
		sendFunc: func(frm []byte) error {
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-03"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sitf1)).
		Happen(swch.Attach(sitf2)).
		Happen(swch.Attach(sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward(link.AddrFromStr("00-00-00-00-00-01"), na.Frame{
		Src: link.AddrFromStr("11-11-11-11-11-11"),
		Dest: link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById(link.AddrFromStr("00-00-00-00-00-01"))
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when packet comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnSameId(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()

	// dest address exists on table with interface addr "00-00-00-00-00-01"
	table.Update(link.AddrFromStr("00-00-00-00-00-01"), link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockInterface{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-01"),
	}
	sitf2 := &mockInterface{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-02"),
	}
	sitf3 := &mockInterface{
		sendFunc: func(pkt []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-03"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sitf1)).
		Happen(swch.Attach(sitf2)).
		Happen(swch.Attach(sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward(link.AddrFromStr("00-00-00-00-00-01"), na.Frame{
		Src: link.AddrFromStr("11-11-11-11-11-11"),
		Dest: link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById(link.AddrFromStr("00-00-00-00-00-01"))
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}

// TestSwitch_Forward_WhenReceiverExistOnSameId tests when packet comes from the same
// interface id with the id that exists on table with key of receiver address
func TestSwitch_Forward_WhenReceiverExistOnDifferentId(t *testing.T) {
	gob.Register(link.Addr(""))

	table := link.NewSwitchTable()

	// dest address exists on table with interface id "2"
	table.Update(link.AddrFromStr("00-00-00-00-00-02"), link.AddrFromStr("33-33-33-33-33-33"))

	swch := link.NewSwitchWithTable(table)
	sitf1 := &mockInterface{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-01"),
	}
	sitf2 := &mockInterface{
		sendFunc: func(frm []byte) error {
			// frame need to be forwarded
			assertFrame(t, frm, "11-11-11-11-11-11", "33-33-33-33-33-33", "hello")
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-02"),
	}
	sitf3 := &mockInterface{
		sendFunc: func(frm []byte) error {
			// frame must be discard
			t.Fail()
			return nil
		},
		addr: link.AddrFromStr("00-00-00-00-00-03"),
	}
	if err := errors.Multiple().
		Happen(swch.Attach(sitf1)).
		Happen(swch.Attach(sitf2)).
		Happen(swch.Attach(sitf3)).
		Return(); err != nil && len(swch.ItfList) != 3 {
		t.Fatalf("failed to attach switch interface: %v", err)
	}
	if err := swch.Forward(link.AddrFromStr("00-00-00-00-00-01"), na.Frame{
		Src: link.AddrFromStr("11-11-11-11-11-11"),
		Dest: link.AddrFromStr("33-33-33-33-33-33"),
		Payload: []byte("hello"),
	}); err != nil {
		t.Fatalf("failed to forward packet: %v", err)
	}

	entry1, ok := table.LookupById(link.AddrFromStr("00-00-00-00-00-01"))
	if !ok && !entry1.Addr.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expect '11-11-11-11-11-11' to exist in table, but not exist")
	}
}
