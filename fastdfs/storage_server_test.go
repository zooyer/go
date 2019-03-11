package fastdfs

import (
	"testing"
	"net"
	"fmt"
	"strings"
)

func TestNewStorageServer(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	var addr = listener.Addr().(*net.TCPAddr)
	fmt.Println("IP:", strings.Trim(addr.IP.String(), ":"))
	fmt.Println("PORT:", addr.Port)
	store,err := NewStorageServer(strings.Trim(addr.IP.String(), ":"), addr.Port, 1)
	if err != nil {
		panic(err)
	}

	store.Close()
}