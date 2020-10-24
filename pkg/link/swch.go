package link

import (
	"fmt"
	"log"
	"time"

	"github.com/zeroFruit/vnet/pkg/errors"
	"github.com/zeroFruit/vnet/pkg/link/na"

	"github.com/zeroFruit/vnet/pkg/types"
)

type ForwardEntry struct {
	// Incoming is interface hardware address attached to switch
	Incoming types.HwAddr

	// Addr is destination node address
	Addr     types.HwAddr

	// Time is timestamp when this entry is created
	Time     time.Time
}

type FrameForwardTable struct {
	entries []ForwardEntry
}

func NewSwitchTable() *FrameForwardTable {
	return &FrameForwardTable{
		entries: make([]ForwardEntry, 0),
	}
}

func (t *FrameForwardTable) Update(incoming types.HwAddr, addr types.HwAddr) {
	idxToRemove := -1
	for i, e := range t.entries {
		if e.Addr.Equal(addr) {
			idxToRemove = i
		}
	}
	if idxToRemove != -1 {
		t.entries = append(t.entries[:idxToRemove], t.entries[idxToRemove+1:]...)
	}
	t.entries = append(t.entries, ForwardEntry{Incoming: incoming, Addr: addr, Time: time.Now()})
}

func (t *FrameForwardTable) LookupByAddr(key types.HwAddr) (ForwardEntry, bool) {
	for _, e := range t.entries {
		if e.Addr.Equal(key) {
			return e, true
		}
	}
	return ForwardEntry{}, false
}

func (t *FrameForwardTable) LookupById(incoming types.HwAddr) (ForwardEntry, bool) {
	for _, e := range t.entries {
		if e.Incoming.Equal(incoming) {
			return e, true
		}
	}
	return ForwardEntry{}, false
}

func (t *FrameForwardTable) Entries() []ForwardEntry {
	return t.entries
}

// FrameForwarder forwards frame based on where this frame comes from and frame destination
type FrameForwarder interface {
	Forward(incoming types.HwAddr, frame na.Frame) error
}

type Switch struct {
	ItfList map[types.HwAddr]Interface
	Table   *FrameForwardTable
	frmDec  *FrameDecoder
	frmEnc  *FrameEncoder
	quit    chan struct{}
}

func NewSwitch() *Switch {
	return &Switch{
		ItfList: make(map[types.HwAddr]Interface),
		Table:   NewSwitchTable(),
		frmDec:  NewFrameDecoder(),
		frmEnc:  NewFrameEncoder(),
		quit:    make(chan struct{}),
	}
}

func NewSwitchWithTable(table *FrameForwardTable) *Switch {
	return &Switch{
		ItfList: make(map[types.HwAddr]Interface),
		Table:   table,
		quit:    make(chan struct{}),
	}
}

func (s *Switch) handle(fd *na.FrameData) error {
	frame, err := s.frmDec.Decode(fd.Buf)
	if err != nil {
		return err
	}
	return s.Forward(fd.Incoming, frame)
}

func (s *Switch) Attach(itf Interface) error {
	if _, ok := s.ItfList[itf.Address()]; ok {
		return fmt.Errorf("already exist interface: %s", itf.Address())
	}
	s.ItfList[itf.Address()] = itf
	return nil
}

// Forward receives address of interface it receives frame, address of sender
// and frame to send to receiver. Based on id, address it determines whether to
// broadcast frame or forward it to others, otherwise just discard frame.
func (s *Switch) Forward(incoming types.HwAddr, frame na.Frame) error {
	s.Table.Update(incoming, frame.Src)
	frm, err := s.frmEnc.Encode(frame)
	if err != nil {
		return err
	}
	entry, ok := s.Table.LookupByAddr(frame.Dest)
	if !ok {
		return s.broadcastExcept(incoming, frm)
	}
	if entry.Incoming.Equal(incoming) {
		log.Printf("discard frame from id: %s, src: %s, dest: %s\n", incoming, frame.Src, frame.Dest)
		return nil
	}
	return s.ItfList[entry.Incoming].Send(frm)
}

// broadcastExcept sends frame to other interfaces except the interface
// with the id value given by parameter
func (s *Switch) broadcastExcept(incoming types.HwAddr, frm []byte) error {
	err := errors.Multiple()
	for addr, itf := range s.ItfList {
		if !addr.Equal(incoming) {
			err = err.Happen(itf.Send(frm))
		}
	}
	return err.Return()
}

func (s *Switch) Shutdown() {
	s.quit <- struct{}{}
}
