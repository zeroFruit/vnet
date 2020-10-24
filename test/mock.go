package test

type MockNetHandler struct {
	HandleFunc func(pl []byte)
}

func (h *MockNetHandler) Handle(pl []byte) {
	h.HandleFunc(pl)
}
