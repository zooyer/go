package fastdfs

import (
	"net"
	"strings"
	"fmt"
	"os"
	"runtime/debug"
)

type TrackerClient struct {
	trackerGroup *TrackerGroup
	errno byte
}

/**
 * constructor with global tracker group
 */
func NewTrackerClient() *TrackerClient {
	return &TrackerClient{
		trackerGroup:GTrackerGroup,
	}
}

/**
 * constructor with specified tracker group
 *
 * @param tracker_group the tracker group object
 */
func NewTrackerClientByGroup(trackerGroup *TrackerGroup) *TrackerClient {
	return &TrackerClient{
		trackerGroup:trackerGroup,
	}
}

/**
 * get the error code of last call
 *
 * @return the error code of last call
 */
func (t *TrackerClient) GetErrorCode() byte {
	return t.errno
}

/**
 * get a connection to tracker server
 *
 * @return tracker server Socket object, return null if fail
 */
func (t *TrackerClient) GetConnection() (*TrackerServer, error) {
	return t.trackerGroup.GetConnection()
}

/**
 * query storage server to upload file
 *
 * @param trackerServer the tracker server
 * @return storage server Socket object, return null if fail
 */
func (t *TrackerClient) GetStoreStorage(trackerServer *TrackerServer) (*StorageServer, error) {
	const groupName = ""

	return t.GetStoreStorageByGroup(trackerServer, groupName)
}

/**
 * query storage server to upload file
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name to upload file to, can be empty
 * @return storage server object, return null if fail
 */
func (t *TrackerClient) GetStoreStorageByGroup(trackerServer *TrackerServer, groupName string) (*StorageServer, error) {
	var (
		header []byte
		ipAddr string
		port int
		cmd byte
		outLen int
		bNewConnection bool
		storePath byte
		trackerSocket net.Conn
	)
	var err error

	if trackerServer == nil {
		if trackerServer,err = t.GetConnection(); err != nil {
			return nil, err
		}
		if trackerServer == nil {
			// todo return nil or error?
			// java is return null.
			return nil, nil
		}
		bNewConnection = true
	} else {
		bNewConnection = false
	}
	defer func() {
		if bNewConnection {
			if err = trackerServer.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
				debug.PrintStack()
			}
		}
	}()

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return nil, err
	}

	if groupName == "" || len(groupName) == 0 {
		cmd = TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE
		outLen = 0
	} else {
		cmd = TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE
		outLen = FDFS_GROUP_NAME_MAX_LEN
	}
	if header,err = PackHeader(cmd, int64(outLen), 0); err != nil {
		return nil, err
	}
	if _,err = trackerSocket.Write(header); err != nil {
		return nil, err
	}

	if groupName != "" && len(groupName) > 0 {
		var bGroupName []byte
		var bs []byte
		var groupLen int

		if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
			return nil, err
		}
		bGroupName = make([]byte, FDFS_GROUP_NAME_MAX_LEN)

		if len(bs) <= FDFS_GROUP_NAME_MAX_LEN {
			groupLen = len(bs)
		} else {
			groupLen = FDFS_GROUP_NAME_MAX_LEN
		}
		copy(bGroupName[:groupLen], bs)
		if _,err = trackerSocket.Write(bGroupName); err != nil {
			return nil, err
		}
	}

	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, TRACKER_QUERY_STORAGE_STORE_BODY_LEN)
	if err != nil {
		return nil, err
	}
	t.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	ipAddr = strings.TrimSpace(string(pkgInfo.Body[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1]))
	port = int(Buff2long(pkgInfo.Body, FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1))
	storePath = pkgInfo.Body[TRACKER_QUERY_STORAGE_STORE_BODY_LEN - 1]

	return NewStorageServerByByte(ipAddr, port, storePath)
}

/**
 * query storage servers to upload file
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name to upload file to, can be empty
 * @return storage servers, return null if fail
 */
func (t *TrackerClient) GetStoreStorages(trackerServer *TrackerServer, groupName string) ([]*StorageServer, error) {
	var (
		header []byte
		ipAddr string
		port int
		cmd byte
		outLen int
		bNewConnection bool
		trackerSocket net.Conn
	)
	var err error

	if trackerServer == nil {
		if trackerServer,err = t.GetConnection(); err != nil {
			return nil, err
		}
		if trackerServer == nil {
			// todo return nil or error?
			// java is return null.
			return nil, nil
		}
		bNewConnection = true
	} else {
		bNewConnection = false
	}
	defer func() {
		if bNewConnection {
			if err = trackerServer.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
				debug.PrintStack()
			}
		}
	}()

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return nil, err
	}

	if groupName == "" || len(groupName) == 0 {
		cmd = TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL
		outLen = 0
	} else {
		cmd = TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL
		outLen = FDFS_GROUP_NAME_MAX_LEN
	}
	if header,err = PackHeader(cmd, int64(outLen), 0); err != nil {
		return nil, err
	}
	if _,err = trackerSocket.Write(header); err != nil {
		return nil, err
	}

	if groupName != "" && len(groupName) > 0 {
		var bGroupName []byte
		var bs []byte
		var groupLen int

		if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
			return nil, err
		}
		bGroupName = make([]byte, FDFS_GROUP_NAME_MAX_LEN)

		if len(bs) <= FDFS_GROUP_NAME_MAX_LEN {
			groupLen = len(bs)
		} else {
			groupLen = FDFS_GROUP_NAME_MAX_LEN
		}
		copy(bGroupName[:groupLen], bs)
		if _,err = trackerSocket.Write(bGroupName); err != nil {
			return nil, err
		}
	}

	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, -1)
	if err != nil {
		return nil, err
	}
	t.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	if len(pkgInfo.Body) < 0 {
		t.errno = ERR_NO_EINVAL
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var ipPortLen = len(pkgInfo.Body) - FDFS_GROUP_NAME_MAX_LEN + 1
	const recordLength = FDFS_IPADDR_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE

	if ipPortLen % recordLength != 0 {
		t.errno = ERR_NO_EINVAL
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var serverCount = ipPortLen / recordLength
	if serverCount > 16 {
		t.errno = ERR_NO_EINVAL
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var results = make([]*StorageServer, serverCount)
	var storePath = pkgInfo.Body[len(pkgInfo.Body) - 1]
	var offset = FDFS_GROUP_NAME_MAX_LEN

	for i := 0; i < serverCount; i++ {
		ipAddr = strings.TrimSpace(string(pkgInfo.Body[offset:offset + FDFS_IPADDR_SIZE - 1]))
		offset += FDFS_IPADDR_SIZE - 1
		port = int(Buff2long(pkgInfo.Body, offset))
		offset += FDFS_PROTO_PKG_LEN_SIZE

		if results[i],err = NewStorageServerByByte(ipAddr, port, storePath); err != nil {
			return nil, err
		}
	}

	return results, nil
}

/**
 * query storage server to download file
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name of storage server
 * @param filename      filename on storage server
 * @return storage server Socket object, return null if fail
 */
func (t *TrackerClient) GetFetchStorage(trackerServer *TrackerServer, groupName, filename string) (*StorageServer, error) {
	var servers,err = t.GetStorages(trackerServer, TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, groupName, filename)
	if err != nil {
		return nil, err
	}
	if servers == nil {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	return NewStorageServer(servers[0].GetIpAddr(), servers[0].GetPort(), 0)
}

/**
 * query storage server to update file (delete file or set meta data)
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name of storage server
 * @param filename      filename on storage server
 * @return storage server Socket object, return null if fail
 */
func (t *TrackerClient) GetUpdateStorage(trackerServer *TrackerServer, groupName, filename string) (*StorageServer, error) {
	var servers,err = t.GetStorages(trackerServer, TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE, groupName, filename)
	if err != nil {
		return nil, err
	}
	if servers == nil {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	return NewStorageServer(servers[0].GetIpAddr(), servers[0].GetPort(), 0)
}

/**
 * get storage servers to download file
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name of storage server
 * @param filename      filename on storage server
 * @return storage servers, return null if fail
 */
func (t *TrackerClient) GetFetchStorages(trackerServer *TrackerServer, groupName, filename string) ([]*ServerInfo, error) {
	return t.GetStorages(trackerServer, TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL, groupName, filename)
}

/**
 * query storage server to download file
 *
 * @param trackerServer the tracker server
 * @param cmd           command code, ProtoCommon.TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE or
 *                      ProtoCommon.TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE
 * @param groupName     the group name of storage server
 * @param filename      filename on storage server
 * @return storage server Socket object, return null if fail
 */
func (t *TrackerClient) GetStorages(trackerServer *TrackerServer, cmd byte, groupName, filename string) ([]*ServerInfo, error) {
	var (
		header []byte
		bFileName []byte
		bGroupName []byte
		bs []byte

		length int
		ipAddr string
		port int
		bNewConnection bool
		trackerSocket net.Conn
	)
	var err error

	if trackerServer == nil {
		if trackerServer,err = t.GetConnection(); err != nil {
			return nil, err
		}
		if trackerServer == nil {
			// todo return nil or error?
			// java is return null.
			return nil, nil
		}
		bNewConnection = true
	} else {
		bNewConnection = false
	}
	defer func() {
		if bNewConnection {
			if err = trackerServer.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
				debug.PrintStack()
			}
		}
	}()

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return nil, err
	}

	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return nil, err
	}
	bGroupName = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if bFileName,err = ConvertBytesToUTF8([]byte(filename), GCharset); err != nil {
		return nil, err
	}

	if len(bs) <= FDFS_GROUP_NAME_MAX_LEN {
		length = len(bs)
	} else {
		length = FDFS_GROUP_NAME_MAX_LEN
	}
	copy(bGroupName[:length], bs)

	if header,err = PackHeader(cmd, int64(FDFS_GROUP_NAME_MAX_LEN + len(bFileName)), 0); err != nil {
		return nil, err
	}

	var wholePkg = make([]byte, len(header) + len(bGroupName) + len(bFileName))
	copy(wholePkg, header)
	copy(wholePkg[len(header):], bGroupName)
	copy(wholePkg[len(header) + len(bGroupName):], bFileName)
	if _,err = trackerSocket.Write(wholePkg); err != nil {
		return nil, err
	}

	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, -1)
	if err != nil {
		return nil, err
	}
	t.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	if len(pkgInfo.Body) < TRACKER_QUERY_STORAGE_FETCH_BODY_LEN {
		return nil, fmt.Errorf("invalid body length: %d", len(pkgInfo.Body))
	}

	if (len(pkgInfo.Body) - TRACKER_QUERY_STORAGE_FETCH_BODY_LEN) % (FDFS_IPADDR_SIZE - 1) != 0 {
		return nil, fmt.Errorf("invalid body length: %d", len(pkgInfo.Body))
	}

	var serverCount = 1 + (len(pkgInfo.Body) - TRACKER_QUERY_STORAGE_FETCH_BODY_LEN) / (FDFS_IPADDR_SIZE - 1)
	ipAddr = strings.TrimSpace(string(pkgInfo.Body[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1]))
	var offset = FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1
	port = int(Buff2long(pkgInfo.Body, offset))
	offset += FDFS_PROTO_PKG_LEN_SIZE

	var servers = make([]*ServerInfo, serverCount)
	servers[0] = NewServerInfo(ipAddr, port)
	for i := 1; i < serverCount; i++ {
		servers[i] = NewServerInfo(strings.TrimSpace(string(pkgInfo.Body[offset: offset + FDFS_IPADDR_SIZE - 1])), port)
		offset += FDFS_IPADDR_SIZE - 1
	}

	return servers, nil
}

/**
 * query storage server to download file
 *
 * @param trackerServer the tracker server
 * @param file_id       the file id(including group name and filename)
 * @return storage server Socket object, return null if fail
 */
func (t *TrackerClient) GetFetchStorage1(trackerServer *TrackerServer, fileId string) (*StorageServer, error) {
	var parts = make([]string, 2)
	t.errno = SplitFileId(fileId, parts)
	if t.errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	return t.GetFetchStorage(trackerServer, parts[0], parts[1])
}

/**
 * get storage servers to download file
 *
 * @param trackerServer the tracker server
 * @param file_id       the file id(including group name and filename)
 * @return storage servers, return null if fail
 */
func (t *TrackerClient) GetFetchStorages1(trackerServer *TrackerServer, fileId string) ([]*ServerInfo, error) {
	var parts = make([]string, 2)
	t.errno = SplitFileId(fileId, parts)
	if t.errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	return t.GetFetchStorages(trackerServer, parts[0], parts[1])
}

/**
 * list groups
 *
 * @param trackerServer the tracker server
 * @return group stat array, return null if fail
 */
func (t *TrackerClient) ListGroups(trackerServer *TrackerServer) ([]StructGroupStat, error) {
	var (
		header []byte
		//ipAddr string
		//port int
		//cmd byte
		//outLen int
		bNewConnection bool
		//storePath byte
		trackerSocket net.Conn
	)
	var err error

	if trackerServer == nil {
		if trackerServer,err = t.GetConnection(); err != nil {
			return nil, err
		}
		if trackerServer == nil {
			// todo return nil or error?
			// java is return null.
			return nil, nil
		}
		bNewConnection = true
	} else {
		bNewConnection = false
	}
	defer func() {
		if bNewConnection {
			if err = trackerServer.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
				debug.PrintStack()
			}
		}
	}()

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return nil, err
	}

	if header,err = PackHeader(TRACKER_PROTO_CMD_SERVER_LIST_GROUP, 0, 0); err != nil {
		return nil, err
	}
	if _,err = trackerSocket.Write(header); err != nil {
		return nil, err
	}
	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, -1)
	if err != nil {
		return nil, err
	}
	t.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var decoder = NewProtoStructDecoder()
	inters,err := decoder.Decode(pkgInfo.Body, StructGroupStat{}, GetGroupFieldsTotalSize())
	if err != nil {
		return nil, err
	}
	var stats = make([]StructGroupStat, len(inters))
	for i,_ := range inters {
		stats[i] = inters[i].(StructGroupStat)
	}
	return stats, nil
}

/**
 * query storage server stat info of the group
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name of storage server
 * @return storage server stat array, return null if fail
 */
func (t *TrackerClient) ListStorages(trackerServer *TrackerServer, groupName string) ([]StructStorageStat, error) {
	const storageIpaddr = ""

	return t.ListStoragesByIpAddress(trackerServer, groupName, storageIpaddr)
}

/**
 * query storage server stat info of the group
 *
 * @param trackerServer the tracker server
 * @param groupName     the group name of storage server
 * @param storageIpAddr the storage server ip address, can be null or empty
 * @return storage server stat array, return null if fail
 */
func (t *TrackerClient) ListStoragesByIpAddress(trackerServer *TrackerServer, groupName, storageIpAddr string) ([]StructStorageStat, error) {
	var (
		header []byte
		bGroupName []byte
		bs []byte
		length int
		bNewConnection bool
		trackerSocket net.Conn
	)
	var err error

	if trackerServer == nil {
		if trackerServer,err = t.GetConnection(); err != nil {
			return nil, err
		}
		if trackerServer == nil {
			// todo return nil or error?
			// java is return null.
			return nil, nil
		}
		bNewConnection = true
	} else {
		bNewConnection = false
	}
	defer func() {
		if bNewConnection {
			if err = trackerServer.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
				debug.PrintStack()
			}
		}
	}()

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return nil, err
	}
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return nil, err
	}
	bGroupName = make([]byte, FDFS_GROUP_NAME_MAX_LEN)

	if len(bs) <= FDFS_GROUP_NAME_MAX_LEN {
		length = len(bs)
	} else {
		length = FDFS_GROUP_NAME_MAX_LEN
	}

	copy(bGroupName[:length], bs)

	var ipAddrLen int
	var bIpAddr []byte
	if storageIpAddr != "" && len(storageIpAddr) > 0 {
		if bIpAddr,err = ConvertBytesToUTF8([]byte(storageIpAddr), GCharset); err != nil {
			return nil, err
		}
		if len(bIpAddr) < FDFS_IPADDR_SIZE {
			ipAddrLen = len(bIpAddr)
		} else {
			ipAddrLen = FDFS_IPADDR_SIZE - 1
		}
	} else {
		bIpAddr = nil
		ipAddrLen = 0
	}

	if header,err = PackHeader(TRACKER_PROTO_CMD_SERVER_LIST_STORAGE, int64(FDFS_GROUP_NAME_MAX_LEN + ipAddrLen), 0); err != nil {
		return nil, err
	}
	var wholePkg = make([]byte, len(header) + len(bGroupName) + ipAddrLen)
	copy(wholePkg, header)
	copy(wholePkg[len(header):], bGroupName)
	if ipAddrLen > 0 {
		copy(wholePkg[len(header) + len(bGroupName):], bIpAddr)
	}
	if _,err = trackerSocket.Write(wholePkg); err != nil {
		return nil, err
	}

	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, -1)
	t.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var decoder = NewProtoStructDecoder()
	inter,err := decoder.Decode(pkgInfo.Body, StructStorageStat{}, GetStorageFieldsTotalSize())
	if err != nil {
		return nil, err
	}
	var stats = make([]StructStorageStat, len(inter))
	for i,_ := range inter {
		stats[i] = inter[i].(StructStorageStat)
	}

	return stats, nil
}

/**
 * delete a storage server from the tracker server
 *
 * @param trackerServer the connected tracker server
 * @param groupName     the group name of storage server
 * @param storageIpAddr the storage server ip address
 * @return true for success, false for fail
 */
func (t *TrackerClient) deleteStorage(trackerServer *TrackerServer, groupName, storageIpAddr string) (bool, error) {
	var (
		header []byte
		bGroupName []byte
		bs []byte
		length int
		trackerSocket net.Conn
	)
	var err error

	if trackerSocket,err = trackerServer.GetSocket(); err != nil {
		return false, err
	}
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return false, err
	}
	bGroupName = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if len(bs) <= FDFS_GROUP_NAME_MAX_LEN {
		length = len(bs)
	} else {
		length = FDFS_GROUP_NAME_MAX_LEN
	}
	copy(bGroupName[:length], bs)

	var ipAddrLen int
	bIpAddr,err:= ConvertBytesToUTF8([]byte(storageIpAddr), GCharset)
	if err != nil {
		return false, err
	}
	if len(bIpAddr) < FDFS_IPADDR_SIZE {
		ipAddrLen = FDFS_IPADDR_SIZE
	} else {
		ipAddrLen = FDFS_IPADDR_SIZE - 1
	}

	if header,err = PackHeader(TRACKER_PROTO_CMD_SERVER_DELETE_STORAGE, int64(FDFS_GROUP_NAME_MAX_LEN + ipAddrLen), 0); err != nil {
		return false, err
	}
	var wholePkg = make([]byte, len(header) + len(bGroupName) + ipAddrLen)
	copy(wholePkg, header)
	copy(wholePkg[len(header):], bGroupName)
	copy(wholePkg[len(header) + len(bGroupName):], bIpAddr)
	if _,err = trackerSocket.Write(wholePkg); err != nil {
		return false, err
	}

	pkgInfo,err := RecvPackage(trackerSocket, TRACKER_PROTO_CMD_RESP, 0)
	t.errno = pkgInfo.Errno

	return pkgInfo.Errno == 0, nil
}

/**
 * delete a storage server from the global FastDFS cluster
 *
 * @param groupName     the group name of storage server
 * @param storageIpAddr the storage server ip address
 * @return true for success, false for fail
 */
func (t *TrackerClient) DeleteStorage(groupName, storageIpAddr string) (bool, error) {
	return t.DeleteStorageByTrackerGroup(GTrackerGroup, groupName, storageIpAddr)
}

/**
 * delete a storage server from the FastDFS cluster
 *
 * @param trackerGroup  the tracker server group
 * @param groupName     the group name of storage server
 * @param storageIpAddr the storage server ip address
 * @return true for success, false for fail
 */
func (t *TrackerClient) DeleteStorageByTrackerGroup(trackerGroup *TrackerGroup, groupName, storageIpAddr string) (bool, error) {
	var (
		serverIndex int
		notFoundCount int
		trackerServer *TrackerServer
	)
	var err error
	
	notFoundCount = 0
	for serverIndex = 0; serverIndex < len(trackerGroup.TrackerServers); serverIndex++ {
		if trackerServer,err = trackerGroup.GetConnection(); err != nil {
			t.errno = ECONNREFUSED
			// todo return nil or error?
			// java is return null.
			return false, nil
		}

		storageStats,err := t.ListStoragesByIpAddress(trackerServer, groupName, storageIpAddr)
		if err != nil {
			trackerServer.Close()
			return false, err
		}
		if storageStats == nil {
			if t.errno == ERR_NO_ENOENT {
				notFoundCount++
			} else {
				trackerServer.Close()
				return false, nil
			}
		} else if len(storageStats) == 0 {
			notFoundCount++
		} else if storageStats[0].GetStatus() == FDFS_STORAGE_STATUS_ONLINE || storageStats[0].GetStatus() == FDFS_STORAGE_STATUS_ACTIVE {
			t.errno = ERR_NO_EBUSY
			trackerServer.Close()
			return false, nil
		}
		trackerServer.Close()
	}

	if notFoundCount == len(trackerGroup.TrackerServers) {
		t.errno = ERR_NO_ENOENT
		return false, nil
	}

	notFoundCount = 0
	for serverIndex = 0; serverIndex < len(trackerGroup.TrackerServers); serverIndex++ {
		if trackerServer,err = trackerGroup.GetConnectionByIndex(serverIndex); err != nil {
			fmt.Fprintln(os.Stderr, "connect to server ", trackerGroup.TrackerServers[serverIndex].String(), " fail")
			t.errno = ECONNREFUSED
			// todo return nil or error?
			// java is return null.
			return false, nil
		}

		if ok,err := t.deleteStorage(trackerServer, groupName, storageIpAddr); err != nil {
			trackerServer.Close()
			return false, err
		} else if !ok {
			if t.errno != 0 {
				if t.errno == ERR_NO_ENOENT {
					notFoundCount++
				} else if t.errno != ERR_NO_EALREADY {
					trackerServer.Close()
					return false, nil
				}
			}
		}
		trackerServer.Close()
	}

	if notFoundCount == len(trackerGroup.TrackerServers) {
		t.errno = ERR_NO_ENOENT
		return false, nil
	}

	if t.errno == ERR_NO_ENOENT {
		t.errno = 0
	}

	return t.errno == 0, nil
}


