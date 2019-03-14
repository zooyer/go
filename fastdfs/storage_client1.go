package fastdfs

import (
	"strings"
	"fmt"
)

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

/**
 * upload file to storage server (by file name)
 *
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadFile1(localFilename, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadFile(localFilename, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file name)
 *
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadFileByGroup1(groupName, localFilename, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.uploadFileByGroup(groupName, localFilename, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file name)
 *
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadBuffer1(fileBuff []byte, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadBuffer(fileBuff, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadBufferByGroup1(groupName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadBufferByGroup(groupName, fileBuff, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by callback)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_size     the file size
 * @param callback      the write data callback object
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadCallback1(groupName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadCallback(groupName, fileSize, callback, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload appender file to storage server (by file name)
 *
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadAppenderFile1(localFilename, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadAppenderFile(localFilename, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload appender file to storage server (by file name)
 *
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadAppenderFileByGroup1(groupName, localFilename, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadAppenderFileByGroup(groupName, localFilename, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadAppenderBuffer1(fileBuff []byte, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadAppenderBuffer(fileBuff, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadAppenderBufferByGroup1(groupName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadAppenderBufferByGroup(groupName, fileBuff, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload appender file to storage server (by callback)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_size     the file size
 * @param callback      the write data callback object
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadAppenderCallback1(groupName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts,err = s.UploadAppenderCallback(groupName, fileSize, callback, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file name, slave file mode)
 *
 * @param master_file_id the master file id to generate the slave file
 * @param prefix_name    the prefix name to generate the slave file
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadMasterFile1(masterFileId, prefixName, localFilename, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts = make([]string, 2)
	var err error
	s.errno = SplitFileId(masterFileId, parts)
	if s.errno != 0 {
		return "", nil
	}
	parts,err = s.UploadMasterFile(parts[0], parts[1], prefixName, localFilename, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file buff, slave file mode)
 *
 * @param master_file_id the master file id to generate the slave file
 * @param prefix_name    the prefix name to generate the slave file
 * @param file_buff      file content/buff
 * @param file_ext_name  file ext name, do not include dot(.)
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadMasterBuffer1(masterFileId, prefixName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts = make([]string, 2)
	var err error
	s.errno = SplitFileId(masterFileId, parts)
	if s.errno != 0 {
		return "", nil
	}
	parts,err = s.UploadMasterBuffer(parts[0], parts[1], prefixName, fileBuff, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by file buff, slave file mode)
 *
 * @param master_file_id the master file id to generate the slave file
 * @param prefix_name    the prefix name to generate the slave file
 * @param file_buff      file content/buff
 * @param file_ext_name  file ext name, do not include dot(.)
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadMasterOffsetBuffer1(groupName, masterFileId, prefixName string, fileBuff []byte, offset, length int, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts = make([]string, 2)
	var err error
	s.errno = SplitFileId(masterFileId, parts)
	if s.errno != 0 {
		return "", nil
	}
	parts,err = s.UploadMasterOffsetBuffer(parts[0], parts[1], prefixName, fileBuff, offset, length, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * upload file to storage server (by callback)
 *
 * @param master_file_id the master file id to generate the slave file
 * @param prefix_name    the prefix name to generate the slave file
 * @param file_size      the file size
 * @param callback       the write data callback object
 * @param file_ext_name  file ext name, do not include dot(.)
 * @param meta_list      meta info array
 * @return file id(including group name and filename) if success, <br>
 * return null if fail
 */
func (s *StorageClient1) UploadMasterCallback1(groupName, masterFileId, prefixName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) (string, error) {
	var parts = make([]string, 2)
	var err error
	s.errno = SplitFileId(masterFileId, parts)
	if s.errno != 0 {
		return "", nil
	}
	parts,err = s.UploadMasterCallback(parts[0], parts[1], prefixName, fileSize, callback, fileExtName, metaList)
	if err == nil && parts != nil {
		return parts[0] + SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR + parts[1], nil
	}

	return "", err
}

/**
 * append file to storage server (by file name)
 *
 * @param appender_file_id the appender file id
 * @param local_filename   local filename to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) AppendFile1(appenderFileId, localFilename string) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}
	return s.AppendFile(parts[0], parts[1], localFilename)
}

/**
 * append file to storage server (by file buff)
 *
 * @param appender_file_id the appender file id
 * @param file_buff        file content/buff
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) AppendBuffer1(appenderFileId string, fileBuffer []byte) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.AppendBuffer(parts[0], parts[1], fileBuffer)
}

/**
 * append file to storage server (by file buff)
 *
 * @param appender_file_id the appender file id
 * @param file_buff        file content/buffer
 * @param offset           start offset of the buffer
 * @param length           the length of the buffer to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) AppendOffsetBuffer1(appenderFileId string, fileBuffer []byte, offset, length int) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.AppendOffsetBuffer(parts[0], parts[1], fileBuffer, offset, length)
}

/**
 * append file to storage server (by callback)
 *
 * @param appender_file_id the appender file id
 * @param file_size        the file size
 * @param callback         the write data callback object
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) AppendCallback1(appenderFileId string, fileSize int, callback UploadCallback) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.AppendCallback(parts[0], parts[1], fileSize, callback)
}

/**
 * modify appender file to storage server (by file name)
 *
 * @param appender_file_id the appender file id
 * @param file_offset      the offset of appender file
 * @param local_filename   local filename to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) ModifyFile1(appenderFileId string, fileOffset int, localFilename string) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.ModifyFile(parts[0], parts[1], fileOffset, localFilename)
}

/**
 * modify appender file to storage server (by file buff)
 *
 * @param appender_file_id the appender file id
 * @param file_offset      the offset of appender file
 * @param file_buff        file content/buff
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) ModifyBuffer1(appenderFileId string, fileOffset int, fileBuffer []byte) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.ModifyBuffer(parts[0], parts[1], fileOffset, fileBuffer)
}

/**
 * modify appender file to storage server (by file buff)
 *
 * @param appender_file_id the appender file id
 * @param file_offset      the offset of appender file
 * @param file_buff        file content/buff
 * @param buffer_offset    start offset of the buff
 * @param buffer_length    the length of buff to modify
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) ModifyOffsetBuffer1(appenderFileId string, fileOffset int, fileBuff []byte, bufferOffset, bufferLength int) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.ModifyOffsetBuffer(parts[0], parts[1], fileOffset, fileBuff, bufferOffset, bufferLength)
}

/**
 * modify appender file to storage server (by callback)
 *
 * @param appender_file_id the appender file id
 * @param file_offset      the offset of appender file
 * @param modify_size      the modify size
 * @param callback         the write data callback object
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient1) ModifyCallback1(appenderFileId string, fileOffset, modifySize int, callback UploadCallback) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(appenderFileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.ModifyCallback(parts[0], parts[1], fileOffset, modifySize, callback)
}

/**
 * delete file from storage server
 *
 * @param file_id the file id(including group name and filename)
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient1) DeleteFile1(fileId string) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.DeleteFile(parts[0], parts[1])
}

/**
 * truncate appender file to size 0 from storage server
 *
 * @param appender_file_id the appender file id
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient1) TruncateFile1(fileId string) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.TruncateFile(parts[0], parts[1])
}

/**
 * truncate appender file from storage server
 *
 * @param appender_file_id    the appender file id
 * @param truncated_file_size truncated file size
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient1) TruncateFileBySize1(fileId string, truncatedFileSize int) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.TruncateFileBySize(parts[0], parts[1], truncatedFileSize)
}

/**
 * download file from storage server
 *
 * @param file_id the file id(including group name and filename)
 * @return file content/buffer, return null if fail
 */
func (s *StorageClient1) DownloadBuffer1(fileId string) ([]byte, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadOffsetBuffer1(fileId, fileOffset, downloadBytes)
}

/**
 * download file from storage server
 *
 * @param file_id        the file id(including group name and filename)
 * @param file_offset    the start offset of the file
 * @param download_bytes download bytes, 0 for remain bytes from offset
 * @return file content/buff, return null if fail
 */
func (s *StorageClient1) DownloadOffsetBuffer1(fileId string, fileOffset, downloadBytes int) ([]byte, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return nil, fmt.Errorf("errno:%d", s.errno)
	}

	return s.DownloadOffsetBuffer(parts[0], parts[1], fileOffset, downloadBytes)
}

/**
 * download file from storage server
 *
 * @param file_id        the file id(including group name and filename)
 * @param local_filename the filename on local
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient1) DownloadFile1(fileId, localFilename string) (int, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadFileByOffsetBuffer1(fileId, fileOffset, downloadBytes, localFilename)
}

/**
 * download file from storage server
 *
 * @param file_id        the file id(including group name and filename)
 * @param file_offset    the start offset of the file
 * @param download_bytes download bytes, 0 for remain bytes from offset
 * @param local_filename the filename on local
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient1) DownloadFileByOffsetBuffer1(fileId string, fileOffset, downloadBytes int, localFilename string) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.DownloadFileByOffsetBuffer(parts[0], parts[1], fileOffset, downloadBytes, localFilename)
}

/**
 * download file from storage server
 *
 * @param file_id  the file id(including group name and filename)
 * @param callback the callback object, will call callback.recv() when data arrive
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient1) DownloadCallback1(fileId string, callback DownloadCallback) (int, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadCallbackByOffsetBuffer1(fileId, fileOffset, downloadBytes, callback)
}

/**
 * download file from storage server
 *
 * @param file_id        the file id(including group name and filename)
 * @param file_offset    the start offset of the file
 * @param download_bytes download bytes, 0 for remain bytes from offset
 * @param callback       the callback object, will call callback.recv() when data arrive
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient1) DownloadCallbackByOffsetBuffer1(fileId string, fileOffset, downloadBytes int, callback DownloadCallback) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.DownloadCallbackByOffsetBuffer(parts[0], parts[1], fileOffset, downloadBytes, callback)
}

/**
 * get all metadata items from storage server
 *
 * @param file_id the file id(including group name and filename)
 * @return meta info array, return null if fail
 */
func (s *StorageClient1) GetMetadata1(fileId string) ([]NameValuePair, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return nil, fmt.Errorf("errno:%d", s.errno)
	}

	return s.GetMetadata(parts[0], parts[1])
}

/**
 * set metadata items to storage server
 *
 * @param file_id   the file id(including group name and filename)
 * @param meta_list meta item array
 * @param op_flag   flag, can be one of following values: <br>
 *                  <ul><li> ProtoCommon.STORAGE_SET_METADATA_FLAG_OVERWRITE: overwrite all old
 *                  metadata items</li></ul>
 *                  <ul><li> ProtoCommon.STORAGE_SET_METADATA_FLAG_MERGE: merge, insert when
 *                  the metadata item not exist, otherwise update it</li></ul>
 * @return 0 for success, !=0 fail (error code)
 */
func (s *StorageClient1) SetMetadata1(fileId string, metaList []NameValuePair, opFlag byte) (int, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return s.SetMetadata(parts[0], parts[1], metaList, opFlag)
}

/**
 * get file info from storage server
 *
 * @param file_id the file id(including group name and filename)
 * @return FileInfo object for success, return null for fail
 */
func (s *StorageClient1) QueryFileInfo1(fileId string) (*FileInfo, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return nil, fmt.Errorf("errno:%d", s.errno)
	}

	return s.QueryFileInfo(parts[0], parts[1])
}

/**
 * get file info decoded from filename
 *
 * @param file_id the file id(including group name and filename)
 * @return FileInfo object for success, return null for fail
 */
func (s *StorageClient1) GetFileInfo1(fileId string) (*FileInfo, error) {
	var parts = make([]string, 2)
	s.errno = SplitFileId(fileId, parts)
	if s.errno != 0 {
		return nil, fmt.Errorf("errno:%d", s.errno)
	}

	return s.GetFileInfo(parts[0], parts[1])
}