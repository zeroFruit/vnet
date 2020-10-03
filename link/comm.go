<<<<<<< HEAD:phy/comm.go
<<<<<<< HEAD:pkg/link/comm.go
package link
=======
package phy
>>>>>>> 18e4a4c... feat: add arp table, separate net address:phy/comm.go
=======
package link
>>>>>>> 6a51a9b... feat: remove phy package and migrate to link:link/comm.go

import (
	"fmt"
	"net"
	"sync"
	"time"

<<<<<<< HEAD:phy/comm.go
<<<<<<< HEAD:pkg/link/comm.go
	"github.com/zeroFruit/vnet/pkg/link/internal"
=======
	"github.com/zeroFruit/vnet/phy/internal"
>>>>>>> 18e4a4c... feat: add arp table, separate net address:phy/comm.go
=======
	"github.com/zeroFruit/vnet/link/internal"
>>>>>>> 6a51a9b... feat: remove phy package and migrate to link:link/comm.go
)

const (
	// udpPacketBufSize is used to buffer incoming packets during read
	// operations.
	udpPacketBufSize = 65536

	// udpRecvBufSize is a large buffer size that we attempt to set UDP
	// sockets to in order to handle a large volume of messages.
	udpRecvBufSize = 2 * 1024 * 1024
)

type Datagram struct {
	Buf          []byte
	From         string
	HardwareAddr Addr
	Timestamp    time.Time
}

type packetConn interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
}

type NetworkAdapter interface {
	Send(buf []byte, addr string) (time.Time, error)
	Recv() <-chan *Datagram
}

type UDPTransport struct {
	ip         internal.Addr
	port       int
	packetCh   chan *Datagram
	packetConn packetConn
	shutdown   bool
	lock       sync.RWMutex
}

func ListenDatagram(ip internal.Addr, port int) (*net.UDPConn, error) {
	udpAddr := &net.UDPAddr{
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

func NewNetworkAdapter(ip internal.Addr, port int) (*UDPTransport, error) {
	pc, err := ListenDatagram(ip, port)
	if err != nil {
		return nil, err
	}
	t := UDPTransport{
		ip:         ip,
		port:       port,
		packetCh:   make(chan *Datagram),
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
		p.packetCh <- &Datagram{
			Buf:       buf[:n],
			Timestamp: now,
		}
	}
}

func (p *UDPTransport) Recv() <-chan *Datagram {
	return p.packetCh
}

func (p *UDPTransport) Send(buf []byte, addr string) (time.Time, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", string(addr))
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
