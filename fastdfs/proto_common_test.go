package fastdfs

import (
	"testing"
	"fmt"
	"bytes"
	"net"
	"time"
)

func TestPackHeader(t *testing.T) {
	head,err := PackHeader(0xAA, 0x1FFFFFFFFFFFFFFF, 0xBB)
	if err != nil {
		panic(err)
	}
	fmt.Printf("header : %2X\n", head)
}

func TestRecvHeader(t *testing.T) {
	head,err := PackHeader(FDFS_PROTO_CMD_ACTIVE_TEST, 1024, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%2x\n", head)

	info,err := RecvHeader(bytes.NewReader(head), FDFS_PROTO_CMD_ACTIVE_TEST, 1024)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}

func TestRecvPackage(t *testing.T) {
	head,err := PackHeader(FDFS_PROTO_CMD_ACTIVE_TEST, 10, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%2x\n", head)

	pkg,err := RecvPackage(bytes.NewReader(append(head, []byte("HelloWorld")...)), FDFS_PROTO_CMD_ACTIVE_TEST, 10)
	if err != nil {
		panic(err)
	}

	fmt.Println("pkg errno:", pkg.Errno)
	fmt.Println("pkg length:", len(pkg.Body))
	fmt.Println(string(pkg.Body))
}

func TestSplitMetadata(t *testing.T) {
	var vals = SplitMetadata("abc"+ FDFS_FIELD_SEPERATOR + "ads")
	fmt.Println("vals len:", len(vals))
	for i,_ := range vals {
		fmt.Println(vals[i].GetName())
		fmt.Println(vals[i].GetValue())
	}
}

func TestPackMetadata(t *testing.T) {
	var vals = make([]NameValuePair, 10)
	for i := 0; i < 10; i++ {
		vals[i] = *NewNameValuePair(fmt.Sprintf("key-%d", i + 1), fmt.Sprintf("val-%d", i + 1))
	}
	fmt.Println(PackMetadata(vals))
}

func TestCloseSocket(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go func() {
		conn,err := listener.Accept()
		if err != nil {
			panic(err)
		}
		var buf = make([]byte, 4096)
		for {
			n,err := conn.Read(buf)
			if err != nil {
				break
			}
			fmt.Println(fmt.Sprintf("%02x", buf[:n]))
		}
		conn.Close()
	}()

	conn,err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		panic(err)
	}
	CloseSocket(conn)
	time.Sleep(time.Second)
}

func TestActiveTest(t *testing.T) {
	listener,err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go func() {
		conn,err := listener.Accept()
		if err != nil {
			panic(err)
		}
		var buf = make([]byte, 4096)
		for {
			n,err := conn.Read(buf)
			if err != nil {
				break
			}
			fmt.Println(fmt.Sprintf("%02x", buf[:n]))
			res,err := PackHeader(TRACKER_PROTO_CMD_RESP, 0, 0)
			if err != nil {
				panic(err)
			}
			if _,err = conn.Write(res); err != nil {
				break
			}
		}
		conn.Close()
	}()

	conn,err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		panic(err)
	}
	if ok,err := ActiveTest(conn); err != nil {
		panic(err)
	} else {
		fmt.Println("is ok:", ok)
	}
	time.Sleep(time.Second)
}

func TestLong2Buff(t *testing.T) {
	var num = int64(^uint64(0) >> 1)
	fmt.Println(fmt.Sprintf("%2x", Long2Buff(num)))
}

func TestBuff2long(t *testing.T) {
	var buf = make([]byte, 8)
	for i := 0; i < len(buf); i++ {
		buf[i] = 255
	}

	fmt.Println(Buff2long(buf, 0))
	fmt.Println(Buff2int32(buf, 0))
}

func TestBuff2int32(t *testing.T) {
	var buf = []byte{127, 255, 255, 255}
	fmt.Println("int32 max:", int32(^uint32(0) >> 1))
	fmt.Println("int32 max mem:", fmt.Sprintf("%02x", int32(^uint32(0) >> 1)))
	fmt.Println("buff to int32:", Buff2int32(buf, 0))
	fmt.Println("buff to int32 mem:", fmt.Sprintf("%02x", Buff2int32(buf, 0)))
}

func TestMd5(t *testing.T) {
	fmt.Println(Md5([]byte("HelloWorld")))
}

func TestGetIpAddress(t *testing.T) {
	fmt.Println(GetIpAddress([]byte{192, 168, 1, 110}, 0))
}

func TestGetToken(t *testing.T) {
	token,err := GetToken("gourp1\\M00/00/00/ab/2f/ba2b.dat", 1502344576, "===SECRET-KEY===")
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", token)
}

func TestGenSlaveFilename(t *testing.T) {
	name,err := GenSlaveFilename("asfjsdkfldsahfdsjfjlsdfjslfjasfjsdkgourp1\\M00/00/00/ab/2f/ba2b.dat", "==", "--")
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}