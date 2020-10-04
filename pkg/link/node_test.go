package link_test

import (
	"sync"
	"testing"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"
)

func TestNetworkTopology(t *testing.T) {
	// setup node
	node1 := link.NewNode()
	node2 := link.NewNode()
	node3 := link.NewNode()

	// setup interface
	intf1 := link.NewInterface(40001, "11-11-11-11-11-11", node1)
	intf2 := link.NewInterface(40002, "11-11-11-11-11-12", node1)
	attachInterface(t, node1, intf1)
	attachInterface(t, node1, intf2)

	intf3 := link.NewInterface(40003, "11-11-11-11-11-13", node2)
	intf4 := link.NewInterface(40004, "11-11-11-11-11-14", node2)
	attachInterface(t, node2, intf3)
	attachInterface(t, node2, intf4)

	intf5 := link.NewInterface(40005, "11-11-11-11-11-15", node3)
	intf6 := link.NewInterface(40006, "11-11-11-11-11-16", node3)
	attachInterface(t, node3, intf5)
	attachInterface(t, node3, intf6)

	// setup link
	link1 := link.NewLink(1)
	attachLink(t, intf1, link1)
	attachLink(t, intf3, link1)

	link2 := link.NewLink(2)
	attachLink(t, intf4, link2)
	attachLink(t, intf5, link2)

	link3 := link.NewLink(3)
	attachLink(t, intf6, link3)
	attachLink(t, intf2, link3)

	testNode(t, node1, "11-11-11-11-11-11", 1)
	testNode(t, node1, "11-11-11-11-11-12", 3)

	testNode(t, node2, "11-11-11-11-11-13", 1)
	testNode(t, node2, "11-11-11-11-11-14", 2)

	testNode(t, node3, "11-11-11-11-11-15", 2)
	testNode(t, node3, "11-11-11-11-11-16", 3)
}

// TODO: close UDP connection
// TODO: this test should moved to `test/` package
func TestNodeSendReceive(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	sender := link.NewNode()
	receiver := link.NewNode()

	defer func() {
		receiver.Shutdown()
		sender.Shutdown()
	}()

	intf1 := link.NewInterface(40001, "11-11-11-11-11-11", sender)
	attachInterface(t, sender, intf1)

	intf2 := link.NewInterface(40002, "11-11-11-11-11-13", receiver)
	attachInterface(t, receiver, intf2)

	link1 := link.NewLink(1)
	attachLink(t, intf1, link1)
	attachLink(t, intf2, link1)

	go func() {
		time.Sleep(time.Millisecond * 300)
		if err := sender.Send("11-11-11-11-11-11", []byte{1, 2, 3}); err != nil {
			t.Fatalf("failed to send message: %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
	// TODO: after changing `handleData` function, test received data by injecting mock object
	time.Sleep(time.Second * 1)
}

func attachInterface(t *testing.T, node *link.Node, itf link.Interface) {
	node.AttachInterface(itf)
}

func attachLink(t *testing.T, itf link.Interface, link *link.Link) {
	if err := itf.AttachLink(link); err != nil {
		t.Fatalf("failed to attach link: %v", err)
	}
}

func testNode(t *testing.T, node *link.Node, addr link.Addr, cost uint) {
	intf1_, err := node.InterfaceOfAddr(addr)
	if err != nil {
		t.Fatalf("interface not exist with address: %s", addr)
	}
	link1_ := intf1_.GetLink()
	if link1_ == nil {
		t.Fatalf("link not exist on interface with address: %s", intf1_.Address())
	}
	otherLink, err := link1_.GetOtherInterface(addr)
	if err != nil {
		t.Fatalf("otherLink not exist: %v", err)
	}
	originLink, err := link1_.GetOtherInterface(otherLink.Address())
	if err != nil {
		t.Fatalf("link1_ not exist: %v", err)
	}
	if !originLink.Address().Equal(intf1_.Address()) {
		t.Fatalf("address is not equal")
	}
	if link1_.GetCost() != cost {
		t.Fatalf("link with address %s expected cost is %d, but got %d", addr, cost, link1_.GetCost())
	}
}
