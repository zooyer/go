package properties

import (
	"testing"
	"os"
	"fmt"
	"runtime/debug"
	"unsafe"
)

func TestNewLineReader(t *testing.T) {
	file,err := os.Open("test/log4j2.properties")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lr := NewLineReader(file)
	n := lr.readLine()
	fmt.Println("readLine len:", n)
}

func TestNewProperties(t *testing.T) {
	var p = NewProperties2()

	file,err := os.Open("test/log4j2.properties")
	if err != nil {
		//debug.PrintStack()
		fmt.Println(string(debug.Stack()))
	}
	defer file.Close()
	if err = p.Load(file); err != nil {
		panic(err)
	}

	p.List(os.Stdout)
}

func TestNewProperties3(t *testing.T) {
	var p = NewProperties()
	file,err := os.Open("test/log4j2.xml")
	if err != nil {
		panic(err)
	}
	if err = p.LoadFromXML(file); err != nil {
		panic(err)
	}
	fmt.Println(p.GetProperty("status"))
	fmt.Println(p.Size())
	p.Remove("status")
	fmt.Println(p.GetProperty("status"))
	fmt.Println(p.StringPropertyNames())
	fmt.Println(p.Keys())

	//p.List(os.Stdout)

	p.Store(os.Stdout, "")
}

func TestLen(t *testing.T) {
	ch := make(chan []byte, 1024)
	fmt.Println(len(ch))
	ch <- []byte("")
	fmt.Println(len(ch))

	var data []byte
	fmt.Println(fmt.Sprintf("%02x", data))
}

func TestSizeof(t *testing.T) {
	var num int
	var addr uintptr
	fmt.Println(unsafe.Sizeof(num))
	fmt.Println(unsafe.Sizeof(addr))
}