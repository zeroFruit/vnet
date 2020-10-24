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
	wg.Add(1)

	node1, node2, node3, swch := network.Type2()

	// pre-define switch table
	swch.Table.Update(link.AddrFromStr("00-00-00-00-00-02"), link.AddrFromStr("22-22-22-22-22-22"))

	// setup mock network layer handler
	node1.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			test.Fatalf("this should not be called")
		},
	})
	node2.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			fmt.Println("node2:", string(pl))
			wg.Done()
		},
	})
	node3.RegisterNetHandler(&test.MockNetHandler{
		HandleFunc: func(pl []byte) {
			test.Fatalf("this should not be called")
		},
	})
	if err := node1.Send(link.AddrFromStr("22-22-22-22-22-22"), []byte("hello")); err != nil {
		test.Fatalf("failed to send frame from node1: %v", err)
	}
	if test.WaitTimeout(wg, 1*time.Second) {
		os.Exit(1)
	}
}
