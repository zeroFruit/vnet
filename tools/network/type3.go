package network

import (
	"github.com/zeroFruit/vnet/pkg/link"
)

func Type3() (host1 *link.Host, host2 *link.Host, host3 *link.Host,
	swch1 *link.Switch, swch2 *link.Switch) {
	// setup node
	host1 = link.NewHost()
	host2 = link.NewHost()
	host3 = link.NewHost()
	swch1 = link.NewSwitch()
	swch2 = link.NewSwitch()

	// setup interface
	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), host1)
	attachInterface(host1, intf1)

	intf2 := link.NewInterface(40002, link.AddrFromStr("22-22-22-22-22-22"), host2)
	attachInterface(host2, intf2)

	intf3 := link.NewInterface(40003, link.AddrFromStr("33-33-33-33-33-33"), host3)
	attachInterface(host3, intf3)

	sp11 := link.NewSwitchPort(40004, swch1)
	sp12 := link.NewSwitchPort(40005, swch1)
	sp13 := link.NewSwitchPort(40006, swch1)
	attachSwchInterface(swch1, sp11, "1")
	attachSwchInterface(swch1, sp12, "2")
	attachSwchInterface(swch1, sp13, "3")

	sp21 := link.NewSwitchPort(40007, swch2)
	sp22 := link.NewSwitchPort(40008, swch2)
	attachSwchInterface(swch2, sp21, "1")
	attachSwchInterface(swch2, sp22, "2")

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(sp11, link1)

	link2 := link.NewLink(1)
	attachLink(intf2, link2)
	attachLink(sp12, link2)

	link3 := link.NewLink(1)
	attachLink(sp13, link3)
	attachLink(sp21, link3)

	link4 := link.NewLink(1)
	attachLink(sp22, link4)
	attachLink(intf3, link4)
	return
}
