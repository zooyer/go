package properties

import (
	"testing"
	"fmt"
)

func TestNewHashtable(t *testing.T) {
	h := NewHashtable()
	h.Init()
	h.Put("key1", "val1")
	h.Put("key2", 1000)
	h.Put("key3", 'a')
	fmt.Println(h.Size())
	fmt.Println(h.Keys())
	fmt.Println(h.Get("key1"))
	fmt.Println(h.Get("key2"))
	fmt.Println(h.Get("key3"))
	h.Remove("key2")
	fmt.Println(h.Size())
	fmt.Println(h.Keys())
}
