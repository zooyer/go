package fastdfs

import (
	"testing"
	"fmt"
	"properties"
)

func TestGlobalClient(t *testing.T) {
	var trackerServers = "10.0.11.101:22122,10.0.11.102:22122"
	if err := InitByTrackers(trackerServers); err != nil {
		panic(err)
	}
	fmt.Println("ClientGlobal.configInfo() : " + ConfigInfo())

	var propFilePath = "test/fastdfs-client.properties.sample"
	if err := InitByPropertiesFile(propFilePath); err != nil {
		panic(err)
	}
	fmt.Println("ClientGlobal.configInfo() : " + ConfigInfo())

	var props = properties.NewProperties()
	props.Put(PropKeyTrackerServers, "10.0.11.101:22122,10.0.11.102:22122")
	if err := InitByProperties(props); err != nil {
		panic(err)
	}
	fmt.Println("ClientGlobal.configInfo() : " + ConfigInfo())
}