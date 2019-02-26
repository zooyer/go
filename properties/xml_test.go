package properties

import (
	"testing"
	"encoding/xml"
	"io/ioutil"
	"fmt"
	"os"
	"sync"
)

func TestXML(t *testing.T) {
	data,err := ioutil.ReadFile("test/test.xml")
	if err != nil {
		panic(err)
	}
	var m XMLProperties
	if err = xml.Unmarshal(data, &m); err != nil {
		panic(err)
	}

	x,err := m.ToXML("")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(x))
	return

	for _,v := range m.Entry {
		//fmt.Println(v.Key + " = " + v.CData)
		fmt.Println(v.Key + " = " + v.CDATA)
	}

	return

	x,err = xml.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(err)
	}
	file,err := os.Create("test2.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(xml.Header))
	file.Write([]byte(Header))
	file.Write(x)
}

func TestString(t *testing.T) {
	str := "hello你好"

	fmt.Println(len(str))
	fmt.Println(len([]byte(str)))

	fmt.Println(string([]byte{str[5], str[6], str[7]}))
}

func TestDefer(t *testing.T) {
	mutex := sync.Mutex{}
	f := func() int {
		mutex.Lock()
		defer mutex.Unlock()

		rf := func() int {
			fmt.Println("return func...")
			return 0
		}

		return rf()
	}

	f()

}