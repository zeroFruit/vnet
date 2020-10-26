package link

import (
	"errors"
	"log"
	"strconv"

	"github.com/zeroFruit/vnet/pkg/types"

	"github.com/zeroFruit/vnet/pkg/link/internal"
	"github.com/zeroFruit/vnet/pkg/link/na"
)

type Transmitter interface {
	Transmit(frame []byte) error
}

type Interface interface {
	Transmitter
	EndPoint
	Address() types.HwAddr
}

// FrameDataHandler receives frame data comes from network adapter. Each node should have
// FrameDataHandler to process and validate frame.
type FrameDataHandler interface {
	handle(frame *na.FrameData) error
}

type NetworkAdapter struct {
	internalIP   internal.Addr
	internalPort int
	Addr         types.HwAddr
	link         *Link
	adapter      na.Card
	handler      FrameDataHandler
	quit         chan struct{}
}

func NewInterface(port int, hwAddr types.HwAddr, handler FrameDataHandler) *NetworkAdapter {
	itf := &NetworkAdapter{
		internalPort: port,
		internalIP:   internal.DefaultAddr,
		Addr:         hwAddr,
		handler:      handler,
		quit:         make(chan struct{}),
	}
	return itf
}

func (i *NetworkAdapter) GetLink() *Link {
	return i.link
}

func (i *NetworkAdapter) AttachLink(link *Link) error {
	if i.link != nil {
		return errors.New("link already exist")
	}
	if err := link.AttachEndpoint(i); err != nil {
		return err
	}
	adapter, err := na.New(i.internalIP, i.internalPort)
	if err != nil {
		return err
	}
	i.adapter = adapter
	i.link = link
	go i.sink()
	return nil
}

func (i *NetworkAdapter) Transmit(frame []byte) error {
	receiver, err := i.link.Opposite(i.Id())
	if err != nil {
		return err
	}
	i.adapter.Send(frame, receiver.InternalAddress().String())
	return nil
}

func (i *NetworkAdapter) sink() {
	for {
		select {
		case data := <-i.adapter.Recv():
			data.Incoming = i.Addr
			if err := i.handler.handle(data); err != nil {
				log.Fatal(err)
			}
		case <-i.quit:
			return
		}
	}
}

func (i *NetworkAdapter) Address() types.HwAddr {
	return i.Addr
}

func (i *NetworkAdapter) Id() Id {
	return Id(i.Addr.String())
}

func (i *NetworkAdapter) InternalAddress() Addr {
	return Addr(string(i.internalIP) + ":" + strconv.Itoa(i.internalPort))
}

func (i *NetworkAdapter) shutdown() {
	i.quit <- struct{}{}
}
