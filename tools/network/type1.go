package network

import (
	"fmt"

	"github.com/zeroFruit/vnet/test"

	"github.com/zeroFruit/vnet/pkg/link"
)

func attachInterface(node *link.Host, itf link.Interface) {
	node.AttachInterface(itf)
}

func attachSwchInterface(swch *link.Switch, itf link.Interface) {
	test.ShouldSuccess(func() error {
		return swch.Attach(itf)
	})
}

func attachLink(itf link.Interface, link *link.Link) {
	if err := itf.AttachLink(link); err != nil {
		panic(fmt.Sprintf("failed to attach link: %v", err))
	}
}

func Type1() (node1 *link.Host, node2 *link.Host) {
	// setup node
	node1 = link.NewNode()
	node2 = link.NewNode()

	// setup interface
	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), node1)
	attachInterface(node1, intf1)

	intf2 := link.NewInterface(40002, link.AddrFromStr("11-11-11-11-11-12"), node2)
	attachInterface(node2, intf2)

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(intf2, link1)
	return
}
