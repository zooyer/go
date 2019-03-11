package fastdfs

import (
	"testing"
	"os"
	"fmt"
)

func TestNewDownLoadStream(t *testing.T) {
	var d = NewDownLoadStream(os.Stdout)
	if _,err := d.Recv(10, []byte("HelloWorld"), 10); err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(d.currentBytes)
}