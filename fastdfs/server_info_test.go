package fastdfs

import (
	"testing"
	"fmt"
	"net"
	"runtime"
)

func TestNewServerInfo(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go func() {
		conn,err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("new connection:", conn.RemoteAddr())
		conn.Close()
	}()

	runtime.Gosched()

	info := NewServerInfo("127.0.0.1", listener.Addr().(*net.TCPAddr).Port)
	fmt.Println(info)
	conn,err := info.Connect()
	if err != nil {
		panic(err)
	}
	runtime.Gosched()
	conn.Close()
}