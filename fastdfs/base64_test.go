package fastdfs

import "testing"

func TestNewBase64(t *testing.T) {
	var b64 = NewBase64()

	// encode

	var str64 = "CwUEFYoAAAADjQMC7ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267ELJiY6w05267EI="
	// decode
	theBytes,err := b64.Decode(str64)
	if err != nil {
		panic(err)
	}

	show(theBytes)
}
