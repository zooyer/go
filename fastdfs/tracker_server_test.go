package fastdfs

import (
	"testing"
	"net"
)

func TestNewTrackerServer(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conn,err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		panic(err)
	}

	tracker := NewTrackerServer(conn, conn.LocalAddr())
	tracker.Close()
}