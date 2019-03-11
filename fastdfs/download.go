package fastdfs

import "io"

type DownloadCallback interface {
	/**
	 * recv file content callback function, may be called more than once when the file downloaded
	 *
	 * @param file_size file size
	 * @param data      data buff
	 * @param bytes     data bytes
	 * @return 0 success, return none zero(errno) if fail
	 */
	 Recv(fileSize int, data []byte, bytes int) (n int, err error)
}

type DownloadStream struct {
	out            io.Writer
	currentBytes   int
}

func NewDownLoadStream(out io.Writer) *DownloadStream {
	var download = new(DownloadStream)
	download.out= out

	return download
}

/**
 * recv file content callback function, may be called more than once when the file downloaded
 *
 * @param fileSize file size
 * @param data     data buff
 * @param bytes    data bytes
 * @return 0 success, return none zero(errno) if fail
 */
func (d *DownloadStream) Recv(fileSize int, data []byte, bytes int) (n int, err error) {
	n,err = d.out.Write(data[:bytes])
	if err == nil {
		d.currentBytes += bytes
		if d.currentBytes == fileSize {
			d.currentBytes = 0
		}
	}

	return n, err
}