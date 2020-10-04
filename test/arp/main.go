package main

import (
	"fmt"
	"time"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/net"
	"github.com/zeroFruit/vnet/test/network"
)

/*
     +-----------------+
     |                 |
     |                 |11-11-11-11-11-11
     |      Node 1     ----------------------------------+
     |                 | NetAddr: 1.1.1.1                |
     |                 |                                 |
     +--------|--------+                                 |
              | 11-11-11-11-11-12                        |
              | NetAddr: 1.1.1.2                         | 11-11-11-11-11-13
              |                                          | NetAddr: 2.2.2.1
              |                                 +--------|--------+
              |                                 |                 |
              |                                 |                 |
              |                                 |      Node 2     |
              |                                 |                 |
              |                                 |                 |
              |                                 +--------|--------+
              |
              |      +-----------------+
              |      |                 |
              |      |                 |
              +-------     Node 3     --
   11-11-11-11-11-16 |                 |
	NetAddr: 3.3.3.1 |                 |
                     +-----------------+
*/
func main() {
	node1, node2, node3 := network.Type1()

	net1 := net.NewNode(node1)
	net2 := net.NewNode(node2)
	net3 := net.NewNode(node3)

	for _, fn := range [](func() error){
		func() error {
			return net1.UpdateAddr(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("1.1.1.1"))
		},
		func() error {
			return net1.UpdateAddr(link.AddrFromStr("11-11-11-11-11-12"), net.AddrFromStr("1.1.1.2"))
		},
		func() error {
			return net2.UpdateAddr(link.AddrFromStr("11-11-11-11-11-13"), net.AddrFromStr("2.2.2.1"))
		},
		func() error {
			return net3.UpdateAddr(link.AddrFromStr("11-11-11-11-11-16"), net.AddrFromStr("3.3.3.1"))
		},
	} {
		ShouldSuccess(fn)
	}

	svc1 := arp.New(arp.AdaptNode(net1))
	arp.New(arp.AdaptNode(net2))
	arp.New(arp.AdaptNode(net3))

	errs := svc1.Broadcast(net.AddrFromStr("2.2.2.1"))
	if len(errs) != 0 {
		Fatalf("expected errors when broadcast is 0 but got %d", len(errs))
	}
	time.Sleep(3 * time.Second)
}

func Fatalf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a))
}

func ShouldSuccess(fn func() error) {
	if err := fn(); err != nil {
		Fatalf("fn should success but failed with err: %v", err)
	}
}
