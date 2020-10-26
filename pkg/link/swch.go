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
	// Incoming is port id attached to switch
	Incoming Id

	// Addr is destination node address
	Addr types.HwAddr

	// Time is timestamp when this entry is created
	Time time.Time
}

type FrameForwardTable struct {
	entries []ForwardEntry
}

func NewSwitchTable() *FrameForwardTable {
	return &FrameForwardTable{
		entries: make([]ForwardEntry, 0),
	}
}

func (t *FrameForwardTable) Update(incoming Id, addr types.HwAddr) {
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

func (t *FrameForwardTable) LookupById(incoming Id) (ForwardEntry, bool) {
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
	Forward(incoming Id, frame na.Frame) error
}

type Switch struct {
	PortList map[Id]Port
	Table    *FrameForwardTable
	frmDec   *FrameDecoder
	frmEnc   *FrameEncoder
	quit     chan struct{}
}

func NewSwitch() *Switch {
	return &Switch{
		PortList: make(map[Id]Port),
		Table:    NewSwitchTable(),
		frmDec:   NewFrameDecoder(),
		frmEnc:   NewFrameEncoder(),
		quit:     make(chan struct{}),
	}
}

func NewSwitchWithTable(table *FrameForwardTable) *Switch {
	return &Switch{
		PortList: make(map[Id]Port),
		Table:    table,
		quit:     make(chan struct{}),
	}
}

func (s *Switch) Attach(port Port) error {
	if !port.Registered() {
		return fmt.Errorf("port is not registered")
	}
	if _, ok := s.PortList[port.Id()]; ok {
		return fmt.Errorf("already exist interface: %s", port.Id())
	}
	s.PortList[port.Id()] = port
	return nil
}

// Forward receives address of interface it receives frame, address of sender
// and frame to send to receiver. Based on id, address it determines whether to
// broadcast frame or forward it to others, otherwise just discard frame.
func (s *Switch) Forward(incoming Id, frame na.Frame) error {
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
	p := s.PortList[entry.Incoming]
	if err := p.Transmit(frm); err != nil {
		return err
	}
	return nil
}

// broadcastExcept sends frame to other interfaces except the interface
// with the id value given by parameter
func (s *Switch) broadcastExcept(incoming Id, frm []byte) error {
	err := errors.Multiple()
	for id, itf := range s.PortList {
		if !id.Equal(incoming) {
			err = err.Happen(itf.Transmit(frm))
		}
	}
	return err.Return()
}

func (s *Switch) Shutdown() {
	s.quit <- struct{}{}
}
