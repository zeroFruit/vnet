package network

import (
	"fmt"

	"github.com/zeroFruit/vnet/link"
)

func attachInterface(node *link.Node, itf link.Interface) {
	node.AttachInterface(itf)
}

func attachLink(itf link.Interface, link *link.Link) {
	if err := itf.AttachLink(link); err != nil {
		panic(fmt.Sprintf("failed to attach link: %v", err))
	}
}

func Type1() (node1 *link.Node, node2 *link.Node, node3 *link.Node) {
	// setup node
	node1 = link.NewNode()
	node2 = link.NewNode()
	node3 = link.NewNode()

	// setup interface
	intf1 := link.NewInterface(40001, "11-11-11-11-11-11", node1.DataSink())
	intf2 := link.NewInterface(40002, "11-11-11-11-11-12", node1.DataSink())
	attachInterface(node1, intf1)
	attachInterface(node1, intf2)

	intf3 := link.NewInterface(40003, "11-11-11-11-11-13", node2.DataSink())
	intf4 := link.NewInterface(40004, "11-11-11-11-11-14", node2.DataSink())
	attachInterface(node2, intf3)
	attachInterface(node2, intf4)

	intf5 := link.NewInterface(40005, "11-11-11-11-11-15", node3.DataSink())
	intf6 := link.NewInterface(40006, "11-11-11-11-11-16", node3.DataSink())
	attachInterface(node3, intf5)
	attachInterface(node3, intf6)

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(intf3, link1)

	link2 := link.NewLink(2)
	attachLink(intf4, link2)
	attachLink(intf5, link2)

	link3 := link.NewLink(3)
	attachLink(intf6, link3)
	attachLink(intf2, link3)
	return
}
