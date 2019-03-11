package fastdfs

import "strings"

const SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR = "/"

/**
 * Storage client for 1 field file id: combined group name and filename
 *
 * @author Happy Fish / YuQing
 * @version Version 1.21
 */
type StorageClient1 struct {
	StorageClient
}

/**
 * constructor
 *
 * @param trackerServer the tracker server, can be null
 * @param storageServer the storage server, can be null
 */
func NewStorageClient1(trackerServer *TrackerServer, storageServer *StorageServer) *StorageClient1 {
	return &StorageClient1{
		StorageClient:*NewStorageClientByServer(trackerServer, storageServer),
	}
}

func SplitFileId(fileId string, results []string) byte {
	var pos = strings.Index(fileId, SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR)
	if pos <= 0 || pos == len(fileId) - 1 {
		return ERR_NO_EINVAL
	}

	results[0] = fileId[:pos] //group name
	results[1] = fileId[pos + 1:] //file name

	return 0
}