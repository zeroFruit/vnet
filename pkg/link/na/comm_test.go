package na_test

import (
	"sync"
	"testing"
	"time"

	"github.com/zeroFruit/vnet/pkg/link/na"
)

// TestDatagramTransport tests whether Card could send bytes data to others properly.
// But it listens UDP port internally so skip when unit-testing.
func TestDatagramTransport(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}
	wg.Add(1)
	sender, err := na.New("127.0.0.1", 40000)
	if err != nil {
		t.Fatalf("failed to create network adapter: %v", err)
	}
	receiver, err := na.New("127.0.0.1", 40001)
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
		wg.Done()
	}()

	time.Sleep(time.Millisecond * 300)
	sender.Send([]byte{'a'}, "127.0.0.1:40001")
	wg.Wait()
}
