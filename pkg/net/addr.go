package net

import (
	"github.com/zeroFruit/vnet/pkg/types"
)

type Addr string

func (a Addr) Equal(o types.NetAddr) bool {
	return a == o
}

func (a Addr) String() string {
	return string(a)
}

func AddrFromStr(s string) types.NetAddr {
	return Addr(s)
}
