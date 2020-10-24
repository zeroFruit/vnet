package main

import (
	"os"
	"sync"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/test"
	"github.com/zeroFruit/vnet/tools/network"
)

type mockNetHandler struct {
	handleFunc func(pl []byte)
}

func (h *mockNetHandler) Handle(pl []byte) {
	h.handleFunc(pl)
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	host1, host2 := network.Type1()
	host2.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			if string(pl) != "hello" {
				test.Fatalf("expected data is 'data', but got '%s'", string(pl))
			}
			wg.Done()
		},
	})
	if err := host1.Send(link.AddrFromStr("11-11-11-11-11-12"), []byte("hello")); err != nil {
		test.Fatalf("failed to send payload: %v", err)
	}
	if test.WaitTimeout(wg, 1*time.Second) {
		os.Exit(1)
	}
}
