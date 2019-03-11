package fastdfs

import (
	"testing"
	"fmt"
)

func TestNewProtoStructDecoder(t *testing.T) {
	var buf = []byte("keynumstrbooint")
	aaa,err := NewProtoStructDecoder().Decode(buf, &StructBase{}, 3)
	if err != nil {
		panic(err)
	}

	for i,_ := range aaa {
		fmt.Println(aaa[i])
	}
}
