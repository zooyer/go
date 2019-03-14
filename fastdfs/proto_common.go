package fastdfs

import (
	"io"
	"fmt"
	"strings"
	"bytes"
	"net"
	"strconv"
	"crypto/md5"
)

const (
	FDFS_PROTO_CMD_QUIT = 82
	TRACKER_PROTO_CMD_SERVER_LIST_GROUP = 91
	TRACKER_PROTO_CMD_SERVER_LIST_STORAGE = 92
	TRACKER_PROTO_CMD_SERVER_DELETE_STORAGE = 93
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE = 103
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE = 104
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL = 105
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL = 106
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL = 107
	TRACKER_PROTO_CMD_RESP = 100
	FDFS_PROTO_CMD_ACTIVE_TEST = 111
	STORAGE_PROTO_CMD_UPLOAD_FILE = 11
	STORAGE_PROTO_CMD_DELETE_FILE = 12
	STORAGE_PROTO_CMD_SET_METADATA = 13
	STORAGE_PROTO_CMD_DOWNLOAD_FILE = 14
	STORAGE_PROTO_CMD_GET_METADATA = 15
	STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE = 21
	STORAGE_PROTO_CMD_QUERY_FILE_INFO = 22
	STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE = 23  //create appender fil
	STORAGE_PROTO_CMD_APPEND_FILE = 24  //append file
	STORAGE_PROTO_CMD_MODIFY_FILE = 34  //modify appender file
	STORAGE_PROTO_CMD_TRUNCATE_FILE = 36  //truncate appender file
	STORAGE_PROTO_CMD_RESP = TRACKER_PROTO_CMD_RESP
	FDFS_STORAGE_STATUS_INIT = 0
	FDFS_STORAGE_STATUS_WAIT_SYNC = 1
	FDFS_STORAGE_STATUS_SYNCING = 2
	FDFS_STORAGE_STATUS_IP_CHANGED = 3
	FDFS_STORAGE_STATUS_DELETED = 4
	FDFS_STORAGE_STATUS_OFFLINE = 5
	FDFS_STORAGE_STATUS_ONLINE = 6
	FDFS_STORAGE_STATUS_ACTIVE = 7
	FDFS_STORAGE_STATUS_NONE = 99
)

/**
 * for overwrite all old metadata
 */
const STORAGE_SET_METADATA_FLAG_OVERWRITE = 'O'

/**
 * for replace, insert when the meta item not exist, otherwise update it
 */
const (
	STORAGE_SET_METADATA_FLAG_MERGE = 'M'
	FDFS_PROTO_PKG_LEN_SIZE = 8
	FDFS_PROTO_CMD_SIZE = 1
	FDFS_GROUP_NAME_MAX_LEN = 16
	FDFS_IPADDR_SIZE = 16
	FDFS_DOMAIN_NAME_MAX_SIZE = 128
	FDFS_VERSION_SIZE = 6
	FDFS_STORAGE_ID_MAX_SIZE = 16
	FDFS_RECORD_SEPERATOR = "\u0001"
	FDFS_FIELD_SEPERATOR = "\u0002"
	TRACKER_QUERY_STORAGE_FETCH_BODY_LEN = FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN = FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE + FDFS_PROTO_PKG_LEN_SIZE

	FDFS_FILE_EXT_NAME_MAX_LEN = 6
	FDFS_FILE_PREFIX_MAX_LEN = 16
	FDFS_FILE_PATH_LEN = 10
	FDFS_FILENAME_BASE64_LENGTH = 27
	FDFS_TRUNK_FILE_INFO_LEN = 16
	ERR_NO_ENOENT = 2
	ERR_NO_EIO = 5
	ERR_NO_EBUSY = 16
	ERR_NO_EINVAL = 22
	ERR_NO_ENOSPC = 28
	ECONNREFUSED = 61
	ERR_NO_EALREADY = 114
	INFINITE_FILE_SIZE = 256 * 1024 * 1024 * 1024 * 1024 * 1024
	APPENDER_FILE_SIZE = INFINITE_FILE_SIZE
	TRUNK_FILE_MARK_SIZE = 512 * 1024 * 1024 * 1024 * 1024 * 1024
	NORMAL_LOGIC_FILENAME_LENGTH = FDFS_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH + FDFS_FILE_EXT_NAME_MAX_LEN + 1
	TRUNK_LOGIC_FILENAME_LENGTH = NORMAL_LOGIC_FILENAME_LENGTH + FDFS_TRUNK_FILE_INFO_LEN
	PROTO_HEADER_CMD_INDEX = FDFS_PROTO_PKG_LEN_SIZE
	PROTO_HEADER_STATUS_INDEX = FDFS_PROTO_PKG_LEN_SIZE + 1
)

func GetStorageStatusCaption(status byte) string {
	switch status {
	case FDFS_STORAGE_STATUS_INIT:
		return "INIT"
	case FDFS_STORAGE_STATUS_WAIT_SYNC:
		return "WAIT_SYNC"
	case FDFS_STORAGE_STATUS_SYNCING:
		return "SYNCING"
	case FDFS_STORAGE_STATUS_IP_CHANGED:
		return "IP_CHANGED"
	case FDFS_STORAGE_STATUS_DELETED:
		return "DELETED"
	case FDFS_STORAGE_STATUS_OFFLINE:
		return "OFFLINE"
	case FDFS_STORAGE_STATUS_ONLINE:
		return "ONLINE"
	case FDFS_STORAGE_STATUS_ACTIVE:
		return "ACTIVE"
	case FDFS_STORAGE_STATUS_NONE:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

/**
 * pack header by FastDFS transfer protocol
 *
 * @param cmd     which command to send
 * @param pkg_len package body length
 * @param errno   status code, should be (byte)0
 * @return packed byte buffer
 */
func PackHeader(cmd byte, pkgLen int64, errno byte) ([]byte, error) {
	var header []byte
	var hexLen []byte

	header = make([]byte, FDFS_PROTO_PKG_LEN_SIZE + 2)
	hexLen = Long2Buff(pkgLen)
	copy(header, hexLen)

	header[PROTO_HEADER_CMD_INDEX] = cmd
	header[PROTO_HEADER_STATUS_INDEX] = errno

	return header, nil
}

/**
 * receive pack header
 *
 * @param in              input stream
 * @param expect_cmd      expect response command
 * @param expect_body_len expect response package body length
 * @return RecvHeaderInfo: errno and pkg body length
 */
func RecvHeader(in io.Reader, expectCmd byte, expectBodyLen int64) (*RecvHeaderInfo, error) {
	var header []byte
	var bytes int
	var pkgLen int64
	var err error

	header = make([]byte, FDFS_PROTO_PKG_LEN_SIZE + 2)

	bytes,err = in.Read(header)
	if err != nil {
		return nil, err
	}
	if bytes != len(header) {
		return nil, fmt.Errorf("recv package size %d != %d", bytes, len(header))
	}

	if header[PROTO_HEADER_CMD_INDEX] != expectCmd {
		return nil, fmt.Errorf("recv cmd: %d is not correct, expect cmd: %d", header[PROTO_HEADER_CMD_INDEX], expectCmd)
	}

	if header[PROTO_HEADER_STATUS_INDEX] != 0 {
		return NewRecvHeaderInfo(header[PROTO_HEADER_STATUS_INDEX], 0), nil
	}

	pkgLen = Buff2long(header, 0)
	if pkgLen < 0 {
		return nil, fmt.Errorf("recv body length: %d < 0", pkgLen)
	}
	//fmt.Println("pkgLen:", pkgLen)

	if expectBodyLen >= 0 && pkgLen != expectBodyLen {
		return nil, fmt.Errorf("recv body length: %d is not correct, expect length: %d", pkgLen, expectBodyLen)
	}

	return NewRecvHeaderInfo(0, int(pkgLen)), nil
}

/**
 * receive whole pack
 *
 * @param in              input stream
 * @param expect_cmd      expect response command
 * @param expect_body_len expect response package body length
 * @return RecvPackageInfo: errno and reponse body(byte buff)
 */
func RecvPackage(in io.Reader,  expectCmd byte, expectBodyLen int64) (*RecvPackageInfo, error) {
	var header,err = RecvHeader(in, expectCmd, expectBodyLen)
	if err != nil {
		return nil, err
	}

	if header.Errno != 0 {
		return NewRecvPackageInfo(header.Errno, nil), nil
	}

	var body = make([]byte, header.BodyLen)
	var totalBytes = 0
	var remainBytes = header.BodyLen
	var bytes int

	for totalBytes < header.BodyLen {
		// TODO java < 0, what is golang?
		bytes,err = in.Read(body[totalBytes:totalBytes + remainBytes])
		if err != nil || bytes < 0 {
			break
		}

		totalBytes += bytes
		remainBytes -= bytes
	}

	if totalBytes != header.BodyLen {
		return nil, fmt.Errorf("recv package size %d != %d", totalBytes, header.BodyLen)
	}

	return NewRecvPackageInfo(0, body), nil
}

/**
 * split metadata to name value pair array
 *
 * @param meta_buff metadata
 * @return name value pair array
 */
func SplitMetadata(metaBuff string) []NameValuePair {
	return SplitMetadata2(metaBuff, FDFS_RECORD_SEPERATOR, FDFS_FIELD_SEPERATOR)
}

/**
 * split metadata to name value pair array
 *
 * @param meta_buff       metadata
 * @param recordSeperator record/row seperator
 * @param filedSeperator  field/column seperator
 * @return name value pair array
 */
func SplitMetadata2(metaBuff, recordSeperator, filedSeperator string) []NameValuePair {
	var rows []string
	var cols []string
	var metaList []NameValuePair

	rows = strings.Split(metaBuff, recordSeperator)
	metaList = make([]NameValuePair, len(rows))

	for i := 0; i < len(rows); i++ {
		cols = strings.SplitN(rows[i], filedSeperator, 2)
		metaList[i] = *NewNameValuePair(cols[0], "")
		if len(cols) == 2 {
			metaList[i].SetValue(cols[1])
		}
	}

	return metaList
}

/**
 * pack metadata array to string
 *
 * @param meta_list metadata array
 * @return packed metadata
 */
func PackMetadata(metaList []NameValuePair) string {
	if len(metaList) == 0 {
		return ""
	}

	var sb = bytes.NewBufferString("")
	sb.WriteString(metaList[0].GetName())
	sb.WriteString(FDFS_FIELD_SEPERATOR)
	sb.WriteString(metaList[0].GetValue())

	for i := 1; i < len(metaList); i++ {
		sb.WriteString(FDFS_RECORD_SEPERATOR)
		sb.WriteString(metaList[i].GetName())
		sb.WriteString(FDFS_FIELD_SEPERATOR)
		sb.WriteString(metaList[i].GetValue())
	}

	return sb.String()
}

/**
 * send quit command to server and close socket
 *
 * @param sock the Socket object
 */
func CloseSocket(conn net.Conn) error {
	var header []byte
	header,err := PackHeader(FDFS_PROTO_CMD_QUIT, 0, 0)
	if err != nil {
		return err
	}
	// TODO write not need return.
	if _,err = conn.Write(header); err != nil {
		return err
	}
	return conn.Close()
}

/**
 * send ACTIVE_TEST command to server, test if network is ok and the server is alive
 *
 * @param sock the Socket object
 */
func ActiveTest(conn net.Conn) (bool, error) {
	var header []byte
	header,err := PackHeader(FDFS_PROTO_CMD_ACTIVE_TEST, 0, 0)
	if err != nil {
		return false, err
	}
	if _,err = conn.Write(header); err != nil {
		return false, err
	}

	headerInfo,err := RecvHeader(conn, TRACKER_PROTO_CMD_RESP, 0)
	if err != nil {
		return false, err
	}
	if headerInfo.Errno == 0 {
		return true, nil
	}

	return false, nil
}


/**
 * long convert to buff (big-endian)
 *
 * @param n long number
 * @return 8 bytes buff
 */
func Long2Buff(n int64) []byte {
	var bs = make([]byte, 8)

	bs[0] = byte((n >> 56) & 0xFF)
	bs[1] = byte((n >> 48) & 0xFF)
	bs[2] = byte((n >> 40) & 0xFF)
	bs[3] = byte((n >> 32) & 0xFF)
	bs[4] = byte((n >> 24) & 0xFF)
	bs[5] = byte((n >> 16) & 0xFF)
	bs[6] = byte((n >> 8) & 0xFF)
	bs[7] = byte(n & 0xFF)

	return bs
}

/**
 * buff convert to long
 *
 * @param bs     the buffer (big-endian)
 * @param offset the start position based 0
 * @return long number
 */
func Buff2long(bs []byte, offset int) int64 {
	var long int64 = 0

	for i := 0; i < 8; i++ {
		long |= int64(bs[offset + i]) << uint64(56 - i * 8)
	}

	// java code
	//for i := 0; i < 8; i++ {
	//	if int8(bs[offset + i]) >= 0 {
	//		long |= int64(int8(bs[offset + i])) << uint64(56 - i * 8)
	//	} else {
	//		long |= (256 + int64(int8(bs[offset + i]))) << uint64(56 - i * 8)
	//	}
	//}

	return long
}

/**
 * buff convert to int
 *
 * @param bs     the buffer (big-endian)
 * @param offset the start position based 0
 * @return int number
 */
func Buff2int32(bs []byte, offset int) int32 {
	var num int32 = 0

	for i := 0; i < 4; i++ {
		num |= int32(bs[offset + i]) << uint32(24 - i * 8)
	}

	// java code
	//for i := 0; i < 4; i++ {
	//	if int8(bs[offset + i]) >= 0 {
	//		num |= int32(int8(bs[offset + i])) << uint32(24 - i * 8)
	//	} else {
	//		num |= (256 + int32(int8(bs[offset + i]))) << uint32(24 - i * 8)
	//	}
	//}

	return num
}

/**
 * buff convert to ip address
 *
 * @param bs     the buffer (big-endian)
 * @param offset the start position based 0
 * @return ip address
 */
func GetIpAddress(bs []byte, offset int) string {
	if bs[0] == 0 || bs[3] == 0 {
		return ""
	}
	var n int
	var sbResult = bytes.NewBufferString("")
	for i := offset; i < offset + 4; i++ {
		n = int(bs[i])
		if sbResult.Len() > 0 {
			sbResult.WriteString(".")
		}
		sbResult.WriteString(strconv.Itoa(n))
	}

	// java code
	//for i := offset; i < offset + 4; i++ {
	//	if int8(bs[i]) >= 0 {
	//		n = int(bs[i])
	//	} else {
	//		n = 256 + int(int8(bs[i]))
	//	}
	//	if sbResult.Len() > 0 {
	//		sbResult.WriteString(".")
	//	}
	//	sbResult.WriteString(strconv.Itoa(n))
	//}

	return sbResult.String()
}

/**
 * md5 function
 *
 * @param source the input buffer
 * @return md5 string
 */
func Md5(source []byte) string {
	var hexDigits = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	var md5num = md5.Sum(source)

	return fmt.Sprintf("%2x", md5num)

	// TODO java
	var str = make([]byte, 32)
	var k = 0
	for i := 0; i < 16; i++ {
		str[k] = hexDigits[md5num[i] >> 4 & 0xf]
		k++
		str[k] = hexDigits[md5num[i] & 0xf]
		k++
	}

	return string(str)
}

/**
 * get token for file URL
 *
 * @param remote_filename the filename return by FastDFS server
 * @param ts              unix timestamp, unit: second
 * @param secret_key      the secret key
 * @return token string
 */
func GetToken(remoteFilename string, ts int, secretKey string) (string, error) {
	bsFilename,err := ConvertUTF8ToBytes([]byte(remoteFilename), GCharset)
	if err != nil {
		return "", err
	}
	bsKey,err := ConvertUTF8ToBytes([]byte(secretKey), GCharset)
	if err != nil {
		return "", err
	}
	bsTimestamp,err := ConvertUTF8ToBytes([]byte(strconv.Itoa(ts)), GCharset)
	if err != nil {
		return "", err
	}

	var buff = make([]byte, len(bsFilename) + len(bsKey) + len(bsTimestamp))
	copy(buff, bsFilename)
	copy(buff[len(bsFilename):], bsKey)
	copy(buff[len(bsFilename) + len(bsKey):], bsTimestamp)

	return Md5(buff), nil
}

/**
 * generate slave filename
 *
 * @param master_filename the master filename to generate the slave filename
 * @param prefix_name     the prefix name to generate the slave filename
 * @param ext_name        the extension name of slave filename, null for same as the master extension name
 * @return slave filename string
 */
func GenSlaveFilename(masterFilename, prefixName, extName string) (string, error) {
	var trueExtName string
	var dotIndex int

	if len(masterFilename) < 28 + FDFS_FILE_EXT_NAME_MAX_LEN {
		return "", fmt.Errorf("master filename \"%s\" is invalid", masterFilename)
	}

	var fromIndex = len(masterFilename) - (FDFS_FILE_EXT_NAME_MAX_LEN + 1)
	dotIndex = strings.IndexByte(masterFilename[fromIndex:], '.') + fromIndex
	if extName != "" {
		if len(extName) == 0 {
			return "", nil
		} else if extName[0] == '.' {
			trueExtName = extName
		} else {
			trueExtName = "." + extName
		}
	} else {
		if dotIndex < 0 {
			trueExtName = ""
		} else {
			trueExtName = masterFilename[dotIndex:]
		}
	}

	if len(trueExtName) == 0 && prefixName == "-m" {
		return "", fmt.Errorf("prefix_name \"%s\" is invalid", prefixName)
	}

	if dotIndex < 0 {
		return masterFilename + prefixName + trueExtName, nil
	}

	return masterFilename[:dotIndex] + prefixName + trueExtName, nil
}



/**
 * receive package info
 */
type RecvPackageInfo struct {
	Errno    byte
	Body     []byte
}

func NewRecvPackageInfo(errno byte, body []byte) *RecvPackageInfo {
	return &RecvPackageInfo{
		Errno:errno,
		Body:body,
	}
}


/**
 * receive header info
 */
type RecvHeaderInfo struct {
	Errno     byte
	BodyLen   int
}

func NewRecvHeaderInfo(errno byte, bodyLen int) *RecvHeaderInfo {
	return &RecvHeaderInfo{
		Errno:errno,
		BodyLen:bodyLen,
	}
}