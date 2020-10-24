package link_test

import (
	"sync"
	"testing"
	"time"

	"github.com/zeroFruit/vnet/pkg/link"
)

type mockNetHandler struct {
	handleFunc func(rawPl []byte)
}

func (h *mockNetHandler) Handle(pl []byte) {
	h.handleFunc(pl)
}

func TestNetworkTopology(t *testing.T) {
	t.Skip()
	// setup node
	node1 := link.NewHost()
	node2 := link.NewHost()

	// setup interface
	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), node1)
	intf2 := link.NewInterface(40002, link.AddrFromStr("11-11-11-11-11-12"), node2)
	attachInterface(t, node1, intf1)
	attachInterface(t, node2, intf2)

	// setup link
	link1 := link.NewLink(1)
	attachLink(t, intf1, link1)
	attachLink(t, intf2, link1)

	testNode(t, node1, "11-11-11-11-11-11", 1)
	testNode(t, node2, "11-11-11-11-11-12", 1)
}

func TestNodeSendReceive(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}
	wg.Add(1)

	sender := link.NewHost()
	receiver := link.NewHost()

	intf1 := link.NewInterface(40001, link.AddrFromStr("11-11-11-11-11-11"), sender)
	intf2 := link.NewInterface(40002, link.AddrFromStr("11-11-11-11-11-12"), receiver)
	attachInterface(t, sender, intf1)
	attachInterface(t, receiver, intf2)

	link1 := link.NewLink(1)
	attachLink(t, intf1, link1)
	attachLink(t, intf2, link1)

	fh := &mockNetHandler{
		handleFunc: func(pl []byte) {
			if string(pl) != "hello" {
				t.Fatalf("expected pl is 'hello', but got '%s'", string(pl))
			}
			time.Sleep(1)
			wg.Done()
		},
	}
	receiver.RegisterNetHandler(fh)
	if err := sender.Send(link.AddrFromStr("11-11-11-11-11-12"), []byte("hello")); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
	wg.Wait()
}

func attachInterface(t *testing.T, node *link.Host, itf link.Interface) {
	node.AttachInterface(itf)
}

func attachLink(t *testing.T, itf link.Interface, link *link.Link) {
	if err := itf.AttachLink(link); err != nil {
		t.Fatalf("failed to attach link: %v", err)
	}
}

func testNode(t *testing.T, node *link.Host, addr link.Addr, cost uint) {
	intf1_ := node.Interface
	link1_ := intf1_.GetLink()
	if link1_ == nil {
		t.Fatalf("link not exist on interface with address: %s", intf1_.Address())
	}
	otherLink, err := link1_.GetOtherInterface(intf1_.Address())
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
