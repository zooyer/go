package fastdfs

import (
	"time"
	"strings"
)

type StructBaseInterface interface {
	SetFields(bs []byte, offset int)
	stringValue(bs []byte, offset int, fieldInfo *FieldInfo) string
	int64Value(bs []byte, offset int, fieldInfo *FieldInfo) int64
	longValue(bs []byte, offset int, fieldInfo *FieldInfo) int64
	intValue(bs []byte, offset int, fieldInfo *FieldInfo) int
	int32Value(bs []byte, offset int, fieldInfo *FieldInfo) int32
	byteValue(bs []byte, offset int, fieldInfo *FieldInfo) byte
	boolValue(bs []byte, offset int, fieldInfo *FieldInfo) bool
	dateValue(bs []byte, offset int, fieldInfo *FieldInfo) time.Time
}

type StructBase struct {
	_ struct{}
}

/**
 * set fields
 *
 * @param bs     byte array
 * @param offset start offset
 */
func (s *StructBase) SetFields(bs []byte, offset int) {

}

func (s *StructBase) stringValue(bs []byte, offset int, fieldInfo *FieldInfo) string {
	if bytes,err := ConvertBytesToUTF8(bs, GCharset); err == nil {
		bs = bytes
	}

	return strings.TrimSpace(string(bs[offset + fieldInfo.offset:offset + fieldInfo.offset + fieldInfo.size]))
}

func (s *StructBase) int64Value(bs []byte, offset int, fieldInfo *FieldInfo) int64 {
	return Buff2long(bs, offset + fieldInfo.offset)
}

func (s *StructBase) longValue(bs []byte, offset int, fieldInfo *FieldInfo) int64 {
	return Buff2long(bs, offset + fieldInfo.offset)
}

func (s *StructBase) intValue(bs []byte, offset int, fieldInfo *FieldInfo) int {
	return int(Buff2long(bs, offset + fieldInfo.offset))
}

func (s *StructBase) int32Value(bs []byte, offset int, fieldInfo *FieldInfo) int32 {
	return Buff2int32(bs, offset + fieldInfo.offset)
}

func (s *StructBase) byteValue(bs []byte, offset int, fieldInfo *FieldInfo) byte {
	return bs[offset + fieldInfo.offset]
}

func (s *StructBase) boolValue(bs []byte, offset int, fieldInfo *FieldInfo) bool {
	return bs[offset + fieldInfo.offset] != 0
}

func (s *StructBase) dateValue(bs []byte, offset int, fieldInfo *FieldInfo) time.Time {
	return time.Unix(Buff2long(bs, offset + fieldInfo.offset), 0)
}

type FieldInfo struct {
	name   string
	offset int
	size   int
}

func NewFieldInfo(name string, offset, size int) *FieldInfo {
	return &FieldInfo{
		name   : name,
		offset : offset,
		size   : size,
	}
}