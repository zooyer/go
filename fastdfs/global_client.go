package fastdfs

import (
	"errors"
	"net"
	"strings"
	"fmt"
	"strconv"
	"github.com/go/properties"
	"os"
	"time"
)

const (
	ConfKeyConnectTimeout      = "connect_timeout"
	ConfKeyNetworkTimeout      = "network_timeout"
	ConfKeyCharset              = "charset"
	ConfKeyHttpAntiStealToken  = "http.anti_steal_token"
	ConfKeyHttpSecretKey        = "http.secret_key"
	ConfKeyHttpTrackerHttpPort = "http.tracker_http_port"
	ConfKeyTrackerServer        = "tracker_server"
)

const (
	PropKeyConnectTimeoutInSeconds = "fastdfs.connect_timeout_in_seconds"
	PropKeyNetworkTimeoutInSeconds = "fastdfs.network_timeout_in_seconds"
	PropKeyCharset                   = "fastdfs.charset"
	PropKeyHttpAntiStealToken      = "fastdfs.http_anti_steal_token"
	PropKeyHttpSecretKey            = "fastdfs.http_secret_key"
	PropKeyHttpTrackerHttpPort     = "fastdfs.http_tracker_http_port"
	PropKeyTrackerServers           = "fastdfs.tracker_servers"
)

const (
	DefaultConnectTimeout      = 5  //second
	DefaultNetworkTimeout      = 30 //second
	DefaultCharset              = "UTF-8"
	DefaultHttpAntiStealToken  = false
	DefaultHttpSecretKey        = "FastDFS1234567890"
	DefaultHttpTrackerHttpPort = 80
)

var (
	GConnectTimeout = DefaultConnectTimeout * 1000 //millisecond
	GNetworkTimeout = DefaultNetworkTimeout * 1000 //millisecond
	GCharset = DefaultCharset
	GAntiStealToken = DefaultHttpAntiStealToken //if anti-steal token
	GSecretKey = DefaultHttpSecretKey //generage token secret key
	GTrackerHttpPort = DefaultHttpTrackerHttpPort
	GTrackerGroup *TrackerGroup
)

/**
 * load global variables
 *
 * @param conf_filename config filename
 */
func Init(confFilename string) error {
	var iniReader *IniFileReader
	var szTrackerServers []string
	var parts []string

	iniReader,err := NewIniFileReader(confFilename)
	if err != nil {
		return err
	}

	GConnectTimeout = iniReader.GetIntValue("connect_timeout", DefaultConnectTimeout)
	if GConnectTimeout < 0 {
		GConnectTimeout = DefaultConnectTimeout
	}
	GConnectTimeout *= 1000 //millisecond

	GNetworkTimeout = iniReader.GetIntValue("network_timeout", DefaultNetworkTimeout)
	if GNetworkTimeout < 0 {
		GNetworkTimeout = DefaultNetworkTimeout
	}
	GNetworkTimeout *= 1000 //millisecond

	GCharset = iniReader.GetStrValue("charset")
	if GCharset == "" || len(GCharset) == 0 {
		GCharset = "ISO8859-1"
	}

	szTrackerServers = iniReader.GetValues("tracker_server")
	if szTrackerServers == nil {
		return errors.New("item \"tracker_server\" in " + confFilename + " not found")
	}

	var trackerServers = make([]net.Addr, len(szTrackerServers))
	for i := 0; i < len(szTrackerServers); i++ {
		parts = strings.SplitN(szTrackerServers[i], "\\:", 2)
		if len(parts) != 2 {
			return errors.New("the value of item \"tracker_server\" is invalid, the correct format is host:port")
		}

		port,err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return err
		}
		trackerServers[i],err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", strings.TrimSpace(parts[0]), port))
		if err != nil {
			return err
		}
	}

	GTrackerGroup = NewTrackerGroup(trackerServers)

	GTrackerHttpPort = iniReader.GetIntValue("http.tracker_http_port", 80)
	GAntiStealToken = iniReader.GetBoolValue("http.anti_steal_token", false)
	if GAntiStealToken {
		GSecretKey = iniReader.GetStrValue("http.secret_key")
	}

	return nil
}

/**
 * load from properties file
 *
 * @param propsFilePath properties file path, eg:
 *                      "fastdfs-client.properties"
 *                      "config/fastdfs-client.properties"
 *                      "/opt/fastdfs-client.properties"
 *                      "C:\\Users\\James\\config\\fastdfs-client.properties"
 *                      properties文件至少包含一个配置项 fastdfs.tracker_servers 例如：
 *                      fastdfs.tracker_servers = 10.0.11.245:22122,10.0.11.246:22122
 *                      server的IP和端口用冒号':'分隔
 *                      server之间用逗号','分隔
 */
func InitByPropertiesFile(propsFilePath string) error {
	var props = properties.NewProperties()
	file,err := os.Open(propsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = props.Load(file); err != nil {
		return err
	}

	return InitByProperties(props)
}

func InitByProperties(props *properties.Properties) error {
	var trackerServersConf,err = props.GetProperty(PropKeyTrackerServers)
	if err != nil {
		return err
	}
	if trackerServersConf == "" || len(trackerServersConf) == 0 {
		return fmt.Errorf("configure item %s is required", PropKeyTrackerServers)
	}

	if err = InitByTrackers(strings.TrimSpace(trackerServersConf)); err != nil {
		return err
	}

	connectTimeoutInSecondsConf,_ := props.GetProperty(PropKeyConnectTimeoutInSeconds)
	networkTimeoutInSecondsConf,_ := props.GetProperty(PropKeyNetworkTimeoutInSeconds)
	charsetConf,_ := props.GetProperty(PropKeyCharset)
	httpAntiStealTokenConf,_ := props.GetProperty(PropKeyHttpAntiStealToken)
	httpSecretKeyConf,_ := props.GetProperty(PropKeyHttpSecretKey)
	httpTrackerHttpPortConf,_ := props.GetProperty(PropKeyHttpTrackerHttpPort)
	if connectTimeoutInSecondsConf != "" && len(strings.TrimSpace(connectTimeoutInSecondsConf)) != 0 {
		if GConnectTimeout,err = strconv.Atoi(strings.TrimSpace(connectTimeoutInSecondsConf)); err != nil {
			return err
		}
		GConnectTimeout *= 1000
	}
	if networkTimeoutInSecondsConf != "" && len(strings.TrimSpace(networkTimeoutInSecondsConf)) != 0 {
		if GNetworkTimeout,err = strconv.Atoi(strings.TrimSpace(networkTimeoutInSecondsConf)); err != nil {
			return err
		}
		GNetworkTimeout *= 1000
	}
	if charsetConf != "" && len(strings.TrimSpace(charsetConf)) != 0 {
		GCharset = strings.TrimSpace(charsetConf)
	}
	if httpAntiStealTokenConf != "" && len(strings.TrimSpace(httpAntiStealTokenConf)) != 0 {
		if GAntiStealToken,err = strconv.ParseBool(httpAntiStealTokenConf); err != nil {
			return err
		}
	}
	if httpSecretKeyConf != "" && len(strings.TrimSpace(httpSecretKeyConf)) != 0 {
		GSecretKey = strings.TrimSpace(httpSecretKeyConf)
	}
	if httpTrackerHttpPortConf != "" && len(strings.TrimSpace(httpTrackerHttpPortConf)) != 0 {
		if GTrackerHttpPort,err = strconv.Atoi(strings.TrimSpace(httpTrackerHttpPortConf)); err != nil {
			return err
		}
	}

	return nil
}

/**
 * load from properties file
 *
 * @param trackerServers 例如："10.0.11.245:22122,10.0.11.246:22122"
 *                       server的IP和端口用冒号':'分隔
 *                       server之间用逗号','分隔
 */
func InitByTrackers(trackerServers string) error {
	var list = make([]net.Addr, 0, len(trackerServers))
	var spr1 = ","
	var spr2 = ":"
	var arr1 = strings.Split(strings.TrimSpace(trackerServers), spr1)
	for _,addrStr := range arr1 {
		var arr2 = strings.Split(strings.TrimSpace(addrStr), spr2)
		var host = strings.TrimSpace(arr2[0])
		var port,err = strconv.Atoi(strings.TrimSpace(arr2[1]))
		if err != nil {
			return err
		}
		addr,err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return err
		}
		list = append(list, addr)
	}

	return InitByTrackersAddr(list)
}

func InitByTrackersAddr(trackerAddresses []net.Addr) error {
	GTrackerGroup = NewTrackerGroup(trackerAddresses)
	// TODO error
	return nil
}

/**
 * construct Socket object
 *
 * @param ip_addr ip address or hostname
 * @param port    port number
 * @return connected Socket object
*/
func GetSocket(ipAddr string, port int) (net.Conn, error) {
	conn,err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ipAddr, port), time.Duration(GConnectTimeout) * time.Microsecond)
	if err != nil {
		return nil, err
	}
	// TODO set read timeout.

	return conn, nil
}

/**
 * construct Socket object
 *
 * @param addr InetSocketAddress object, including ip address and port
 * @return connected Socket object
 */
func GetSocketAddr(addr net.Addr) (net.Conn, error) {
	conn,err := net.DialTimeout("tcp", addr.String(), time.Duration(GConnectTimeout) * time.Microsecond)
	if err != nil {
		return nil, err
	}

	// TODO set read timeout

	return conn, nil
}

func GetGConnectTimeout() int {
	return GConnectTimeout
}

func SetGConnectTimeout(connectTimeout int) {
	GConnectTimeout = connectTimeout
}

func GetGNetworkTimeout() int {
	return GNetworkTimeout
}

func SetGNetworkTimeout(networkTimeout int) {
	GNetworkTimeout = networkTimeout
}

func GetGCharset() string {
	return GCharset
}

func SetGCharset(charset string) {
	GCharset = charset
}

func GetGTrackerHttpPort() int {
	return GTrackerHttpPort
}

func SetGTrackerHttpPort(trackerHttpPort int) {
	GTrackerHttpPort = trackerHttpPort
}

func GetGAntiStealToken() bool {
	return GAntiStealToken
}

func IsGAntiStealToken() bool {
	return GAntiStealToken
}

func SetGAntiStealToken(antiStealToken bool) {
	GAntiStealToken = antiStealToken
}

func GetGSecretKey() string {
	return GSecretKey
}

func SetGSecretKey(secretKey string) {
	GSecretKey = secretKey
}

func GetGTrackerGroup() *TrackerGroup {
	return GTrackerGroup
}

func SetGTrackerGroup(trackerGroup *TrackerGroup) {
	GTrackerGroup = trackerGroup
}

func ConfigInfo() string {
	var trackerServers = ""
	if GTrackerGroup != nil {
		var trackerAddresses = GTrackerGroup.TrackerServers
		for _,inetSocketAddress := range trackerAddresses {
			if len(trackerServers) > 0 {
				trackerServers += ","
			}
			trackerServers += inetSocketAddress.String()
		}
	}

	return "{" +
		"\n  GConnectTimeout(ms) = " + strconv.Itoa(GConnectTimeout) +
		"\n  GNetworkTimeout(ms) = " + strconv.Itoa(GNetworkTimeout) +
		"\n  GCharset = " + GCharset +
		"\n  GAntiStealToken = " + strconv.FormatBool(GAntiStealToken) +
		"\n  GSecretKey = " + GSecretKey +
		"\n  GTrackerHttpPort = " + strconv.Itoa(GTrackerHttpPort) +
		"\n  trackerServers = " + trackerServers +
		"\n}"
}
