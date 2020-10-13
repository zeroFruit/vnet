package main

import (
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/link/na"
	"github.com/zeroFruit/vnet/test"
	"github.com/zeroFruit/vnet/tools/network"
	"os"
	"sync"
	"time"
)

/*

   11-11-11-11-11-11                  11-11-11-11-11-12
       +-------+                          +-------+
       | node1 ---------------------------- node2 |
       |       |                          |       |
       +-------+                          +-------+
*/

type mockNetFrameHandler struct {
	handleFunc func(frame na.Frame)
}

func (h *mockNetFrameHandler) Handle(frame na.Frame) {
	h.handleFunc(frame)
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	node1, node2 := network.Type1()
	node2.RegisterNetHandler(&mockNetFrameHandler{
		handleFunc: func(frame na.Frame) {
			if !frame.Src.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
				test.Fatalf("expected src address is '11-11-11-11-11-11', but got '%s'", frame.Src)
			}
			if !frame.Dest.Equal(link.AddrFromStr("11-11-11-11-11-12")) {
				test.Fatalf("expected src address is '11-11-11-11-11-12', but got '%s'", frame.Dest)
			}
			if string(frame.Payload) != "data" {
				test.Fatalf("expected payload is 'data', but got %s", string(frame.Payload))
			}
			wg.Done()
		},
	})
	if err := node1.Send(link.AddrFromStr("11-11-11-11-11-12"), []byte("data")); err != nil {
		test.Fatalf("failed to send payload: %v", err)
	}
	if test.WaitTimeout(wg, 1 * time.Second) {
		os.Exit(1)
	}
}
