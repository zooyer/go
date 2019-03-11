package fastdfs

import (
	"testing"
	"fmt"
	"time"
)

func TestNewFileInfo(t *testing.T) {
	info := NewFileInfo(1024, time.Now().Unix(), 123, "127.0.0.1")
	fmt.Println(info)
}