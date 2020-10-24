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

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	node1, node2, node3, _ := network.Type2()
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
			fmt.Println("node3:", string(pl))
			wg.Done()
		},
	})
	if err := node1.Send(link.AddrFromStr("22-22-22-22-22-22"), []byte("hello")); err != nil {
		test.Fatalf("failed to send frame from node1: %v", err)
	}
	if test.WaitTimeout(wg, 1*time.Second) {
		os.Exit(1)
	}
}
