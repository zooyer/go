package fastdfs

const (
	GROUP_FIELD_INDEX_GROUP_NAME = 0
	GROUP_FIELD_INDEX_TOTAL_MB = 1
	GROUP_FIELD_INDEX_FREE_MB = 2
	GROUP_FIELD_INDEX_TRUNK_FREE_MB = 3
	GROUP_FIELD_INDEX_STORAGE_COUNT = 4
	GROUP_FIELD_INDEX_STORAGE_PORT = 5
	GROUP_FIELD_INDEX_STORAGE_HTTP_PORT = 6
	GROUP_FIELD_INDEX_ACTIVE_COUNT = 7
	GROUP_FIELD_INDEX_CURRENT_WRITE_SERVER = 8
	GROUP_FIELD_INDEX_STORE_PATH_COUNT = 9
	GROUP_FIELD_INDEX_SUBDIR_COUNT_PER_PATH = 10
	GROUP_FIELD_INDEX_CURRENT_TRUNK_FILE_ID = 11
)

var groupFieldsTotalSize int
var groupFieldsArray = make([]*FieldInfo, 12)

func init() {
	var offset = 0
	groupFieldsArray[GROUP_FIELD_INDEX_GROUP_NAME] = NewFieldInfo("groupName", offset, FDFS_GROUP_NAME_MAX_LEN + 1)
	offset += FDFS_GROUP_NAME_MAX_LEN + 1

	groupFieldsArray[GROUP_FIELD_INDEX_TOTAL_MB] = NewFieldInfo("totalMB", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_FREE_MB] = NewFieldInfo("freeMB", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_TRUNK_FREE_MB] = NewFieldInfo("trunkFreeMB", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_COUNT] = NewFieldInfo("storageCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_PORT] = NewFieldInfo("storagePort", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_HTTP_PORT] = NewFieldInfo("storageHttpPort", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_ACTIVE_COUNT] = NewFieldInfo("activeCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_CURRENT_WRITE_SERVER] = NewFieldInfo("currentWriteServer", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_STORE_PATH_COUNT] = NewFieldInfo("storePathCount", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_SUBDIR_COUNT_PER_PATH] = NewFieldInfo("subdirCountPerPath", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsArray[GROUP_FIELD_INDEX_CURRENT_TRUNK_FILE_ID] = NewFieldInfo("currentTrunkFileId", offset, FDFS_PROTO_PKG_LEN_SIZE)
	offset += FDFS_PROTO_PKG_LEN_SIZE

	groupFieldsTotalSize = offset
}

type StructGroupStat struct {
	StructBase

	groupName string        //name of this group
	totalMB int64           //total disk storage in MB
	freeMB int64            //free disk space in MB
	trunkFreeMB int64       //trunk free space in MB
	storageCount int        //storage server count
	storagePort int         //storage server port
	storageHttpPort int     //storage server HTTP port
	activeCount int         //active storage server count
	currentWriteServer int  //current storage server index to upload file
	storePathCount int      //store base path count of each storage server
	subdirCountPerPath int  //sub dir count per store path
	currentTrunkFileId int  //current trunk file id
}

/**
 * get fields total size
 *
 * @return fields total size
 */
func GetGroupFieldsTotalSize() int {
	return groupFieldsTotalSize
}

/**
 * get group name
 *
 * @return group name
 */
func (s *StructGroupStat) GetGroupName() string {
	return s.groupName
}

/**
 * get total disk space in MB
 *
 * @return total disk space in MB
 */
func (s *StructGroupStat) GetTotalMB() int64 {
	return s.totalMB
}

/**
 * get free disk space in MB
 *
 * @return free disk space in MB
 */
func (s *StructGroupStat) GetFreeMB() int64 {
	return s.freeMB
}

/**
 * get trunk free space in MB
 *
 * @return trunk free space in MB
 */
func (s *StructGroupStat) GetTrunkFreeMB() int64 {
	return s.trunkFreeMB
}

/**
 * get storage server count in this group
 *
 * @return storage server count in this group
 */
func (s *StructGroupStat) GetStorageCount() int {
	return s.storageCount
}

/**
 * get active storage server count in this group
 *
 * @return active storage server count in this group
 */
func (s *StructGroupStat) GetActiveCount() int {
	return s.activeCount
}

/**
 * get storage server port
 *
 * @return storage server port
 */
func (s *StructGroupStat) GetStoragePort() int {
	return s.storagePort
}

/**
 * get storage server HTTP port
 *
 * @return storage server HTTP port
 */
func (s *StructGroupStat) GetStorageHttpPort() int {
	return s.storageHttpPort
}

/**
 * get current storage server index to upload file
 *
 * @return current storage server index to upload file
 */
func (s *StructGroupStat) GetCurrentWriteServer() int {
	return s.currentWriteServer
}

/**
 * get store base path count of each storage server
 *
 * @return store base path count of each storage server
 */
func (s *StructGroupStat) GetStorePathCount() int {
	return s.storePathCount
}

/**
 * get sub dir count per store path
 *
 * @return sub dir count per store path
 */
func (s *StructGroupStat) GetSubdirCountPerPath() int {
	return s.subdirCountPerPath
}

/**
 * get current trunk file id
 *
 * @return current trunk file id
 */
func (s *StructGroupStat) GetCurrentTrunkFileId() int {
	return s.currentTrunkFileId
}

/**
 * set fields
 *
 * @param bs     byte array
 * @param offset start offset
 */
func (s *StructGroupStat) SetFields(bs []byte, offset int) {
	s.groupName = s.StructBase.stringValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_GROUP_NAME])
	s.totalMB = s.StructBase.longValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_TOTAL_MB])
	s.freeMB = s.StructBase.longValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_FREE_MB])
	s.trunkFreeMB = s.StructBase.longValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_TRUNK_FREE_MB])
	s.storageCount = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_COUNT])
	s.storagePort = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_PORT])
	s.storageHttpPort = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_STORAGE_HTTP_PORT])
	s.activeCount = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_ACTIVE_COUNT])
	s.currentWriteServer = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_CURRENT_WRITE_SERVER])
	s.storePathCount = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_STORE_PATH_COUNT])
	s.subdirCountPerPath = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_SUBDIR_COUNT_PER_PATH])
	s.currentTrunkFileId = s.StructBase.intValue(bs, offset, groupFieldsArray[GROUP_FIELD_INDEX_CURRENT_TRUNK_FILE_ID])
}
