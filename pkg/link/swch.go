package link

import (
	"fmt"
	"github.com/zeroFruit/vnet/pkg/errors"
	"github.com/zeroFruit/vnet/pkg/link/na"
	"log"
	"time"

	"github.com/zeroFruit/vnet/pkg/types"
)

type SwitchEntry struct {
	Id   string
	Addr types.HwAddr
	Time time.Time
}

type SwitchTable struct {
	entries []SwitchEntry
}

func NewSwitchTable() *SwitchTable {
	return &SwitchTable{
		entries: make([]SwitchEntry, 0),
	}
}

func (t *SwitchTable) Update(id string, addr types.HwAddr) {
	idxToRemove := -1
	for i, e := range t.entries {
		if e.Addr.Equal(addr) {
			idxToRemove = i
		}
	}
	if idxToRemove != -1 {
		t.entries = append(t.entries[:idxToRemove], t.entries[idxToRemove+1:]...)
	}
	t.entries = append(t.entries, SwitchEntry{Id: id, Addr: addr, Time: time.Now()})
}

func (t *SwitchTable) LookupByAddr(key types.HwAddr) (SwitchEntry, bool) {
	for _, e := range t.entries {
		if e.Addr.Equal(key) {
			return e, true
		}
	}
	return SwitchEntry{}, false
}

func (t *SwitchTable) LookupById(id string) (SwitchEntry, bool) {
	for _, e := range t.entries {
		if e.Id == id {
			return e, true
		}
	}
	return SwitchEntry{}, false
}

func (t *SwitchTable) Entries() []SwitchEntry {
	return t.entries
}

type PacketForwarder interface {
	Forward(id string, src types.HwAddr, dest types.HwAddr, pkt []byte) error
}

type AnonymInterface interface {
	GetLink() *Link
	AttachLink(link *Link) error
	Send(frame []byte) error
}

type UDPBasedSwitchInterface struct {
	*UDPBasedInterface
	id        string
	forwarder PacketForwarder
	frmDec    *FrameDecoder
}

func NewSwitchInterface(port int, hwAddr types.HwAddr, id string, forwarder PacketForwarder) *UDPBasedSwitchInterface {
	si := &UDPBasedSwitchInterface{
		id:        id,
		forwarder: forwarder,
		frmDec:    NewFrameDecoder(),
	}
	i := NewInterface(port, hwAddr, si.handler)
	si.UDPBasedInterface = i
	return si
}

// handle receives datagram from UDPBasedInterface than forward it to PacketForwarder
func (si *UDPBasedSwitchInterface) handle(fd *na.FrameData) error {
	frame, err := si.frmDec.Decode(fd.Buf)
	if err != nil {
		return err
	}
	return si.forwarder.Forward(si.id, frame.Src, frame.Dest, frame.Payload)
}

type Switch struct {
	ItfList map[string]AnonymInterface
	table   *SwitchTable
	quit    chan struct{}
}

func NewSwitch() *Switch {
	return &Switch{
		ItfList: make(map[string]AnonymInterface),
		table:   NewSwitchTable(),
		quit:    make(chan struct{}),
	}
}

func NewSwitchWithTable(table *SwitchTable) *Switch {
	return &Switch{
		ItfList: make(map[string]AnonymInterface),
		table:   table,
		quit:    make(chan struct{}),
	}
}

func (s *Switch) Attach(id string, itf AnonymInterface) error {
	if _, ok := s.ItfList[id]; ok {
		return fmt.Errorf("already exist interface: %s", id)
	}
	s.ItfList[id] = itf
	return nil
}

// Forward receives id of interface it receives packet, address of sender
// and packet to send to receiver. Based on id, address it determines whether to
// broadcast packet or forward it to others, otherwise just discard packet.
func (s *Switch) Forward(id string, src types.HwAddr, dest types.HwAddr, pkt []byte) error {
	s.table.Update(id, src)

	entry, ok := s.table.LookupByAddr(dest)
	if !ok {
		return s.broadcastExcept(id, pkt)
	}
	if entry.Id == id {
		log.Printf("discard packet from id: %s, src: %s, dest: %s\n", id, src, dest)
		return nil
	}
	return s.ItfList[entry.Id].Send(pkt)
}

// broadcastExcept sends packet to other interfaces except the interface
// with the id value given by parameter
func (s *Switch) broadcastExcept(id string, pkt []byte) error {
	err := errors.Multiple()
	for itfId, itf := range s.ItfList {
		if itfId != id {
			err = err.Happen(itf.Send(pkt))
		}
	}
	return err.Return()
}

func (s *Switch) Shutdown() {
	s.quit <- struct{}{}
}
