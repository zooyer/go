package fastdfs

import (
	"testing"
	"fmt"
)

func TestNewIniFileReader(t *testing.T) {
	r,err := NewIniFileReader("test/fdfs_client.conf.sample")
	if err != nil {
		panic(err)
	}

	fmt.Println(r.GetConfFilename())
	fmt.Println(r.GetIntValue("connect_timeout", 0))
	fmt.Println(r.GetIntValue("network_timeout", 0))
	fmt.Println(r.GetStrValue("charset"))
	fmt.Println(r.GetIntValue("http.tracker_http_port", 0))
	fmt.Println(r.GetBoolValue("http.anti_steal_token", false))
	fmt.Println(r.GetStrValue("http.secret_key"))
	fmt.Println(r.GetValues("tracker_server"))
}