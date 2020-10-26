package link

import (
	"errors"
	"log"
	"strconv"

	"github.com/zeroFruit/vnet/pkg/link/internal"
	"github.com/zeroFruit/vnet/pkg/link/na"
)

// Port can transmit data and can be point of link. But it has no hardware
// address. Before using Port, it must register its own Id
type Port interface {
	Transmitter
	EndPoint
	Register(id Id)
	Registered() bool
}

type SwitchPort struct {
	internalIP   internal.Addr
	internalPort int
	id           Id
	link         *Link
	adapter      na.Card
	quit         chan struct{}
	frmDec       *FrameDecoder
	frmForwarder FrameForwarder
}

func NewSwitchPort(port int, frmForwarder FrameForwarder) *SwitchPort {
	sp := &SwitchPort{
		internalPort: port,
		internalIP:   internal.DefaultAddr,
		quit:         make(chan struct{}),
		frmDec:       NewFrameDecoder(),
		frmForwarder: frmForwarder,
	}
	return sp
}

func (s *SwitchPort) Register(id Id) {
	s.id = id
}
func (s *SwitchPort) Registered() bool {
	return s.id != ""
}

func (s *SwitchPort) AttachLink(link *Link) error {
	if !s.Registered() {
		return errors.New("id is not registered. before attach link, first register id on your port")
	}
	if s.link != nil {
		return errors.New("link already exist")
	}
	if err := link.AttachEndpoint(s); err != nil {
		return err
	}
	adapter, err := na.New(s.internalIP, s.internalPort)
	if err != nil {
		return err
	}
	s.adapter = adapter
	s.link = link
	go s.sink()
	return nil
}

func (s *SwitchPort) GetLink() *Link {
	return s.link
}

func (s *SwitchPort) Transmit(frame []byte) error {
	receiver, err := s.link.Opposite(s.Id())
	if err != nil {
		return err
	}
	s.adapter.Send(frame, receiver.InternalAddress().String())
	return nil
}

func (s *SwitchPort) sink() {
	for {
		select {
		case data := <-s.adapter.Recv():
			if err := s.handle(data); err != nil {
				log.Fatal(err)
			}
		case <-s.quit:
			return
		}
	}
}

func (s *SwitchPort) handle(fd *na.FrameData) error {
	frame, err := s.frmDec.Decode(fd.Buf)
	if err != nil {
		return err
	}
	return s.frmForwarder.Forward(s.id, frame)
}

func (s *SwitchPort) Id() Id {
	return s.id
}

func (s *SwitchPort) InternalAddress() Addr {
	return Addr(string(s.internalIP) + ":" + strconv.Itoa(s.internalPort))
}

func (s *SwitchPort) shutdown() {
	s.quit <- struct{}{}
}
