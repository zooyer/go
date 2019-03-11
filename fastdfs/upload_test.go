package fastdfs

import (
	"testing"
	"os"
	"fmt"
	"bytes"
)

func TestNewUploadStream(t *testing.T) {
	var str = "HelloWorld"
	var u = NewUploadStream(bytes.NewReader([]byte(str)), len(str))
	if _,err := u.Send(os.Stdout); err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(u.fileSize)
}