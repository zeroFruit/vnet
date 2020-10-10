package main

import (
	"fmt"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"

	"github.com/zeroFruit/vnet/pkg/arp"
	"github.com/zeroFruit/vnet/pkg/net"
	"github.com/zeroFruit/vnet/tools/network"
)

/*

        1.1.1.1                            1.1.1.2
   11-11-11-11-11-11                  11-11-11-11-11-12
       +-------+                          +-------+
       | node1 ---------------------------- node2 |
       |       |                          |       |
       +-------+                          +-------+
*/
func main() {
	node1, node2 := network.Type1()

	net1 := net.NewNode(node1)
	net2 := net.NewNode(node2)

	node1.RegisterNetHandler(net1)
	node2.RegisterNetHandler(net2)

	for _, fn := range []func() error{
		func() error {
			return net1.UpdateAddr(link.AddrFromStr("11-11-11-11-11-11"), net.AddrFromStr("1.1.1.1"))
		},
		func() error {
			return net2.UpdateAddr(link.AddrFromStr("11-11-11-11-11-12"), net.AddrFromStr("1.1.1.2"))
		},
	} {
		ShouldSuccess(fn)
	}

	table := arp.NewTable()
	arp1 := arp.NewWithTable(net1, net.NewArpPayloadEncoder(), table)
	arp2 := arp.New(net2, net.NewArpPayloadEncoder())

	net1.RegisterArp(arp1)
	net2.RegisterArp(arp2)

	errs := arp1.Broadcast(net.AddrFromStr("1.1.1.2"))
	if len(errs) != 0 {
		Fatalf("expected errors when broadcast is 0 but got %d", len(errs))
	}

	time.Sleep(1 * time.Second)

	entry, ok := table.Lookup(arp.Key{NetAddr: net.AddrFromStr("1.1.1.2")})
	if !ok {
		Fatalf("net address not exist on table")
	}
	if !entry.HwAddr.Equal(link.AddrFromStr("11-11-11-11-11-12")) {
		Fatalf("expected hw address is '11-11-11-11-11-12', but got '%s'", entry.HwAddr)
	}
}

func Fatalf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func ShouldSuccess(fn func() error) {
	if err := fn(); err != nil {
		Fatalf("fn should success but failed with err: %v", err)
	}
}
