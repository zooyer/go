package fastdfs

import (
	"time"
	"strconv"
)

type FileInfo struct {
	sourceIpAddr        string
	fileSize            int64
	createTimeStamp     time.Time
	crc32               int
}

/**
 * Constructor
 *
 * @param file_size        the file size
 * @param create_timestamp create timestamp in seconds
 * @param crc32            the crc32 signature
 * @param source_ip_addr   the source storage ip address
 */
func NewFileInfo(fileSize int64, createTimestamp int64, crc32 int, sourceIpAddr string) *FileInfo {
	return &FileInfo{
		fileSize        : fileSize,
		createTimeStamp : time.Unix(createTimestamp, 0),
		crc32           : crc32,
		sourceIpAddr    : sourceIpAddr,
	}
}

/**
 * get the source ip address of the file uploaded to
 *
 * @return the source ip address of the file uploaded to
 */
func (f *FileInfo) GetSourceIpAddr() string {
	return f.sourceIpAddr
}

/**
 * set the source ip address of the file uploaded to
 *
 * @param source_ip_addr the source ip address
 */
func (f *FileInfo) SetSourceIpAddr(sourceIpAddr string) {
	f.sourceIpAddr = sourceIpAddr
}

/**
 * get the file size
 *
 * @return the file size
 */
func (f *FileInfo) GetFileSize() int64 {
	return f.fileSize
}

/**
 * set the file size
 *
 * @param file_size the file size
 */
func (f *FileInfo) SetFileSize(fileSize int64) {
	f.fileSize = fileSize
}

/**
 * get the create timestamp of the file
 *
 * @return the create timestamp of the file
 */
func (f *FileInfo) GetCreateTimestamp() time.Time {
	return f.createTimeStamp
}

/**
 * set the create timestamp of the file
 *
 * @param create_timestamp create timestamp in seconds
 */
func (f *FileInfo) SetCreateTimestamp(createTimestamp int64) {
	f.createTimeStamp = time.Unix(createTimestamp, 0)
}

/**
 * get the file CRC32 signature
 *
 * @return the file CRC32 signature
 */
func (f *FileInfo) GetCrc32() int {
	return f.crc32
}

/**
 * set the create timestamp of the file
 *
 * @param crc32 the crc32 signature
 */
func (f *FileInfo) SetCrc32(crc32 int) {
	f.crc32 = crc32
}

/**
 * to string
 *
 * @return string
 */
func (f *FileInfo) ToString() string {
	return f.String()
}

func (f *FileInfo) String() string {
	var df = "2006-01-02 15:04:05"
	return "source_ip_addr = " + f.sourceIpAddr + ", " +
		"file_size = " + strconv.Itoa(int(f.fileSize)) + ", " +
		"create_timestamp = " + f.createTimeStamp.Format(df) + ", " +
		"crc32 = " + strconv.Itoa(f.crc32)
}