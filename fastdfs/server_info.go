package fastdfs

import (
	"net"
	"fmt"
	"time"
)

type ServerInfo struct {
	ipAddr     string
	port       int
}

/**
 * Constructor
 *
 * @param ip_addr address of the server
 * @param port    the port of the server
 */
func NewServerInfo(ipAddr string, port int) *ServerInfo {
	return &ServerInfo{
		ipAddr : ipAddr,
		port   : port,
	}
}

/**
 * return the ip address
 *
 * @return the ip address
 */
func (s *ServerInfo) GetIpAddr() string {
	return s.ipAddr
}

/**
 * return the port of the server
 *
 * @return the port of the server
 */
func (s *ServerInfo) GetPort() int {
	return s.port
}

/**
 * connect to server
 *
 * @return connected Socket object
 */
func (s *ServerInfo) Connect() (net.Conn, error) {
	// TODO set address reuse.
	// TODO set read timeout.

	return net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.ipAddr	, s.port), time.Duration(GConnectTimeout) * time.Microsecond)
}