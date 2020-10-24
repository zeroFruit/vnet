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

	sintf11 := link.NewInterface(40004, link.AddrFromStr("00-00-00-00-00-01"), swch1)
	sintf12 := link.NewInterface(40005, link.AddrFromStr("00-00-00-00-00-02"), swch1)
	sintf13 := link.NewInterface(40006, link.AddrFromStr("00-00-00-00-00-03"), swch1)
	attachSwchInterface(swch1, sintf11)
	attachSwchInterface(swch1, sintf12)
	attachSwchInterface(swch1, sintf13)

	sintf21 := link.NewInterface(40007, link.AddrFromStr("00-00-00-00-00-11"), swch2)
	sintf22 := link.NewInterface(40008, link.AddrFromStr("00-00-00-00-00-12"), swch2)
	attachSwchInterface(swch2, sintf21)
	attachSwchInterface(swch2, sintf22)

	// setup link
	link1 := link.NewLink(1)
	attachLink(intf1, link1)
	attachLink(sintf11, link1)

	link2 := link.NewLink(1)
	attachLink(intf2, link2)
	attachLink(sintf12, link2)

	link3 := link.NewLink(1)
	attachLink(sintf13, link3)
	attachLink(sintf21, link3)

	link4 := link.NewLink(1)
	attachLink(sintf22, link4)
	attachLink(intf3, link4)
	return
}
