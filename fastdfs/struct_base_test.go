package fastdfs

import (
	"testing"
	"encoding/binary"
	"bytes"
	"fmt"
	"time"
)

type data struct {
	id    int32
	is    byte
	i64   int64
	b     byte
}

func TestStructBase(t *testing.T) {
	var sb = new(StructBase)

	var buf = bytes.NewBuffer(nil)
	var d = data{id:1024, is:1, i64:123456789, b:200}
	if err := binary.Write(buf, binary.BigEndian, d); err != nil {
		panic(err)
	}
	var str = "HelloWorld"
	if err := binary.Write(buf, binary.BigEndian, []byte(str)); err != nil {
		panic(err)
	}
	if err := binary.Write(buf, binary.BigEndian, time.Now().Unix()); err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%02x", buf.Bytes()))

	fmt.Println("id:", sb.int32Value(buf.Bytes(), 0, NewFieldInfo("", 0, 4)))
	fmt.Println("is:", sb.boolValue(buf.Bytes(), 4, NewFieldInfo("is", 0, 1)))
	fmt.Println("i64:", sb.int64Value(buf.Bytes(), 5, NewFieldInfo("i64", 0, 8)))
	fmt.Println("b:", sb.byteValue(buf.Bytes(), 13, NewFieldInfo("b", 0, 1)))
	fmt.Println("str:", sb.stringValue(buf.Bytes(), 14, NewFieldInfo("str", 0, len(str))))
	fmt.Println("time:", sb.dateValue(buf.Bytes(), 14 + len(str), NewFieldInfo("time", 0, 8)))
}