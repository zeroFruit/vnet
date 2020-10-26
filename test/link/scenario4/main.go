package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/test"
	"github.com/zeroFruit/vnet/tools/network"
)

type netHandler struct {
	handleFunc func(pl []byte)
}

func (h *netHandler) Handle(pl []byte) {
	h.handleFunc(pl)
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	host1, host2, host3, _, swch2 := network.Type3()

	// pre-define switch table
	swch2.Table.Update("1", link.AddrFromStr("22-22-22-22-22-22"))

	host1.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			test.Fatalf("this should not be called")
		},
	})
	host2.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			fmt.Println("node2:", string(pl))
			wg.Done()
		},
	})
	host3.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			test.Fatalf("this should not be called")
		},
	})
	if err := host1.Send(link.AddrFromStr("22-22-22-22-22-22"), []byte("hello")); err != nil {
		test.Fatalf("failed to send frame from node1: %v", err)
	}

	go func() {
		// wait one seconds for checking frame discard log
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	if test.WaitTimeout(wg, 2*time.Second) {
		os.Exit(1)
	}
}
