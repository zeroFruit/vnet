package net

type Addr string

func (a Addr) Equal(o Addr) bool {
	return a == o
}

func (a Addr) String() string {
	return string(a)
}