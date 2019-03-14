package fastdfs

import (
	"os"
	"strings"
	"errors"
	"io"
	"net"
	"fmt"
	"runtime/debug"
)

var base64 = NewBase64ByDetailed('-', '_', '.', 0)

type StorageClient struct {
	trackerServer   *TrackerServer
	storageServer   *StorageServer
	errno           byte
}

/**
 * constructor using global settings in class ClientGlobal
 */
func NewStorageClient() *StorageClient {
	return new(StorageClient)
}

/**
 * constructor with tracker server and storage server
 *
 * @param trackerServer the tracker server, can be null
 * @param storageServer the storage server, can be null
 */
func NewStorageClientByServer(trackerServer *TrackerServer, storageServer *StorageServer) *StorageClient {
	var storageClient = new(StorageClient)
	storageClient.trackerServer = trackerServer
	storageClient.storageServer = storageServer

	return storageClient
}

/**
 * get the error code of last call
 *
 * @return the error code of last call
 */
func (s *StorageClient) GetErrorCode() byte {
	return s.errno
}

/**
 * upload file to storage server (by file name)
 *
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadFile(localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.uploadFileByGroup(groupName, localFilename, fileExtName, metaList)
}

/**
 * upload file to storage server (by file name)
 *
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) uploadFileByGroup(groupName, localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const cmd = STORAGE_PROTO_CMD_UPLOAD_FILE

	return s.uploadFileByCmd(cmd, groupName, localFilename, fileExtName, metaList)
}

/**
 * upload file to storage server (by file name)
 *
 * @param cmd            the command
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) uploadFileByCmd(cmd byte, groupName, localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	file,err := os.Open(localFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if fileExtName == "" {
		var nPos = strings.LastIndexByte(localFilename, '.')
		if nPos > 0 && len(localFilename) - nPos <= FDFS_FILE_EXT_NAME_MAX_LEN + 1 {
			fileExtName = localFilename[nPos + 1:]
		}
	}

	stat,err := file.Stat()
	if err != nil {
		return nil, err
	}

	return s.doUploadFile(cmd, groupName, "", "", fileExtName, int(stat.Size()), NewUploadStream(file, int(stat.Size())), metaList)
}

/**
 * upload file to storage server (by file buff)
 *
 * @param file_buff     file content/buff
 * @param offset        start offset of the buff
 * @param length        the length of buff to upload
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadBufferOffset(fileBuff []byte, offset int, length int, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.UploadBufferOffsetByGroup(groupName, fileBuff, offset, length, fileExtName, metaList)
}

/**
 * upload file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param offset        start offset of the buff
 * @param length        the length of buff to upload
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadBufferOffsetByGroup(groupName string, fileBuff []byte, offset int, length int, fileExtName string, metaList []NameValuePair) ([]string, error) {
	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_FILE, groupName, "", "", fileExtName, length, NewUploadBuff(fileBuff, offset, length), metaList)
}

/**
 * upload file to storage server (by file buff)
 *
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadBuffer(fileBuff []byte, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.UploadBufferByGroup(groupName, fileBuff, fileExtName, metaList)
}

/**
 * upload file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadBufferByGroup(groupName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) ([]string, error) {
	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_FILE, groupName, "", "", fileExtName, len(fileBuff), NewUploadBuff(fileBuff, 0, len(fileBuff)), metaList)
}

/**
 * upload file to storage server (by callback)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_size     the file size
 * @param callback      the write data callback object
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadCallback(groupName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const masterFilename = ""
	const prefixName = ""

	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_FILE, groupName, masterFilename, prefixName, fileExtName, fileSize, callback, metaList)
}

/**
 * upload file to storage server (by file name, slave file mode)
 *
 * @param group_name      the group name of master file
 * @param master_filename the master file name to generate the slave file
 * @param prefix_name     the prefix name to generate the slave file
 * @param local_filename  local filename to upload
 * @param file_ext_name   file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list       meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadMasterFile(groupName, masterFilename, prefixName, localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	if groupName == "" || len(groupName) == 0 || masterFilename == "" || len(masterFilename) == 0 || prefixName == "" {
		return nil, errors.New("invalid argument")
	}

	file,err := os.Open(localFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if fileExtName == "" {
		var nPos = strings.LastIndexByte(localFilename, '.')
		if nPos > 0 && len(localFilename) - nPos <= FDFS_FILE_EXT_NAME_MAX_LEN + 1 {
			fileExtName = localFilename[nPos + 1:]
		}
	}

	stat,err := file.Stat()
	if err != nil {
		return nil, err
	}

	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, groupName, masterFilename, prefixName, fileExtName, int(stat.Size()), NewUploadStream(file, int(stat.Size())), metaList)
}

/**
 * upload file to storage server (by file buff, slave file mode)
 *
 * @param group_name      the group name of master file
 * @param master_filename the master file name to generate the slave file
 * @param prefix_name     the prefix name to generate the slave file
 * @param file_buff       file content/buff
 * @param file_ext_name   file ext name, do not include dot(.)
 * @param meta_list       meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadMasterBuffer(groupName, masterFilename, prefixName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) ([]string, error) {
	if groupName == "" || len(groupName) == 0 || masterFilename == "" || len(masterFilename) == 0 || prefixName == "" {
		return nil, errors.New("invalid argument")
	}

	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, groupName, masterFilename, prefixName, fileExtName, len(fileBuff), NewUploadBuff(fileBuff, 0, len(fileBuff)), metaList)
}

/**
 * upload file to storage server (by file buff, slave file mode)
 *
 * @param group_name      the group name of master file
 * @param master_filename the master file name to generate the slave file
 * @param prefix_name     the prefix name to generate the slave file
 * @param file_buff       file content/buff
 * @param offset          start offset of the buff
 * @param length          the length of buff to upload
 * @param file_ext_name   file ext name, do not include dot(.)
 * @param meta_list       meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadMasterOffsetBuffer(groupName, masterFilename, prefixName string, fileBuff []byte, offset, length int, fileExtName string, metaList []NameValuePair) ([]string, error) {
	if groupName == "" || len(groupName) == 0 || masterFilename == "" || len(masterFilename) == 0 || prefixName == "" {
		return nil, errors.New("invalid argument")
	}

	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, groupName, masterFilename, prefixName, fileExtName, length, NewUploadBuff(fileBuff, offset, length), metaList)
}

/**
 * upload file to storage server (by callback, slave file mode)
 *
 * @param group_name      the group name to upload file to, can be empty
 * @param master_filename the master file name to generate the slave file
 * @param prefix_name     the prefix name to generate the slave file
 * @param file_size       the file size
 * @param callback        the write data callback object
 * @param file_ext_name   file ext name, do not include dot(.)
 * @param meta_list       meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadMasterCallback(groupName, masterFilename, prefixName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) ([]string, error) {
	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, groupName, masterFilename, prefixName, fileExtName, fileSize, callback, metaList)
}

/**
 * upload appender file to storage server (by file name)
 *
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderFile(localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.UploadAppenderFileByGroup(groupName, localFilename, fileExtName, metaList)
}

/**
 * upload appender file to storage server (by file name)
 *
 * @param group_name     the group name to upload file to, can be empty
 * @param local_filename local filename to upload
 * @param file_ext_name  file ext name, do not include dot(.), null to extract ext name from the local filename
 * @param meta_list      meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file </li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderFileByGroup(groupName, localFilename, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const cmd = STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE

	return s.uploadFileByCmd(cmd, groupName, localFilename, fileExtName, metaList)
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param file_buff     file content/buff
 * @param offset        start offset of the buff
 * @param length        the length of buff to upload
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderOffsetBuffer(fileBuff []byte, offset, length int, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.UploadAppenderOffsetBufferByGroup(groupName, fileBuff, offset, length, fileExtName, metaList)
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param offset        start offset of the buff
 * @param length        the length of buff to upload
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderOffsetBufferByGroup(groupName string, fileBuff []byte, offset, length int, fileExtName string, metaList []NameValuePair) ([]string, error) {
	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE, groupName, "", "", fileExtName, length, NewUploadBuff(fileBuff, offset, length), metaList)
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderBuffer(fileBuff []byte, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const groupName = ""

	return s.UploadAppenderOffsetBufferByGroup(groupName, fileBuff, 0, len(fileBuff)	, fileExtName, metaList)
}

/**
 * upload appender file to storage server (by file buff)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_buff     file content/buff
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderBufferByGroup(groupName string, fileBuff []byte, fileExtName string, metaList []NameValuePair) ([]string, error) {
	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE, groupName, "", "", fileExtName, len(fileBuff), NewUploadBuff(fileBuff, 0, len(fileBuff)), metaList)
}

/**
 * upload appender file to storage server (by callback)
 *
 * @param group_name    the group name to upload file to, can be empty
 * @param file_size     the file size
 * @param callback      the write data callback object
 * @param file_ext_name file ext name, do not include dot(.)
 * @param meta_list     meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li>results[0]: the group name to store the file</li></ul>
 * <ul><li>results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) UploadAppenderCallback(groupName string, fileSize int, callback UploadCallback, fileExtName string, metaList []NameValuePair) ([]string, error) {
	const masterFilename = ""
	const prefixName = ""

	return s.doUploadFile(STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE, groupName, masterFilename, prefixName, fileExtName, fileSize, callback, metaList)
}

/**
 * append file to storage server (by file name)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param local_filename    local filename to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) AppendFile(groupName, appenderFilename, localFilename string) (int, error) {
	file,err := os.Open(localFilename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat,err := file.Stat()
	if err != nil {
		return 0, err
	}

	return s.doAppendFile(groupName, appenderFilename, int(stat.Size()), NewUploadStream(file, int(stat.Size())))
}

/**
 * append file to storage server (by file buff)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_buff         file content/buff
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) AppendBuffer(groupName, appenderFilename string, fileBuff []byte) (int, error) {
	return s.doAppendFile(groupName, appenderFilename, len(fileBuff), NewUploadBuff(fileBuff, 0, len(fileBuff)))
}

/**
 * append file to storage server (by file buff)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_buff         file content/buff
 * @param offset            start offset of the buff
 * @param length            the length of buff to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) AppendOffsetBuffer(groupName, appenderFilename string, fileBuff []byte, offset, length int) (int, error) {
	return s.doAppendFile(groupName, appenderFilename, length, NewUploadBuff(fileBuff, offset, length))
}

/**
 * append file to storage server (by callback)
 *
 * @param group_name        the group name to append file to
 * @param appender_filename the appender filename
 * @param file_size         the file size
 * @param callback          the write data callback object
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) AppendCallback(groupName, appenderFilename string, fileSize int, callback UploadCallback) (int, error) {
	return s.doAppendFile(groupName, appenderFilename, fileSize, callback)
}

/**
 * modify appender file to storage server (by file name)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_offset       the offset of appender file
 * @param local_filename    local filename to append
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) ModifyFile(groupName, appenderFilename string, fileOffset int, localFilename string) (int, error) {
	file,err := os.Open(localFilename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat,err := file.Stat()
	if err != nil {
		return 0, err
	}

	return s.doModifyFile(groupName, appenderFilename, fileOffset, int(stat.Size()), NewUploadStream(file, int(stat.Size())))
}

/**
 * modify appender file to storage server (by file buff)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_offset       the offset of appender file
 * @param file_buff         file content/buff
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) ModifyBuffer(groupName, appenderFilename string, fileOffset int, fileBuff []byte) (int, error) {
	return s.doModifyFile(groupName, appenderFilename, fileOffset, len(fileBuff), NewUploadBuff(fileBuff, 0, len(fileBuff)))
}

/**
 * modify appender file to storage server (by file buff)
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_offset       the offset of appender file
 * @param file_buff         file content/buff
 * @param buffer_offset     start offset of the buff
 * @param buffer_length     the length of buff to modify
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) ModifyOffsetBuffer(groupName, appenderFilename string, fileOffset int, fileBuff []byte, bufferOffset, bufferLength int) (int, error) {
	return s.doModifyFile(groupName, appenderFilename, fileOffset, bufferLength, NewUploadBuff(fileBuff, bufferOffset, bufferLength))
}

/**
 * modify appender file to storage server (by callback)
 *
 * @param group_name        the group name to modify file to
 * @param appender_filename the appender filename
 * @param file_offset       the offset of appender file
 * @param modify_size       the modify size
 * @param callback          the write data callback object
 * @return 0 for success, != 0 for error (error no)
 */
func (s *StorageClient) ModifyCallback(groupName, appenderFilename string, fileOffset, modifySize int, callback UploadCallback) (int, error) {
	return s.doModifyFile(groupName, appenderFilename, fileOffset, modifySize, callback)
}

/**
 * upload file to storage server
 *
 * @param cmd             the command code
 * @param group_name      the group name to upload file to, can be empty
 * @param master_filename the master file name to generate the slave file
 * @param prefix_name     the prefix name to generate the slave file
 * @param file_ext_name   file ext name, do not include dot(.)
 * @param file_size       the file size
 * @param callback        the write data callback object
 * @param meta_list       meta info array
 * @return 2 elements string array if success:<br>
 * <ul><li> results[0]: the group name to store the file</li></ul>
 * <ul><li> results[1]: the new created filename</li></ul>
 * return null if fail
 */
func (s *StorageClient) doUploadFile(cmd byte, groupName, masterFilename, prefixName, fileExtName string, fileSize int, callback UploadCallback, metaList []NameValuePair) ([]string, error) {
	var (
		header []byte
		extNameBs []byte
		newGroupName string
		remoteFilename string
		bNewConnection bool
		storageSocket  net.Conn
		sizeBytes []byte
		hexLenBytes []byte
		masterFilenameBytes []byte
		bUploadSlave bool
		offset int
		bodyLen int
	)
	var err error

	bUploadSlave = (groupName != "" && len(groupName) > 0) && (masterFilename != "" && len(masterFilename) > 0) && (prefixName != "")
	if bUploadSlave {
		bNewConnection,err = s.newUpdatableStorageConnection(groupName, masterFilename)
	} else {
		bNewConnection,err = s.newWritableStorageConnection(groupName)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()

	if storageSocket,err = s.storageServer.GetSocket(); err != nil {
		return nil, err
	}

	extNameBs = make([]byte, FDFS_FILE_EXT_NAME_MAX_LEN)
	if fileExtName != "" && len(fileExtName) > 0 {
		var bs,err = ConvertBytesToUTF8([]byte(fileExtName), GCharset)
		if err != nil {
			return nil, err
		}
		var extNameLen = len(bs)
		if extNameLen > FDFS_FILE_EXT_NAME_MAX_LEN {
			extNameLen = FDFS_FILE_EXT_NAME_MAX_LEN
		}
		copy(extNameBs, bs[:extNameLen])
	}

	if bUploadSlave {
		if masterFilenameBytes,err = ConvertBytesToUTF8([]byte(masterFilename), GCharset); err != nil {
			return nil, err
		}
		sizeBytes = make([]byte, 2 * FDFS_PROTO_PKG_LEN_SIZE)
		bodyLen = len(sizeBytes) + FDFS_FILE_PREFIX_MAX_LEN + FDFS_FILE_EXT_NAME_MAX_LEN + len(masterFilenameBytes) + fileSize

		hexLenBytes = Long2Buff(int64(len(masterFilename)))
		copy(sizeBytes, hexLenBytes)
		offset = len(hexLenBytes)
	} else {
		masterFilename = ""
		sizeBytes = make([]byte, 1 + 1 * FDFS_PROTO_PKG_LEN_SIZE)
		bodyLen = len(sizeBytes) + FDFS_FILE_EXT_NAME_MAX_LEN + fileSize

		sizeBytes[0] = byte(s.storageServer.GetStorePathIndex())
		offset = 1
	}

	hexLenBytes = Long2Buff(int64(fileSize))
	copy(sizeBytes[offset:], hexLenBytes)

	if header,err = PackHeader(cmd, int64(bodyLen), 0); err != nil {
		return nil, err
	}
	var wholePkg = make([]byte, len(header) + bodyLen - fileSize)
	copy(wholePkg, header)
	copy(wholePkg[len(header):], sizeBytes)
	offset = len(header) + len(sizeBytes)
	if bUploadSlave {
		var prefixNameBs = make([]byte, FDFS_FILE_PREFIX_MAX_LEN)
		var bs,err = ConvertBytesToUTF8([]byte(prefixName), GCharset)
		if err != nil {
			return nil, err
		}
		var prefixNameLen = len(bs)
		if prefixNameLen > FDFS_FILE_PREFIX_MAX_LEN {
			prefixNameLen = FDFS_FILE_PREFIX_MAX_LEN
		}
		if prefixNameLen > 0 {
			copy(prefixNameBs, bs)
		}

		copy(wholePkg[offset:], prefixNameBs)
		offset += len(prefixNameBs)
	}

	copy(wholePkg[offset:], extNameBs)
	offset += len(extNameBs)

	if bUploadSlave {
		copy(wholePkg[offset:], masterFilenameBytes)
		offset += len(masterFilenameBytes)
	}

	if _,err = storageSocket.Write(wholePkg); err != nil {
		return nil, err
	}

	errno,_ := callback.Send(storageSocket)
	s.errno = byte(errno)
	if s.errno != 0 {
		return nil, nil
	}

	pkgInfo,err := RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, -1)
	if err != nil {
		return nil, err
	}
	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		return nil, nil
	}

	if len(pkgInfo.Body) <= FDFS_GROUP_NAME_MAX_LEN {
		return nil, fmt.Errorf("body length: %d <= %d", len(pkgInfo.Body), FDFS_GROUP_NAME_MAX_LEN)
	}

	newGroupName = strings.Trim(string(pkgInfo.Body[:FDFS_GROUP_NAME_MAX_LEN]), " \x00")
	remoteFilename = string(pkgInfo.Body[FDFS_GROUP_NAME_MAX_LEN: FDFS_GROUP_NAME_MAX_LEN + len(pkgInfo.Body) - FDFS_GROUP_NAME_MAX_LEN])
	var results = make([]string, 2)
	results[0] = newGroupName
	results[1] = remoteFilename

	if metaList == nil || len(metaList) == 0 {
		return results, nil
	}

	var result = 0
	result,err = s.SetMetadata(newGroupName, remoteFilename, metaList, STORAGE_SET_METADATA_FLAG_OVERWRITE)
	if err != nil || result != 0 {
		if err != nil {
			result = 5
		}
		s.errno = byte(result)
		s.DeleteFile(newGroupName, remoteFilename)
		// todo return nil or error?
		// java is return null
		return nil, err
	}

	return results, nil
}

/**
 * append file to storage server
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_size         the file size
 * @param callback          the write data callback object
 * @return return true for success, false for fail
 */
func (s *StorageClient) doAppendFile(groupName, appenderFilename string, fileSize int, callback UploadCallback) (int, error) {
	var (
		header []byte
		bNewConnection bool
		storageSocket net.Conn
		hexLenBytes []byte
		appenderFilenameBytes []byte
		offset int
		bodyLen int
	)
	var err error

	if groupName == "" || len(groupName) == 0 || appenderFilename == "" || len(appenderFilename) == 0 {
		s.errno = ERR_NO_EINVAL
		return int(s.errno), errors.New("ERR_NO_EINVAL")
	}

	if bNewConnection,err = s.newUpdatableStorageConnection(groupName, appenderFilename); err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()

	if storageSocket,err = s.storageServer.GetSocket(); err != nil {
		return -1, err
	}

	if appenderFilenameBytes,err = ConvertBytesToUTF8([]byte(appenderFilename), GCharset); err != nil {
		return -1, err
	}
	bodyLen = 2 * FDFS_PROTO_PKG_LEN_SIZE + len(appenderFilenameBytes) + fileSize

	if header,err = PackHeader(STORAGE_PROTO_CMD_APPEND_FILE, int64(bodyLen), 0); err != nil {
		return -1, err
	}
	var wholePkg = make([]byte, len(header) + bodyLen - fileSize)
	copy(wholePkg, header)
	offset = len(header)

	hexLenBytes = Long2Buff(int64(len(appenderFilename)))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	hexLenBytes = Long2Buff(int64(fileSize))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	copy(wholePkg[offset:], appenderFilenameBytes)
	offset += len(appenderFilenameBytes)

	if _,err = storageSocket.Write(wholePkg); err != nil {
		return -1, err
	}
	if n,err := callback.Send(storageSocket); err != nil {
		s.errno = byte(n)
		return n, err
	}

	pkgInfo,err := RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 0)
	if err != nil {
		return -1, err
	}
	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return 0, nil
}

/**
 * modify appender file to storage server
 *
 * @param group_name        the group name of appender file
 * @param appender_filename the appender filename
 * @param file_offset       the offset of appender file
 * @param modify_size       the modify size
 * @param callback          the write data callback object
 * @return return true for success, false for fail
 */
func (s *StorageClient) doModifyFile(groupName, appenderFilename string, fileOffset, modifySize int, callback UploadCallback) (int, error) {
	var (
		header []byte
		bNewConnection bool
		storageSocket net.Conn
		hexLenBytes []byte
		appenderFilenameBytes []byte
		offset int
		bodyLen int
	)
	var err error

	if groupName == "" || len(groupName) == 0 || appenderFilename == "" || len(appenderFilename) == 0 {
		s.errno = ERR_NO_EINVAL
		return int(s.errno), errors.New("ERR_NO_EINVAL")
	}

	if bNewConnection,err = s.newUpdatableStorageConnection(groupName, appenderFilename); err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()

	if storageSocket,err = s.storageServer.GetSocket(); err != nil {
		return -1, err
	}

	if appenderFilenameBytes,err = ConvertBytesToUTF8([]byte(appenderFilename), GCharset); err != nil {
		return -1, err
	}
	bodyLen = 3 * FDFS_PROTO_PKG_LEN_SIZE + len(appenderFilenameBytes) + modifySize

	if header,err = PackHeader(STORAGE_PROTO_CMD_MODIFY_FILE, int64(bodyLen), 0); err != nil {
		return -1, err
	}
	var wholePkg = make([]byte, len(header) + bodyLen - modifySize)
	copy(wholePkg, header)
	offset = len(header)

	hexLenBytes = Long2Buff(int64(len(appenderFilename)))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	hexLenBytes = Long2Buff(int64(fileOffset))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	hexLenBytes = Long2Buff(int64(modifySize))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	copy(wholePkg[offset:], appenderFilenameBytes)
	offset += len(appenderFilenameBytes)

	if _,err = storageSocket.Write(wholePkg); err != nil {
		return -1, err
	}
	if n,err := callback.Send(storageSocket); err != nil {
		s.errno = byte(n)
		return n, err
	}

	pkgInfo,err := RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 0)
	if err != nil {
		return -1, err
	}
	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return 0, nil
}

/**
 * delete file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient) DeleteFile(groupName, remoteFilename string) (int, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return -1, err
	}

	if err = s.sendPackage(STORAGE_PROTO_CMD_DELETE_FILE, groupName, remoteFilename); err != nil {
		return -1, err
	}
	pkgInfo,err := RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 0)
	if err != nil {
		return -1, err
	}

	s.errno = pkgInfo.Errno
	return int(pkgInfo.Errno), nil
}

/**
 * truncate appender file to size 0 from storage server
 *
 * @param group_name        the group name of storage server
 * @param appender_filename the appender filename
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient) TruncateFile(groupName, appenderFilename string) (int, error) {
	const truncatedFileSize = 0

	return s.TruncateFileBySize(groupName, appenderFilename, truncatedFileSize)
}

/**
 * truncate appender file from storage server
 *
 * @param group_name          the group name of storage server
 * @param appender_filename   the appender filename
 * @param truncated_file_size truncated file size
 * @return 0 for success, none zero for fail (error code)
 */
func (s *StorageClient) TruncateFileBySize(groupName, appenderFilename string, truncatedFileSize int) (int, error) {
	var (
		header []byte
		bNewConnection bool
		storageSocket net.Conn
		hexLenBytes []byte
		appenderFilenameBytes []byte
		offset int
		bodyLen int
	)
	var err error

	if groupName == "" || len(groupName) == 0 || appenderFilename == "" || len(appenderFilename) == 0 {
		s.errno = ERR_NO_EINVAL
		return int(s.errno), errors.New("ERR_NO_EINVAL")
	}

	if bNewConnection,err = s.newUpdatableStorageConnection(groupName, appenderFilename); err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()

	if storageSocket,err = s.storageServer.GetSocket(); err != nil {
		return -1, err
	}

	if appenderFilenameBytes,err = ConvertBytesToUTF8([]byte(appenderFilename), GCharset); err != nil {
		return -1, err
	}
	bodyLen = 2 * FDFS_PROTO_PKG_LEN_SIZE + len(appenderFilenameBytes)

	if header,err = PackHeader(STORAGE_PROTO_CMD_TRUNCATE_FILE, int64(bodyLen), 0); err != nil {
		return -1, err
	}
	var wholePkg = make([]byte, len(header) + bodyLen)
	copy(wholePkg, header)
	offset = len(header)

	hexLenBytes = Long2Buff(int64(len(appenderFilename)))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	hexLenBytes = Long2Buff(int64(truncatedFileSize))
	copy(wholePkg[offset:], hexLenBytes)
	offset += len(hexLenBytes)

	copy(wholePkg[offset:], appenderFilenameBytes)
	offset += len(appenderFilenameBytes)

	if _,err = storageSocket.Write(wholePkg); err != nil {
		return -1, err
	}

	pkgInfo,err := RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 0)
	if err != nil {
		return -1, err
	}
	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return 0, nil
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return file content/buff, return null if fail
 */
func (s *StorageClient) DownloadBuffer(groupName, remoteFilename string) ([]byte, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadOffsetBuffer(groupName, remoteFilename, fileOffset, downloadBytes)
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param file_offset     the start offset of the file
 * @param download_bytes  download bytes, 0 for remain bytes from offset
 * @return file content/buff, return null if fail
 */
func (s *StorageClient) DownloadOffsetBuffer(groupName, remoteFilename string, fileOffset, downloadBytes int) ([]byte, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return nil, err
	}

	var pkgInfo *RecvPackageInfo

	if err = s.sendDownloadPackage(groupName, remoteFilename, fileOffset, downloadBytes); err != nil {
		return nil, err
	}
	if pkgInfo,err = RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, -1); err != nil {
		return nil, err
	}

	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		return nil, fmt.Errorf("errno:%d", s.errno)
	}

	return pkgInfo.Body, nil
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param local_filename  filename on local
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient) DownloadFile(groupName, remoteFilename, localFilename string) (int, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadFileByOffsetBuffer(groupName, remoteFilename, fileOffset, downloadBytes, localFilename)
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param file_offset     the start offset of the file
 * @param download_bytes  download bytes, 0 for remain bytes from offset
 * @param local_filename  filename on local
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient) DownloadFileByOffsetBuffer(groupName, remoteFilename string, fileOffset, downloadBytes int, localFilename string) (int, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return -1, err
	}
	var header *RecvHeaderInfo
	file,err := os.Create(localFilename)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	s.errno = 0
	defer func() {
		if err != nil {
			s.errno = ERR_NO_EIO
		}
		if s.errno != 0 {
			os.Remove(localFilename)
		}
	}()
	if err = s.sendDownloadPackage(groupName, remoteFilename, fileOffset, downloadBytes); err != nil {
		return -1, err
	}

	if header,err = RecvHeader(storageSocket, STORAGE_PROTO_CMD_RESP, -1); err != nil {
		return -1, err
	}
	s.errno = header.Errno
	if header.Errno != 0 {
		return int(header.Errno), fmt.Errorf("errno:%d", header.Errno)
	}

	var buff = make([]byte, 256 * 1024)
	var remainBytes = header.BodyLen
	var bytes int

	for remainBytes > 0 {
		var length = remainBytes
		if length > len(buff) {
			length = len(buff)
		}
		if bytes,err = storageSocket.Read(buff[:length]); err != nil {
			return -1, err
		}
		if bytes < 0 {
			return -1, fmt.Errorf("recv package size %d != %d", header.BodyLen - remainBytes, header.BodyLen)
		}

		if _,err = file.Write(buff[:bytes]); err != nil {
			return -1, err
		}

		remainBytes -= bytes
	}

	return 0, nil
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param callback        call callback.recv() when data arrive
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient) DownloadCallback(groupName, remoteFilename string, callback DownloadCallback) (int, error) {
	const fileOffset = 0
	const downloadBytes = 0

	return s.DownloadCallbackByOffsetBuffer(groupName, remoteFilename, fileOffset, downloadBytes, callback)
}

/**
 * download file from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param file_offset     the start offset of the file
 * @param download_bytes  download bytes, 0 for remain bytes from offset
 * @param callback        call callback.recv() when data arrive
 * @return 0 success, return none zero errno if fail
 */
func (s *StorageClient) DownloadCallbackByOffsetBuffer(groupName, remoteFilename string, fileOffset, downloadBytes int, callback DownloadCallback) (int, error) {
	var result int
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return -1, err
	}

	var header *RecvHeaderInfo
	if err = s.sendDownloadPackage(groupName, remoteFilename,fileOffset, downloadBytes); err != nil {
		return -1, err
	}
	if header,err = RecvHeader(storageSocket, STORAGE_PROTO_CMD_RESP, -1); err != nil {
		return -1, err
	}
	s.errno = header.Errno
	if header.Errno != 0 {
		return int(header.Errno), fmt.Errorf("errno:%d", header.Errno)
	}

	var buff = make([]byte, 2 * 1024)
	var remainBytes = header.BodyLen
	var bytes int

	for remainBytes > 0 {
		var length = remainBytes
		if length > len(buff) {
			length = len(buff)
		}
		if bytes,err = storageSocket.Read(buff[:length]); err != nil {
			return -1, err
		}
		if bytes < 0 {
			return -1, fmt.Errorf("recv package size %d != %d", header.BodyLen - remainBytes, header.BodyLen)
		}
		if result,err = callback.Recv(header.BodyLen, buff, bytes); err != nil {
			s.errno = byte(result)
			return result, err
		}

		remainBytes -= bytes
	}

	return 0, nil
}

/**
 * get all metadata items from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return meta info array, return null if fail
 */
func (s *StorageClient) GetMetadata(groupName, remoteFilename string) ([]NameValuePair, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return nil, err
	}

	var pkgInfo *RecvPackageInfo

	if err = s.sendPackage(STORAGE_PROTO_CMD_GET_METADATA, groupName, remoteFilename); err != nil {
		return nil, err
	}
	if pkgInfo,err = RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, -1); err != nil {
		return nil, err
	}

	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return error or nil?
		// java is return null.
		return nil, nil
	}

	str,err := ConvertByteToString(pkgInfo.Body, GCharset)
	if err != nil {
		return nil, err
	}
	return SplitMetadata(str), nil
}

/**
 * set metadata items to storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param meta_list       meta item array
 * @param op_flag         flag, can be one of following values: <br>
 *                        <ul><li> ProtoCommon.STORAGE_SET_METADATA_FLAG_OVERWRITE: overwrite all old
 *                        metadata items</li></ul>
 *                        <ul><li> ProtoCommon.STORAGE_SET_METADATA_FLAG_MERGE: merge, insert when
 *                        the metadata item not exist, otherwise update it</li></ul>
 * @return 0 for success, !=0 fail (error code)
 */
func (s *StorageClient) SetMetadata(groupName, remoteFilename string, metaList []NameValuePair, opFlag byte) (int, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return -1, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return -1, err
	}

	var (
		header []byte
		groupBytes []byte
		filenameBytes []byte
		metaBuff []byte
		bs []byte
		groupLen int
		sizeBytes []byte
		pkgInfo *RecvPackageInfo
	)

	if metaList == nil {
		metaBuff = make([]byte, 0)
	} else {
		if metaBuff,err = ConvertBytesToUTF8([]byte(PackMetadata(metaList)), GCharset); err != nil {
			return -1, err
		}
	}

	if filenameBytes,err = ConvertBytesToUTF8([]byte(remoteFilename), GCharset); err != nil {
		return -1, err
	}
	sizeBytes = make([]byte, 2 * FDFS_PROTO_PKG_LEN_SIZE)

	bs = Long2Buff(int64(len(filenameBytes)))
	copy(sizeBytes, bs)
	bs = Long2Buff(int64(len(metaBuff)))
	copy(sizeBytes[FDFS_PROTO_PKG_LEN_SIZE:], bs)

	groupBytes = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return -1, err
	}

	if len(bs) <= len(groupBytes) {
		groupLen = len(bs)
	} else {
		groupLen = len(groupBytes)
	}

	copy(groupBytes[:groupLen], bs)
	if header,err = PackHeader(STORAGE_PROTO_CMD_SET_METADATA, int64(2 * FDFS_PROTO_PKG_LEN_SIZE + 1 + len(groupBytes) + len(filenameBytes) + len(metaBuff)), 0); err != nil {
		return -1, err
	}
	var wholePkg = make([]byte, len(header) + len(sizeBytes) + 1 + len(groupBytes) + len(filenameBytes))
	copy(wholePkg, header)
	copy(wholePkg[len(header):], sizeBytes)
	wholePkg[len(header) + len(sizeBytes)] = opFlag
	copy(wholePkg[len(header) + len(sizeBytes) + 1:], groupBytes)
	copy(wholePkg[len(header) + len(sizeBytes) + 1 + len(groupBytes):], filenameBytes)
	if _,err = storageSocket.Write(wholePkg); err != nil {
		return -1, err
	}

	if len(metaBuff) > 0 {
		if _,err = storageSocket.Write(metaBuff); err != nil {
			return -1, err
		}
	}

	if pkgInfo,err = RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 0); err != nil {
		return -1, err
	}

	s.errno = pkgInfo.Errno
	if s.errno != 0 {
		return int(s.errno), fmt.Errorf("errno:%d", s.errno)
	}

	return 0, nil
}

/**
 * get file info decoded from the filename, fetch from the storage if necessary
 *
 * @param group_name      the group name
 * @param remote_filename the filename
 * @return FileInfo object for success, return null for fail
 */
func (s *StorageClient) GetFileInfo(groupName, remoteFilename string) (*FileInfo, error) {
	if len(remoteFilename) < FDFS_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH + FDFS_FILE_EXT_NAME_MAX_LEN + 1 {
		s.errno = ERR_NO_EINVAL
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var buff,err = base64.DecodeAuto(remoteFilename[FDFS_FILE_PATH_LEN:FDFS_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH])
	if err != nil {
		return nil, err
	}

	var fileSize = Buff2long(buff, 4 * 2)
	if ((len(remoteFilename) > TRUNK_LOGIC_FILENAME_LENGTH) || ((len(remoteFilename) > NORMAL_LOGIC_FILENAME_LENGTH) && ((fileSize & TRUNK_FILE_MARK_SIZE) == 0))) || ((fileSize & APPENDER_FILE_SIZE) != 0) {
		//slave file or appender file
		return s.QueryFileInfo(groupName, remoteFilename)
	}

	var fileInfo = NewFileInfo(fileSize, 0, 0, GetIpAddress(buff, 0))
	fileInfo.SetCreateTimestamp(int64(Buff2int32(buff, 4)))
	if fileSize >> 63 != 0 {
		fileSize &= 0xFFFFFFFF  //low 32 bits is file size
		fileInfo.SetFileSize(fileSize)
	}
	fileInfo.SetCrc32(int(Buff2int32(buff, 4 * 4)))

	return fileInfo, nil
}

/**
 * get file info from storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return FileInfo object for success, return null for fail
 */
func (s *StorageClient) QueryFileInfo(groupName, remoteFilename string) (*FileInfo, error) {
	var bNewConnection,err = s.newUpdatableStorageConnection(groupName, remoteFilename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if bNewConnection {
			if err = s.storageServer.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				debug.PrintStack()
			}
			s.storageServer = nil
		}
	}()
	storageSocket,err := s.storageServer.GetSocket()
	if err != nil {
		return nil, err
	}

	var (
		header []byte
		groupBytes []byte
		filenameBytes []byte
		bs []byte
		groupLen int
		pkgInfo *RecvPackageInfo
	)

	if filenameBytes,err = ConvertBytesToUTF8([]byte(remoteFilename), GCharset); err != nil {
		return nil, err
	}
	groupBytes = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return nil, err
	}

	if len(bs) <= len(groupBytes) {
		groupLen = len(bs)
	} else {
		groupLen = len(groupBytes)
	}
	copy(groupBytes[:groupLen], bs)

	if header,err = PackHeader(STORAGE_PROTO_CMD_QUERY_FILE_INFO, int64(len(groupBytes) + len(filenameBytes)), 0); err != nil {
		return nil, err
	}

	var wholePkg = make([]byte, len(header) + len(groupBytes) + len(filenameBytes))
	copy(wholePkg, header)
	copy(wholePkg[len(header):], groupBytes)
	copy(wholePkg[len(header) + len(groupBytes):], filenameBytes)

	if _,err = storageSocket.Write(wholePkg); err != nil {
		return nil, err
	}

	if pkgInfo,err = RecvPackage(storageSocket, STORAGE_PROTO_CMD_RESP, 3 * FDFS_PROTO_PKG_LEN_SIZE +	 FDFS_IPADDR_SIZE); err != nil {
		return nil, err
	}

	s.errno = pkgInfo.Errno
	if pkgInfo.Errno != 0 {
		// todo return nil or error?
		// java is return null.
		return nil, nil
	}

	var fileSize = Buff2long(pkgInfo.Body, 0)
	var createTimestamp = Buff2long(pkgInfo.Body, FDFS_PROTO_PKG_LEN_SIZE)
	var crc32 = int(Buff2long(pkgInfo.Body, 2 * FDFS_PROTO_PKG_LEN_SIZE))
	var sourceIpAddr = strings.Trim(string(pkgInfo.Body[3 * FDFS_PROTO_PKG_LEN_SIZE:3 * FDFS_PROTO_PKG_LEN_SIZE + FDFS_IPADDR_SIZE]), " \x00")

	return NewFileInfo(fileSize, createTimestamp, crc32, sourceIpAddr), nil
}

/**
 * check storage socket, if null create a new connection
 *
 * @param group_name the group name to upload file to, can be empty
 * @return true if create a new connection
 */
func (s *StorageClient) newWritableStorageConnection(groupName string) (bool, error) {
	if s.storageServer != nil {
		return false, nil
	}
	var err error
	var tracker = NewTrackerClient()
	if s.storageServer,err = tracker.GetStoreStorageByGroup(s.trackerServer, groupName); err != nil {
		return false, err
	}
	if s.storageServer == nil {
		return false, fmt.Errorf("getStoreStorage fail, errno code: %d", tracker.GetErrorCode())
	}

	return true, nil
}

/**
 * check storage socket, if null create a new connection
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return true if create a new connection
 */
func (s *StorageClient) newReadableStorageConnection(groupName, remoteFilename string) (bool, error) {
	if s.storageServer != nil {
		return false, nil
	}
	var err error
	var tracker = NewTrackerClient()
	if s.storageServer,err = tracker.GetFetchStorage(s.trackerServer, groupName, remoteFilename); err != nil {
		return false, err
	}
	if s.storageServer == nil {
		return false, fmt.Errorf("getStoreStorage fail, errno code: %d", tracker.GetErrorCode())
	}

	return true, nil
}

/**
 * check storage socket, if null create a new connection
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @return true if create a new connection
 */
func (s *StorageClient) newUpdatableStorageConnection(groupName, remoteFilename string) (bool, error) {
	if s.storageServer != nil {
		return false, nil
	}
	var err error
	var tracker = NewTrackerClient()
	if s.storageServer,err = tracker.GetUpdateStorage(s.trackerServer, groupName, remoteFilename); err != nil {
		return false, err
	}
	if s.storageServer == nil {
		return false, fmt.Errorf("getStoreStorage fail, errno code: %d", tracker.GetErrorCode())
	}

	return true, nil
}

/**
 * send package to storage server
 *
 * @param cmd             which command to send
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 */
func (s *StorageClient) sendPackage(cmd byte, groupName, remoteFilename string) error {
	var (
		header []byte
		groupBytes []byte
		filenameBytes []byte
		bs []byte
		groupLen int
	)
	var err error

	groupBytes = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return err
	}
	if filenameBytes,err = ConvertBytesToUTF8([]byte(remoteFilename), GCharset); err != nil {
		return err
	}
	if len(bs) <= len(groupBytes) {
		groupLen = len(bs)
	} else {
		groupLen = len(groupBytes)
	}
	copy(groupBytes[:groupLen], bs)

	if header,err = PackHeader(cmd, int64(len(groupBytes) + len(filenameBytes)), 0); err != nil {
		return err
	}
	var wholePkg = make([]byte, len(header) + len(groupBytes) + len(filenameBytes))
	copy(wholePkg, header)
	copy(wholePkg[len(header):], groupBytes)
	copy(wholePkg[len(header) + len(groupBytes):], filenameBytes)

	conn,err := s.storageServer.GetSocket()
	if err != nil {
		return err
	}
	if _,err = conn.Write(wholePkg); err != nil {
		return err
	}

	return nil
}

/**
 * send package to storage server
 *
 * @param group_name      the group name of storage server
 * @param remote_filename filename on storage server
 * @param file_offset     the start offset of the file
 * @param download_bytes  download bytes
 */
func (s *StorageClient) sendDownloadPackage(groupName, remoteFilename string, fileOffset, downloadBytes int) error {
	var (
		header []byte
		bsOffset []byte
		bsDownBytes []byte
		groupBytes []byte
		filenameBytes []byte
		bs []byte
		groupLen int
	)
	var err error

	bsOffset = Long2Buff(int64(fileOffset))
	bsDownBytes = Long2Buff(int64(downloadBytes))
	groupBytes = make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	if bs,err = ConvertBytesToUTF8([]byte(groupName), GCharset); err != nil {
		return err
	}
	if filenameBytes,err = ConvertBytesToUTF8([]byte(remoteFilename), GCharset); err != nil {
		return err
	}
	if len(bs) <= len(groupBytes) {
		groupLen = len(bs)
	} else {
		groupLen = len(groupBytes)
	}
	copy(groupBytes[:groupLen], bs)

	if header,err = PackHeader(STORAGE_PROTO_CMD_DOWNLOAD_FILE, int64(len(bsOffset) + len(bsDownBytes) + len(groupBytes) + len(filenameBytes)), 0); err != nil {
		return err
	}
	var wholePkg = make([]byte, len(header) + len(bsOffset) + len(bsDownBytes) + len(groupBytes) + len(filenameBytes))
	copy(wholePkg, header)
	copy(wholePkg[len(header):], bsOffset)
	copy(wholePkg[len(header) + len(bsOffset):], bsDownBytes)
	copy(wholePkg[len(header) + len(bsOffset) + len(bsDownBytes):], groupBytes)
	copy(wholePkg[len(header) + len(bsOffset) + len(bsDownBytes) + len(groupBytes):], filenameBytes)

	conn,err := s.storageServer.GetSocket()
	if err != nil {
		return err
	}
	if _,err = conn.Write(wholePkg); err != nil {
		return err
	}

	return nil
}

type UploadBuff struct {
	fileBuff    []byte
	offset      int
	length      int
}

/**
 * constructor
 *
 * @param fileBuff the file buff for uploading
 */
func NewUploadBuff(fileBuff []byte, offset, length int) *UploadBuff {
	return &UploadBuff{
		fileBuff:fileBuff,
		offset:offset,
		length:length,
	}
}

/**
 * send file content callback function, be called only once when the file uploaded
 *
 * @param out output stream for writing file content
 * @return 0 success, return none zero(errno) if fail
 */
func (u *UploadBuff) Send(out io.Writer) (int, error) {
	return out.Write(u.fileBuff[u.offset:u.offset + u.length])
}
