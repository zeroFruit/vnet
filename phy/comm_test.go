<<<<<<< HEAD:pkg/link/comm_test.go
package link_test
=======
package phy_test
>>>>>>> 18e4a4c... feat: add arp table, separate net address:phy/comm_test.go

import (
	"sync"
	"testing"
	"time"

<<<<<<< HEAD:pkg/link/comm_test.go
	"github.com/zeroFruit/vnet/pkg/link"
=======
	"github.com/zeroFruit/vnet/phy"
>>>>>>> 18e4a4c... feat: add arp table, separate net address:phy/comm_test.go
)

func TestDatagramTransport(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}
	wg.Add(1)
<<<<<<< HEAD:pkg/link/comm_test.go
	sender, err := link.NewNetworkAdapter("127.0.0.1", 40000)
	if err != nil {
		t.Fatalf("failed to create network adapter: %v", err)
	}
	receiver, err := link.NewNetworkAdapter("127.0.0.1", 40001)
=======
	sender, err := phy.NewNetworkAdapter("127.0.0.1", 40000)
	if err != nil {
		t.Fatalf("failed to create network adapter: %v", err)
	}
	receiver, err := phy.NewNetworkAdapter("127.0.0.1", 40001)
>>>>>>> 18e4a4c... feat: add arp table, separate net address:phy/comm_test.go
	if err != nil {
		t.Fatalf("failed to create network adapter: %v", err)
	}
	go func() {
		data := <-receiver.Recv()
		if len(data.Buf) != 1 {
			t.Fatalf("expected data buf length is 1, but got %d", len(data.Buf))
		}
		if data.Buf[0] != 'a' {
			t.Fatalf("expected datagram is 'a', but got %c", data.Buf[0])
		}
		if data.From != "" {
			t.Fatalf("expected sender address is '127.0.0.1:40000', but got %s", data.From)
		}
		if data.HardwareAddr != "" {
			t.Fatalf("hardware address is not empty")
		}
		wg.Done()
	}()

	time.Sleep(time.Millisecond * 300)
	sender.Send([]byte{'a'}, "127.0.0.1:40001")
	wg.Wait()
}
