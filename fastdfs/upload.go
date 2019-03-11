package fastdfs

import "io"

type UploadCallback interface {
	/**
	* send file content callback function, be called only once when the file uploaded
	*
	* @param out output stream for writing file content
	* @return 0 success, return none zero(errno) if fail
	*/
	Send(out io.Writer) (int, error)
}

type UploadStream struct {
	inputStream   io.Reader   //input stream for reading
	fileSize      int         //size of the uploaded file
}

/**
 * constructor
 *
 * @param inputStream input stream for uploading
 * @param fileSize    size of uploaded file
 */
func NewUploadStream(inputStream io.Reader, fileSize int) *UploadStream {
	var u = new(UploadStream)
	u.inputStream = inputStream
	u.fileSize = fileSize

	return u
}

/**
 * send file content callback function, be called only once when the file uploaded
 *
 * @param out output stream for writing file content
 * @return 0 success, return none zero(errno) if fail
 */
func (u *UploadStream) Send(out io.Writer) (int, error) {
	var remainBytes = u.fileSize
	var buff = make([]byte, 256 * 1024)
	var bytes int
	var err error
	for remainBytes > 0 {
		var length = 0
		if remainBytes > len(buff) {
			length = len(buff)
		} else {
			length = remainBytes
		}
		if bytes,err = u.inputStream.Read(buff[:length]); err != nil {
			return -1, err
		}

		if _,err = out.Write(buff[:bytes]); err != nil {
			return -1, err
		}
		remainBytes -= bytes
	}

	return 0, nil
}