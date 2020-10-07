package link_test

import (
	"sync"
	"testing"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"
)

func TestDatagramTransport(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}
	wg.Add(1)
	sender, err := link.NewNetworkAdapter("127.0.0.1", 40000)
	if err != nil {
		t.Fatalf("failed to create network adapter: %v", err)
	}
	receiver, err := link.NewNetworkAdapter("127.0.0.1", 40001)
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
