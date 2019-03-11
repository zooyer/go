package fastdfs

import (
	"net"
	"io"
)

type TrackerServer struct {
	conn     net.Conn
	addr     net.Addr
}

/**
 * Constructor
 *
 * @param sock         Socket of server
 * @param inetSockAddr the server info
 */
func NewTrackerServer(conn net.Conn, addr net.Addr) *TrackerServer {
	var trackerServer = new(TrackerServer)
	trackerServer.conn = conn
	trackerServer.addr = addr

	return trackerServer
}

/**
 * get the connected socket
 *
 * @return the socket
 */
func (t *TrackerServer) GetSocket() (net.Conn, error) {
	if t.conn == nil {
		conn,err := GetSocketAddr(t.addr)
		if err != nil {
			return nil, err
		}
		t.conn = conn
	}

	return t.conn, nil
}

/**
  * get the server info
  *
  * @return the server info
  */
func (t *TrackerServer) GetAddress() net.Addr {
	return t.addr
}

func (t *TrackerServer) GetWriter() io.Writer {
	return t.conn
}

func (t *TrackerServer) GetReader() io.Reader {
	return t.conn
}

func (t *TrackerServer) Close() error {
	if t.conn != nil {
		if err := t.conn.Close(); err != nil {
			return err
		}
		t.conn = nil
	}

	return nil
}

// java gc close.
func (t *TrackerServer) finalize() error {
	return t.Close()
}