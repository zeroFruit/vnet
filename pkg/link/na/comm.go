package na

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/zeroFruit/vnet/pkg/link/internal"
)

const (
	// udpPacketBufSize is used to buffer incoming packets during read
	// operations.
	udpPacketBufSize = 65536

	// udpRecvBufSize is a large buffer size that we attempt to set UDP
	// sockets to in order to handle a large volume of messages.
	udpRecvBufSize = 2 * 1024 * 1024
)

type packetConn interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
}

// Card is abstracted network adapter part to simulate bytes transport on
// physical cable. Node's interface uses this interface to send frame
// between nodes.
type Card interface {
	Send(buf []byte, addr string) (time.Time, error)
	Recv() <-chan *FrameData
}

type UDPTransport struct {
	ip         internal.Addr
	port       int
	frameCh    chan *FrameData
	packetConn packetConn
	shutdown   bool
	lock       sync.RWMutex
}

func ListenDatagram(ip internal.Addr, port int) (*net.UDPConn, error) {
	udpAddr := &net.UDPAddr{
		// TODO: remove ip parameters and set this value by internal.DefaultAddr value
		IP:   net.ParseIP(string(ip)),
		Port: port,
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start UDP listener on %s:%d: %v", ip, port, err)
	}
	if err := setUDPRecvBuf(udpConn); err != nil {
		return nil, fmt.Errorf("failed to resize UDP buffer: %v", err)
	}
	return udpConn, nil
}

func New(ip internal.Addr, port int) (*UDPTransport, error) {
	pc, err := ListenDatagram(ip, port)
	if err != nil {
		return nil, err
	}
	t := UDPTransport{
		ip:         ip,
		port:       port,
		frameCh:    make(chan *FrameData),
		packetConn: pc,
		shutdown:   false,
		lock:       sync.RWMutex{},
	}
	go t.listen()
	return &t, nil
}

func (p *UDPTransport) listen() {
	for {
		buf := make([]byte, udpPacketBufSize)
		n, _, err := p.packetConn.ReadFrom(buf)
		now := time.Now()
		if p.shutdown {
			break
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("error reading UDP packet: %v", err))
			continue
		}
		if n < 1 {
			fmt.Println(fmt.Sprintf("error UDP packet too short, got %d bytes", n))
			continue
		}
		p.frameCh <- &FrameData{
			Buf:       buf[:n],
			Timestamp: now,
		}
	}
}

func (p *UDPTransport) Recv() <-chan *FrameData {
	return p.frameCh
}

func (p *UDPTransport) Send(buf []byte, addr string) (time.Time, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return time.Time{}, err
	}
	if _, err := p.packetConn.WriteTo(buf, udpAddr); err != nil {
		return time.Time{}, err
	}
	return time.Now(), nil
}

// setUDPRecvBuf is used to resize the UDP receive window. The function
// attempts to set the read buffer to `udpRecvBuf` but backs off until
// the read buffer can be set.
func setUDPRecvBuf(c *net.UDPConn) error {
	size := udpRecvBufSize
	var err error
	for size > 0 {
		if err = c.SetReadBuffer(size); err == nil {
			return nil
		}
		size = size / 2
	}
	return err
}
