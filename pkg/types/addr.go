package types

type NetAddr interface {
	Equal(o NetAddr) bool
	String() string
}

type HwAddr interface {
	Equal(o HwAddr) bool
	String() string
}
