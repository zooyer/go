package fastdfs

import (
	"reflect"
	"fmt"
)

type ProtoStructDecoder struct {

}

func NewProtoStructDecoder() *ProtoStructDecoder {
	return new(ProtoStructDecoder)
}

func (p *ProtoStructDecoder) Decode(bs []byte, types interface{}, fieldsTotalSize int) ([]interface{}, error) {
	if len(bs) % fieldsTotalSize != 0 {
		return nil, fmt.Errorf("byte array length: %d is invalid", len(bs))
	}

	var count = len(bs) / fieldsTotalSize
	var offset int
	var results = make([]interface{}, count)

	offset = 0
	for i := 0; i < len(results); i++ {
		results[i] = reflect.New(reflect.TypeOf(types).Elem()).Interface()
		results[i].(*StructBase).SetFields(bs, offset)
		offset += fieldsTotalSize
	}

	return results, nil
}