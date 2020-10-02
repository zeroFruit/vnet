package phy_test

import (
	"sync"
	"testing"
	"time"

	"github.com/zeroFruit/vnet/phy"
)

/*
                        +----------+
                    0/4 |          |0/0
       +----------------+   R0_re  +---------------------------+
       |     40.1.1.1/24| 122.1.1.0|20.1.1.1/24                |
       |                +----------+                           |
       |                                                       |
       |                                                       |
       |                                                       |
       |40.1.1.2/24                                            |20.1.1.2/24
       |0/5                                                    |0/1
   +---+---+                                              +----+-----+
   |       |0/3                                        0/2|          |
   | R2_re +----------------------------------------------+    R1_re |
   |       |30.1.1.2/24                        30.1.1.1/24|          |
   +-------+                                              +----------+
*/

func TestNetworkTopology(t *testing.T) {
	// setup node
	node1 := phy.NewNode("R0_re")
	node2 := phy.NewNode("R1_re")
	node3 := phy.NewNode("R2_re")

	// setup interface
	intf1 := phy.NewInterface("0/0", "20.1.1.1", 40001, "11-11-11-11-11-11", node1.DataSink())
	intf2 := phy.NewInterface("0/4", "40.1.1.1", 40002, "11-11-11-11-11-12", node1.DataSink())
	attachInterface(t, node1, intf1)
	attachInterface(t, node1, intf2)

	intf3 := phy.NewInterface("0/1", "20.1.1.2", 40003, "11-11-11-11-11-13", node2.DataSink())
	intf4 := phy.NewInterface("0/2", "30.1.1.1", 40004, "11-11-11-11-11-14", node2.DataSink())
	attachInterface(t, node2, intf3)
	attachInterface(t, node2, intf4)

	intf5 := phy.NewInterface("0/3", "30.1.1.2", 40005, "11-11-11-11-11-15", node3.DataSink())
	intf6 := phy.NewInterface("0/5", "40.1.1.2", 40006, "11-11-11-11-11-16", node3.DataSink())
	attachInterface(t, node3, intf5)
	attachInterface(t, node3, intf6)

	// setup link
	link1 := phy.NewLink("Link1", 1)
	attachLink(t, intf1, link1)
	attachLink(t, intf3, link1)

	link2 := phy.NewLink("Link2", 2)
	attachLink(t, intf4, link2)
	attachLink(t, intf5, link2)

	link3 := phy.NewLink("Link3", 3)
	attachLink(t, intf6, link3)
	attachLink(t, intf2, link3)

	testNode(t, node1, "0/0", "Link1", 1)
	testNode(t, node1, "0/4", "Link3", 3)

	testNode(t, node2, "0/1", "Link1", 1)
	testNode(t, node2, "0/2", "Link2", 2)

	testNode(t, node3, "0/3", "Link2", 2)
	testNode(t, node3, "0/5", "Link3", 3)
}

func TestNodeSendReceive(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	sender := phy.NewNode("R0_re")
	receiver := phy.NewNode("R1_re")

	defer func() {
		receiver.Shutdown()
		sender.Shutdown()
	}()

	intf1 := phy.NewInterface("0/0", "20.1.1.1", 40001, "11-11-11-11-11-11", sender.DataSink())
	attachInterface(t, sender, intf1)

	intf2 := phy.NewInterface("0/1", "20.1.1.2", 40002, "11-11-11-11-11-13", receiver.DataSink())
	attachInterface(t, receiver, intf2)

	link1 := phy.NewLink("Link1", 1)
	attachLink(t, intf1, link1)
	attachLink(t, intf2, link1)

	go func() {
		time.Sleep(time.Millisecond * 300)
		if err := sender.SendTo("0/0", []byte{1, 2, 3}); err != nil {
			t.Fatalf("failed to send message: %v", err)
			t.Fail()
		}
		wg.Done()
	}()
	wg.Wait()
	// TODO: after changing `handleData` function, test received data by injecting mock object
	time.Sleep(time.Second * 1)
}

func attachInterface(t *testing.T, node *phy.Node, intf *phy.Interface) {
	if err := node.AttachInterface(intf); err != nil {
		t.Fatalf("failed to attach interface: %v", err)
	}
}

func attachLink(t *testing.T, intf *phy.Interface, link *phy.Link) {
	if err := intf.AttachLink(link); err != nil {
		t.Fatalf("failed to attach link: %v", err)
	}
}

func testNode(t *testing.T, node *phy.Node, intfId, linkId string, cost uint) {
	intf1_, err := node.GetInterfaceById(phy.IdOf(intfId))
	if err != nil {
		t.Fatalf("%s interface not exist on %s", intfId, node.GetId().Get())
	}
	link1_ := intf1_.GetLink()
	if link1_ == nil {
		t.Fatalf("link %s not exist on %s interface", linkId, intfId)
	}
	if link1_.GetId() != phy.IdOf(linkId) {
		t.Fatalf("link %s invalid id", linkId)
	}
	if link1_.GetCost() != cost {
		t.Fatalf("link %s expected cost is %d, but got %d", linkId, cost, link1_.GetCost())
	}
}
