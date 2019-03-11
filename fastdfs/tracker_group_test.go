package fastdfs

import (
	"testing"
	"net"
)

func TestNewTrackerGroup(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	var addr = make([]net.Addr, 10)
	for i := 0; i < 10; i++ {
		addr[i] = listener.Addr()
	}

	var group = NewTrackerGroup(addr)
	tracker,err := group.GetConnection()
	if err != nil {
		panic(err)
	}
	tracker.Close()

	tracker,err = group.GetConnectionByIndex(9)
	if err != nil {
		panic(err)
	}
	tracker.Close()

	var group2 = group.Clone()

	tracker,err = group2.GetConnection()
	if err != nil {
		panic(err)
	}
	tracker.Close()

	tracker,err = group2.GetConnectionByIndex(9)
	if err != nil {
		panic(err)
	}
	tracker.Close()
}