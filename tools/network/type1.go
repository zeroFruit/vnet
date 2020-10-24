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

func Type1() (host1 *link.Host, host2 *link.Host) {
	// setup node
	host1 = link.NewHost()
	host2 = link.NewHost()

	// setup interface
	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), host1)
	attachInterface(host1, intf1)

	intf2 := link.NewInterface(40002, link.AddrFromStr("11-11-11-11-11-12"), host2)
	attachInterface(host2, intf2)

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(intf2, link1)
	return
}
