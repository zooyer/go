package fastdfs

import (
	"time"
)

const (
	FIELD_INDEX_STATUS = 0
	FIELD_INDEX_ID = 1
	FIELD_INDEX_IP_ADDR = 2
	FIELD_INDEX_DOMAIN_NAME = 3
	FIELD_INDEX_SRC_IP_ADDR = 4
	FIELD_INDEX_VERSION = 5
	FIELD_INDEX_JOIN_TIME = 6
	FIELD_INDEX_UP_TIME = 7
	FIELD_INDEX_TOTAL_MB = 8
	FIELD_INDEX_FREE_MB = 9
	FIELD_INDEX_UPLOAD_PRIORITY = 10
	FIELD_INDEX_STORE_PATH_COUNT = 11
	FIELD_INDEX_SUBDIR_COUNT_PER_PATH = 12
	FIELD_INDEX_CURRENT_WRITE_PATH = 13
	FIELD_INDEX_STORAGE_PORT = 14
	FIELD_INDEX_STORAGE_HTTP_PORT = 15

	FIELD_INDEX_CONNECTION_ALLOC_COUNT = 16
	FIELD_INDEX_CONNECTION_CURRENT_COUNT = 17
	FIELD_INDEX_CONNECTION_MAX_COUNT = 18

	FIELD_INDEX_TOTAL_UPLOAD_COUNT = 19
	FIELD_INDEX_SUCCESS_UPLOAD_COUNT = 20
	FIELD_INDEX_TOTAL_APPEND_COUNT = 21
	FIELD_INDEX_SUCCESS_APPEND_COUNT = 22
	FIELD_INDEX_TOTAL_MODIFY_COUNT = 23
	FIELD_INDEX_SUCCESS_MODIFY_COUNT = 24
	FIELD_INDEX_TOTAL_TRUNCATE_COUNT = 25
	FIELD_INDEX_SUCCESS_TRUNCATE_COUNT = 26
	FIELD_INDEX_TOTAL_SET_META_COUNT = 27
	FIELD_INDEX_SUCCESS_SET_META_COUNT = 28
	FIELD_INDEX_TOTAL_DELETE_COUNT = 29
	FIELD_INDEX_SUCCESS_DELETE_COUNT = 30
	FIELD_INDEX_TOTAL_DOWNLOAD_COUNT = 31
	FIELD_INDEX_SUCCESS_DOWNLOAD_COUNT = 32
	FIELD_INDEX_TOTAL_GET_META_COUNT = 33
	FIELD_INDEX_SUCCESS_GET_META_COUNT = 34
	FIELD_INDEX_TOTAL_CREATE_LINK_COUNT = 35
	FIELD_INDEX_SUCCESS_CREATE_LINK_COUNT = 36
	FIELD_INDEX_TOTAL_DELETE_LINK_COUNT = 37
	FIELD_INDEX_SUCCESS_DELETE_LINK_COUNT = 38
	FIELD_INDEX_TOTAL_UPLOAD_BYTES = 39
	FIELD_INDEX_SUCCESS_UPLOAD_BYTES = 40
	FIELD_INDEX_TOTAL_APPEND_BYTES = 41
	FIELD_INDEX_SUCCESS_APPEND_BYTES = 42
	FIELD_INDEX_TOTAL_MODIFY_BYTES = 43
	FIELD_INDEX_SUCCESS_MODIFY_BYTES = 44
	FIELD_INDEX_TOTAL_DOWNLOAD_BYTES = 45
	FIELD_INDEX_SUCCESS_DOWNLOAD_BYTES = 46
	FIELD_INDEX_TOTAL_SYNC_IN_BYTES = 47
	FIELD_INDEX_SUCCESS_SYNC_IN_BYTES = 48
	FIELD_INDEX_TOTAL_SYNC_OUT_BYTES = 49
	FIELD_INDEX_SUCCESS_SYNC_OUT_BYTES = 50
	FIELD_INDEX_TOTAL_FILE_OPEN_COUNT = 51
	FIELD_INDEX_SUCCESS_FILE_OPEN_COUNT = 52
	FIELD_INDEX_TOTAL_FILE_READ_COUNT = 53
	FIELD_INDEX_SUCCESS_FILE_READ_COUNT = 54
	FIELD_INDEX_TOTAL_FILE_WRITE_COUNT = 55
	FIELD_INDEX_SUCCESS_FILE_WRITE_COUNT = 56
	FIELD_INDEX_LAST_SOURCE_UPDATE = 57
	FIELD_INDEX_LAST_SYNC_UPDATE = 58
	FIELD_INDEX_LAST_SYNCED_TIMESTAMP = 59
	FIELD_INDEX_LAST_HEART_BEAT_TIME = 60
	FIELD_INDEX_IF_TRUNK_FILE = 61
)

var storageFieldsTotalSize int
var storageFieldsArray = make([]*FieldInfo, 62)

func init() {
	var offset = 0

	storageFieldsArray[FIELD_INDEX_STATUS] = NewFieldInfo("status", offset, 1)
	offset += 1

	storageFieldsArray[FIELD_INDEX_ID] = NewFieldInfo("id", offset, FDFS_STORAGE_ID_MAX_SIZE)
	offset += FDFS_STORAGE_ID_MAX_SIZE

	storageFieldsArray[FIELD_INDEX_IP_ADDR] = NewFieldInfo("ipAddr", offset, FDFS_IPADDR_SIZE)
	offset += FDFS_IPADDR_SIZE

	storageFieldsArray[FIELD_INDEX_DOMAIN_NAME] = NewFieldInfo("domainName", offset, FDFS_DOMAIN_NAME_MAX_SIZE)
	offset += FDFS_DOMAIN_NAME_MAX_SIZE

	storageFieldsArray[FIELD_INDEX_SRC_IP_ADDR] = NewFieldInfo("srcIpAddr", offset, FDFS_IPADDR_SIZE)
	offset += FDFS_IPADDR_SIZE

	storageFieldsArray[FIELD_INDEX_VERSION] = NewFieldInfo("version", offset, FDFS_VERSION_SIZE)
	offset += FDFS_VERSION_SIZE

	storageFieldsArray[FIELD_INDEX_JOIN_TIME] = NewFieldInfo("joinTime", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_UP_TIME] = NewFieldInfo("upTime", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_MB] = NewFieldInfo("totalMB", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_FREE_MB] = NewFieldInfo("freeMB", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_UPLOAD_PRIORITY] = NewFieldInfo("uploadPriority", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_STORE_PATH_COUNT] = NewFieldInfo("storePathCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUBDIR_COUNT_PER_PATH] = NewFieldInfo("subdirCountPerPath", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_CURRENT_WRITE_PATH] = NewFieldInfo("currentWritePath", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_STORAGE_PORT] = NewFieldInfo("storagePort", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_STORAGE_HTTP_PORT] = NewFieldInfo("storageHttpPort", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_CONNECTION_ALLOC_COUNT] = NewFieldInfo("connectionAllocCount", offset, 4)
	offset += 4

	storageFieldsArray[FIELD_INDEX_CONNECTION_CURRENT_COUNT] = NewFieldInfo("connectionCurrentCount", offset, 4)
	offset += 4

	storageFieldsArray[FIELD_INDEX_CONNECTION_MAX_COUNT] = NewFieldInfo("connectionMaxCount", offset, 4)
	offset += 4

	storageFieldsArray[FIELD_INDEX_TOTAL_UPLOAD_COUNT] = NewFieldInfo("totalUploadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_UPLOAD_COUNT] = NewFieldInfo("successUploadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_APPEND_COUNT] = NewFieldInfo("totalAppendCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_APPEND_COUNT] = NewFieldInfo("successAppendCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_MODIFY_COUNT] = NewFieldInfo("totalModifyCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_MODIFY_COUNT] = NewFieldInfo("successModifyCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_TRUNCATE_COUNT] = NewFieldInfo("totalTruncateCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_TRUNCATE_COUNT] = NewFieldInfo("successTruncateCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_SET_META_COUNT] = NewFieldInfo("totalSetMetaCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_SET_META_COUNT] = NewFieldInfo("successSetMetaCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_DELETE_COUNT] = NewFieldInfo("totalDeleteCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_DELETE_COUNT] = NewFieldInfo("successDeleteCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_DOWNLOAD_COUNT] = NewFieldInfo("totalDownloadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_DOWNLOAD_COUNT] = NewFieldInfo("successDownloadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_GET_META_COUNT] = NewFieldInfo("totalGetMetaCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_GET_META_COUNT] = NewFieldInfo("successGetMetaCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_CREATE_LINK_COUNT] = NewFieldInfo("totalCreateLinkCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_CREATE_LINK_COUNT] = NewFieldInfo("successCreateLinkCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_DELETE_LINK_COUNT] = NewFieldInfo("totalDeleteLinkCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_DELETE_LINK_COUNT] = NewFieldInfo("successDeleteLinkCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_UPLOAD_BYTES] = NewFieldInfo("totalUploadBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_UPLOAD_BYTES] = NewFieldInfo("successUploadBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_APPEND_BYTES] = NewFieldInfo("totalAppendBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_APPEND_BYTES] = NewFieldInfo("successAppendBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_MODIFY_BYTES] = NewFieldInfo("totalModifyBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_MODIFY_BYTES] = NewFieldInfo("successModifyBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_DOWNLOAD_BYTES] = NewFieldInfo("totalDownloadloadBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_DOWNLOAD_BYTES] = NewFieldInfo("successDownloadloadBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_SYNC_IN_BYTES] = NewFieldInfo("totalSyncInBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_SYNC_IN_BYTES] = NewFieldInfo("successSyncInBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_SYNC_OUT_BYTES] = NewFieldInfo("totalSyncOutBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_SYNC_OUT_BYTES] = NewFieldInfo("successSyncOutBytes", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_FILE_OPEN_COUNT] = NewFieldInfo("totalFileOpenCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_OPEN_COUNT] = NewFieldInfo("successFileOpenCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_FILE_READ_COUNT] = NewFieldInfo("totalFileReadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_READ_COUNT] = NewFieldInfo("successFileReadCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_TOTAL_FILE_WRITE_COUNT] = NewFieldInfo("totalFileWriteCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_WRITE_COUNT] = NewFieldInfo("successFileWriteCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_LAST_SOURCE_UPDATE] = NewFieldInfo("lastSourceUpdate", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_LAST_SYNC_UPDATE] = NewFieldInfo("lastSyncUpdate", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_LAST_SYNCED_TIMESTAMP] = NewFieldInfo("lastSyncedTimestamp", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_LAST_HEART_BEAT_TIME] = NewFieldInfo("lastHeartBeatTime", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	storageFieldsArray[FIELD_INDEX_IF_TRUNK_FILE] = NewFieldInfo("ifTrunkServer", offset, 1)
	offset += 1

	storageFieldsTotalSize = offset
}


type StructStorageStat struct {
	StructBase

	status byte
	id string
	ipAddr string
	srcIpAddr string
	domainName string //http domain name
	version string
	totalMB int64 //total disk storage in MB
	freeMB int64 //free disk storage in MB
	uploadPriority int  //upload priority
	joinTime time.Time //storage join timestamp (create time
	upTime time.Time   //storage service started timestamp
	storePathCount int  //store base path count of each
	subdirCountPerPath int
	storagePort int
	storageHttpPort int //storage http server port
	currentWritePath int //current write path index
	connectionAllocCount int
	connectionCurrentCount int
	connectionMaxCount int
	totalUploadCount int64
	successUploadCount int64
	totalAppendCount int64
	successAppendCount int64
	totalModifyCount int64
	successModifyCount int64
	totalTruncateCount int64
	successTruncateCount int64
	totalSetMetaCount int64
	successSetMetaCount int64
	totalDeleteCount int64
	successDeleteCount int64
	totalDownloadCount int64
	successDownloadCount int64
	totalGetMetaCount int64
	successGetMetaCount int64
	totalCreateLinkCount int64
	successCreateLinkCount int64
	totalDeleteLinkCount int64
	successDeleteLinkCount int64
	totalUploadBytes int64
	successUploadBytes int64
	totalAppendBytes int64
	successAppendBytes int64
	totalModifyBytes int64
	successModifyBytes int64
	totalDownloadloadBytes int64
	successDownloadloadBytes int64
	totalSyncInBytes int64
	successSyncInBytes int64
	totalSyncOutBytes int64
	successSyncOutBytes int64
	totalFileOpenCount int64
	successFileOpenCount int64
	totalFileReadCount int64
	successFileReadCount int64
	totalFileWriteCount int64
	successFileWriteCount int64
	lastSourceUpdate time.Time
	lastSyncUpdate time.Time
	lastSyncedTimestamp time.Time
	lastHeartBeatTime time.Time
	ifTrunkServer bool
}

/**
 * get fields total size
 *
 * @return fields total size
 */
func GetStorageFieldsTotalSize() int {
	return storageFieldsTotalSize
}

/**
 * get storage status
 *
 * @return storage status
 */
func (s *StructStorageStat) GetStatus() byte {
	return s.status
}

/**
 * get storage server id
 *
 * @return storage server id
 */
func (s *StructStorageStat) GetId() string {
	return s.id
}

/**
 * get storage server ip address
 *
 * @return storage server ip address
 */
func (s *StructStorageStat) GetIpAddr() string {
	return s.ipAddr
}

/**
 * get source storage ip address
 *
 * @return source storage ip address
 */
func (s *StructStorageStat) GetSrcIpAddr() string {
	return s.srcIpAddr
}

/**
 * get the domain name of the storage server
 *
 * @return the domain name of the storage server
 */
func (s *StructStorageStat) GetDomainName() string {
	return s.domainName
}

/**
 * get storage version
 *
 * @return storage version
 */
func (s *StructStorageStat) GetVersion() string {
	return s.version
}

/**
 * get total disk space in MB
 *
 * @return total disk space in MB
 */
func (s *StructStorageStat) GetTotalMB() int64 {
	return s.totalMB
}

/**
 * get free disk space in MB
 *
 * @return free disk space in MB
 */
func (s *StructStorageStat) GetFreeMB() int64 {
	return s.freeMB
}

/**
 * get storage server upload priority
 *
 * @return storage server upload priority
 */
func (s *StructStorageStat) GetUploadPriority() int {
	return s.uploadPriority
}

/**
 * get storage server join time
 *
 * @return storage server join time
 */
func (s *StructStorageStat) GetJoinTime() time.Time {
	return s.joinTime
}

/**
 * get storage server up time
 *
 * @return storage server up time
 */
func (s *StructStorageStat) GetUpTime() time.Time {
	return s.upTime
}

/**
 * get store base path count of each storage server
 *
 * @return store base path count of each storage server
 */
func (s *StructStorageStat) GetStorePathCount() int {
	return s.storePathCount
}

/**
 * get sub dir count per store path
 *
 * @return sub dir count per store path
 */
func (s *StructStorageStat) GetSubdirCountPerPath() int {
	return s.subdirCountPerPath
}

/**
 * get storage server port
 *
 * @return storage server port
 */
func (s *StructStorageStat) GetStoragePort() int {
	return s.storagePort
}

/**
 * get storage server HTTP port
 *
 * @return storage server HTTP port
 */
func (s *StructStorageStat) GetStorageHttpPort() int {
	return s.storageHttpPort
}

/**
 * get current write path index
 *
 * @return current write path index
 */
func (s *StructStorageStat) GetCurrentWritePath() int {
	return s.currentWritePath
}

/**
 * get total upload file count
 *
 * @return total upload file count
 */
func (s *StructStorageStat) GetTotalUploadCount() int64 {
	return s.totalUploadCount
}

/**
 * get success upload file count
 *
 * @return success upload file count
 */
func (s *StructStorageStat) GetSuccessUploadCount() int64 {
	return s.successUploadCount
}

/**
 * get total append count
 *
 * @return total append count
 */
func (s *StructStorageStat) GetTotalAppendCount() int64 {
	return s.totalAppendCount
}

/**
 * get success append count
 *
 * @return success append count
 */
func (s *StructStorageStat) GetSuccessAppendCount() int64 {
	return s.successAppendCount
}

/**
 * get total modify count
 *
 * @return total modify count
 */
func (s *StructStorageStat) GetTotalModifyCount() int64 {
	return s.totalModifyCount
}

/**
 * get success modify count
 *
 * @return success modify count
 */
func (s *StructStorageStat) GetSuccessModifyCount() int64 {
	return s.successModifyCount
}

/**
 * get total truncate count
 *
 * @return total truncate count
 */
func (s *StructStorageStat) GetTotalTruncateCount() int64 {
	return s.totalTruncateCount
}

/**
 * get success truncate count
 *
 * @return success truncate count
 */
func (s *StructStorageStat) GetSuccessTruncateCount() int64 {
	return s.successTruncateCount
}

/**
 * get total set meta data count
 *
 * @return total set meta data count
 */
func (s *StructStorageStat) GetTotalSetMetaCount() int64 {
	return s.totalSetMetaCount
}

/**
 * get success set meta data count
 *
 * @return success set meta data count
 */
func (s *StructStorageStat) GetSuccessSetMetaCount() int64 {
	return s.successSetMetaCount
}

/**
 * get total delete file count
 *
 * @return total delete file count
 */
func (s *StructStorageStat) GetTotalDeleteCount() int64 {
	return s.totalDeleteCount
}

/**
 * get success delete file count
 *
 * @return success delete file count
 */
func (s *StructStorageStat) GetSuccessDeleteCount() int64 {
	return s.successDeleteCount
}

/**
 * get total download file count
 *
 * @return total download file count
 */
func (s *StructStorageStat) GetTotalDownloadCount() int64 {
	return s.totalDownloadCount
}

/**
 * get success download file count
 *
 * @return success download file count
 */
func (s *StructStorageStat) GetSuccessDownloadCount() int64 {
	return s.successDownloadCount
}

/**
 * get total get metadata count
 *
 * @return total get metadata count
 */
func (s *StructStorageStat) GetTotalGetMetaCount() int64 {
	return s.totalGetMetaCount
}

/**
 * get success get metadata count
 *
 * @return success get metadata count
 */
func (s *StructStorageStat) GetSuccessGetMetaCount() int64 {
	return s.successGetMetaCount
}

/**
 * get total create linke count
 *
 * @return total create linke count
 */
func (s *StructStorageStat) GetTotalCreateLinkCount() int64 {
	return s.totalCreateLinkCount
}

/**
 * get success create linke count
 *
 * @return success create linke count
 */
func (s *StructStorageStat) GetSuccessCreateLinkCount() int64 {
	return s.successCreateLinkCount
}

/**
 * get total delete link count
 *
 * @return total delete link count
 */
func (s *StructStorageStat) GetTotalDeleteLinkCount() int64 {
	return s.totalDeleteLinkCount
}

/**
 * get success delete link count
 *
 * @return success delete link count
 */
func (s *StructStorageStat) GetSuccessDeleteLinkCount() int64 {
	return s.successDeleteLinkCount
}

/**
 * get total upload file bytes
 *
 * @return total upload file bytes
 */
func (s *StructStorageStat) GetTotalUploadBytes() int64 {
	return s.totalUploadBytes
}

/**
 * get success upload file bytes
 *
 * @return success upload file bytes
 */
func (s *StructStorageStat) GetSuccessUploadBytes() int64 {
	return s.successUploadBytes
}

/**
 * get total append bytes
 *
 * @return total append bytes
 */
func (s *StructStorageStat) GetTotalAppendBytes() int64 {
	return s.totalAppendBytes
}

/**
 * get success append bytes
 *
 * @return success append bytes
 */
func (s *StructStorageStat) GetSuccessAppendBytes() int64 {
	return s.successAppendBytes
}

/**
 * get total modify bytes
 *
 * @return total modify bytes
 */
func (s *StructStorageStat) GetTotalModifyBytes() int64 {
	return s.totalModifyBytes
}

/**
 * get success modify bytes
 *
 * @return success modify bytes
 */
func (s *StructStorageStat) GetSuccessModifyBytes() int64 {
	return s.successModifyBytes
}

/**
 * get total download file bytes
 *
 * @return total download file bytes
 */
func (s *StructStorageStat) GetTotalDownloadloadBytes() int64 {
	return s.totalDownloadloadBytes
}

/**
 * get success download file bytes
 *
 * @return success download file bytes
 */
func (s *StructStorageStat) GetSuccessDownloadloadBytes() int64 {
	return s.successDownloadloadBytes
}

/**
 * get total sync in bytes
 *
 * @return total sync in bytes
 */
func (s *StructStorageStat) GetTotalSyncInBytes() int64 {
	return s.totalSyncInBytes
}

/**
 * get success sync in bytes
 *
 * @return success sync in bytes
 */
func (s *StructStorageStat) GetSuccessSyncInBytes() int64 {
	return s.successSyncInBytes
}

/**
 * get total sync out bytes
 *
 * @return total sync out bytes
 */
func (s *StructStorageStat) GetTotalSyncOutBytes() int64 {
	return s.totalSyncOutBytes
}

/**
 * get success sync out bytes
 *
 * @return success sync out bytes
 */
func (s *StructStorageStat) GetSuccessSyncOutBytes() int64 {
	return s.successSyncOutBytes
}

/**
 * get total file opened count
 *
 * @return total file opened bytes
 */
func (s *StructStorageStat) GetTotalFileOpenCount() int64 {
	return s.totalFileOpenCount
}

/**
 * get success file opened count
 *
 * @return success file opened count
 */
func (s *StructStorageStat) GetSuccessFileOpenCount() int64 {
	return s.successFileOpenCount
}

/**
 * get total file read count
 *
 * @return total file read bytes
 */
func (s *StructStorageStat) GetTotalFileReadCount() int64 {
	return s.totalFileReadCount
}

/**
 * get success file read count
 *
 * @return success file read count
 */
func (s *StructStorageStat) GetSuccessFileReadCount() int64 {
	return s.successFileReadCount
}

/**
 * get total file write count
 *
 * @return total file write bytes
 */
func (s *StructStorageStat) GetTotalFileWriteCount() int64 {
	return s.totalFileWriteCount
}

/**
 * get success file write count
 *
 * @return success file write count
 */
func (s *StructStorageStat) GetSuccessFileWriteCount() int64 {
	return s.successFileWriteCount
}

/**
 * get last source update timestamp
 *
 * @return last source update timestamp
 */
func (s *StructStorageStat) GetLastSourceUpdate() time.Time {
	return s.lastSourceUpdate
}

/**
 * get last synced update timestamp
 *
 * @return last synced update timestamp
 */
func (s *StructStorageStat) GetLastSyncUpdate() time.Time {
	return s.lastSyncUpdate
}

/**
 * get last synced timestamp
 *
 * @return last synced timestamp
 */
func (s *StructStorageStat) GetLastSyncedTimestamp() time.Time {
	return s.lastSyncedTimestamp
}

/**
 * get last heart beat timestamp
 *
 * @return last heart beat timestamp
 */
func (s *StructStorageStat) GetLastHeartBeatTime() time.Time {
	return s.lastHeartBeatTime
}

/**
 * if the trunk server
 *
 * @return true for the trunk server, otherwise false
 */
func (s *StructStorageStat) isTrunkServer() bool {
	return s.ifTrunkServer
}

/**
 * get connection alloc count
 *
 * @return connection alloc count
 */
func (s *StructStorageStat) GetConnectionAllocCount() int {
	return s.connectionAllocCount
}

/**
 * get connection current count
 *
 * @return connection current count
 */
func (s *StructStorageStat) GetConnectionCurrentCount() int {
	return s.connectionCurrentCount
}

/**
 * get connection max count
 *
 * @return connection max count
 */
func (s *StructStorageStat) GetConnectionMaxCount() int {
	return s.connectionMaxCount
}

/**
 * set fields
 *
 * @param bs     byte array
 * @param offset start offset
 */
func (s *StructStorageStat) SetFields(bs []byte, offset int) {
	s.status = s.StructBase.byteValue(bs, offset, storageFieldsArray[FIELD_INDEX_STATUS]);
	s.id = s.StructBase.stringValue(bs, offset, storageFieldsArray[FIELD_INDEX_ID]);
	s.ipAddr = s.StructBase.stringValue(bs, offset, storageFieldsArray[FIELD_INDEX_IP_ADDR]);
	s.srcIpAddr = s.StructBase.stringValue(bs, offset, storageFieldsArray[FIELD_INDEX_SRC_IP_ADDR]);
	s.domainName = s.StructBase.stringValue(bs, offset, storageFieldsArray[FIELD_INDEX_DOMAIN_NAME]);
	s.version = s.StructBase.stringValue(bs, offset, storageFieldsArray[FIELD_INDEX_VERSION]);
	s.totalMB = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_MB]);
	s.freeMB = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_FREE_MB]);
	s.uploadPriority = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_UPLOAD_PRIORITY]);
	s.joinTime = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_JOIN_TIME]);
	s.upTime = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_UP_TIME]);
	s.storePathCount = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_STORE_PATH_COUNT]);
	s.subdirCountPerPath = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUBDIR_COUNT_PER_PATH]);
	s.storagePort = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_STORAGE_PORT]);
	s.storageHttpPort = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_STORAGE_HTTP_PORT]);
	s.currentWritePath = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_CURRENT_WRITE_PATH]);

	s.connectionAllocCount = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_CONNECTION_ALLOC_COUNT]);
	s.connectionCurrentCount = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_CONNECTION_CURRENT_COUNT]);
	s.connectionMaxCount = s.StructBase.intValue(bs, offset, storageFieldsArray[FIELD_INDEX_CONNECTION_MAX_COUNT]);

	s.totalUploadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_UPLOAD_COUNT]);
	s.successUploadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_UPLOAD_COUNT]);
	s.totalAppendCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_APPEND_COUNT]);
	s.successAppendCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_APPEND_COUNT]);
	s.totalModifyCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_MODIFY_COUNT]);
	s.successModifyCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_MODIFY_COUNT]);
	s.totalTruncateCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_TRUNCATE_COUNT]);
	s.successTruncateCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_TRUNCATE_COUNT]);
	s.totalSetMetaCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_SET_META_COUNT]);
	s.successSetMetaCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_SET_META_COUNT]);
	s.totalDeleteCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_DELETE_COUNT]);
	s.successDeleteCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_DELETE_COUNT]);
	s.totalDownloadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_DOWNLOAD_COUNT]);
	s.successDownloadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_DOWNLOAD_COUNT]);
	s.totalGetMetaCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_GET_META_COUNT]);
	s.successGetMetaCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_GET_META_COUNT]);
	s.totalCreateLinkCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_CREATE_LINK_COUNT]);
	s.successCreateLinkCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_CREATE_LINK_COUNT]);
	s.totalDeleteLinkCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_DELETE_LINK_COUNT]);
	s.successDeleteLinkCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_DELETE_LINK_COUNT]);
	s.totalUploadBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_UPLOAD_BYTES]);
	s.successUploadBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_UPLOAD_BYTES]);
	s.totalAppendBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_APPEND_BYTES]);
	s.successAppendBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_APPEND_BYTES]);
	s.totalModifyBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_MODIFY_BYTES]);
	s.successModifyBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_MODIFY_BYTES]);
	s.totalDownloadloadBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_DOWNLOAD_BYTES]);
	s.successDownloadloadBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_DOWNLOAD_BYTES]);
	s.totalSyncInBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_SYNC_IN_BYTES]);
	s.successSyncInBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_SYNC_IN_BYTES]);
	s.totalSyncOutBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_SYNC_OUT_BYTES]);
	s.successSyncOutBytes = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_SYNC_OUT_BYTES]);
	s.totalFileOpenCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_FILE_OPEN_COUNT]);
	s.successFileOpenCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_OPEN_COUNT]);
	s.totalFileReadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_FILE_READ_COUNT]);
	s.successFileReadCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_READ_COUNT]);
	s.totalFileWriteCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_TOTAL_FILE_WRITE_COUNT]);
	s.successFileWriteCount = s.StructBase.longValue(bs, offset, storageFieldsArray[FIELD_INDEX_SUCCESS_FILE_WRITE_COUNT]);
	s.lastSourceUpdate = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_LAST_SOURCE_UPDATE]);
	s.lastSyncUpdate = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_LAST_SYNC_UPDATE]);
	s.lastSyncedTimestamp = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_LAST_SYNCED_TIMESTAMP]);
	s.lastHeartBeatTime = s.StructBase.dateValue(bs, offset, storageFieldsArray[FIELD_INDEX_LAST_HEART_BEAT_TIME]);
	s.ifTrunkServer = s.StructBase.boolValue(bs, offset, storageFieldsArray[FIELD_INDEX_IF_TRUNK_FILE]);
}