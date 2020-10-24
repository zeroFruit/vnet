package network

import (
	"github.com/zeroFruit/vnet/pkg/link"
)

func Type2() (host1 *link.Host, host2 *link.Host, host3 *link.Host, swch *link.Switch) {
	// setup node
	host1 = link.NewHost()
	host2 = link.NewHost()
	host3 = link.NewHost()
	swch = link.NewSwitch()

	// setup interface
	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), host1)
	attachInterface(host1, intf1)

	intf2 := link.NewInterface(40002, link.AddrFromStr("22-22-22-22-22-22"), host2)
	attachInterface(host2, intf2)

	intf3 := link.NewInterface(40003, link.AddrFromStr("33-33-33-33-33-33"), host3)
	attachInterface(host3, intf3)

	sintf1 := link.NewInterface(40004, link.AddrFromStr("00-00-00-00-00-01"), swch)
	sintf2 := link.NewInterface(40005, link.AddrFromStr("00-00-00-00-00-02"), swch)
	sintf3 := link.NewInterface(40006, link.AddrFromStr("00-00-00-00-00-03"), swch)
	attachSwchInterface(swch, sintf1)
	attachSwchInterface(swch, sintf2)
	attachSwchInterface(swch, sintf3)

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(sintf1, link1)

	link2 := link.NewLink(1)
	attachLink(intf2, link2)
	attachLink(sintf2, link2)

	link3 := link.NewLink(1)
	attachLink(intf3, link3)
	attachLink(sintf3, link3)
	return
}
