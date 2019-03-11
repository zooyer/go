package fastdfs

import (
	"net"
	"fmt"
)

type StorageServer struct {
	TrackerServer

	storePathIndex   int
}

/**
 * Constructor
 *
 * @param ip_addr    the ip address of storage server
 * @param port       the port of storage server
 * @param store_path the store path index on the storage server
 */
func NewStorageServer(ipAddr string, port, storePath int) (*StorageServer, error) {
	addr,err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipAddr, port))
	if err != nil {
		return nil, err
	}
	conn,err := GetSocket(ipAddr, port)
	if err != nil {
		return nil, err
	}

	return &StorageServer{
		TrackerServer:*NewTrackerServer(conn, addr),
		storePathIndex:storePath,
	}, nil
}

/**
 * Constructor
 *
 * @param ip_addr    the ip address of storage server
 * @param port       the port of storage server
 * @param store_path the store path index on the storage server
 */
func NewStorageServerByByte(ipAddr string, port int, storePath byte) (*StorageServer, error) {
	addr,err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipAddr, port))
	if err != nil {
		return nil, err
	}
	conn,err := GetSocket(ipAddr, port)
	if err != nil {
		return nil, err
	}

	var storePathIndex = int(storePath)

	// java code
	//var storePathIndex = int(storePath)
	//if int8(storePath) < 0 {
	//	storePathIndex = 256 + int(int8(storePath))
	//}

	return &StorageServer{
		TrackerServer:*NewTrackerServer(conn, addr),
		storePathIndex:storePathIndex,
	}, nil
}

/**
 * @return the store path index on the storage server
 */
func (s *StorageServer) GetStorePathIndex() int {
	return s.storePathIndex
}