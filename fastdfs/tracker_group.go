package fastdfs

import (
	"net"
	"sync"
	"time"
	"fmt"
	"os"
	"runtime/debug"
	"reflect"
)

/**
 * Tracker server group
 *
 * @author Happy Fish / YuQing
 * @version Version 1.17
 */
type TrackerGroup struct {
	TrackerServerIndex     int
	TrackerServers         []net.Addr
	lock                   sync.Mutex
}

/**
 * Constructor
 *
 * @param tracker_servers tracker servers
 */
func NewTrackerGroup(trackerServers []net.Addr) *TrackerGroup {
	var trackerGroup = new(TrackerGroup)
	trackerGroup.TrackerServers = trackerServers
	trackerGroup.TrackerServerIndex = 0

	return trackerGroup
}

/**
 * return connected tracker server
 *
 * @return connected tracker server, null for fail
 */
func (t *TrackerGroup) GetConnectionByIndex(serverIndex int) (*TrackerServer, error) {
	//fmt.Println("timeout:", time.Duration(GConnectTimeout) * time.Microsecond)
	//conn,err := net.Dial("tcp", t.TrackerServers[serverIndex].String())
	conn,err := net.DialTimeout("tcp", t.TrackerServers[serverIndex].String(), time.Duration(GConnectTimeout) * time.Microsecond)
	if err != nil {
		return nil, err
	}

	// TODO set address reused.

	return NewTrackerServer(conn, t.TrackerServers[serverIndex]), nil
}

/**
 * return connected tracker server
 *
 * @return connected tracker server, null for fail
 */
func (t *TrackerGroup) GetConnection() (*TrackerServer, error) {
	var currentIndex int
	t.lock.Lock()
	{
		t.TrackerServerIndex++
		if t.TrackerServerIndex >= len(t.TrackerServers) {
			t.TrackerServerIndex = 0
		}

		currentIndex = t.TrackerServerIndex
	}
	t.lock.Unlock()

	if server,err := t.GetConnectionByIndex(currentIndex); err == nil {
		return server, nil
	} else {
		fmt.Fprintln(os.Stderr, "connect to server " + t.TrackerServers[currentIndex].String() + " fail")
		debug.PrintStack()
	}

	var trackerServer *TrackerServer
	var err error
	for i := 0; i < len(t.TrackerServers); i++ {
		if i == currentIndex {
			continue
		}
		trackerServer,err = t.GetConnectionByIndex(i)
		if err == nil {
			t.lock.Lock()
			if t.TrackerServerIndex == currentIndex {
				t.TrackerServerIndex = i
			}
			t.lock.Unlock()
			return trackerServer, nil
		} else {
			fmt.Fprintln(os.Stderr, "connect to server " + t.TrackerServers[i].String() + " fail")
			debug.PrintStack()
		}
	}

	return trackerServer, err
}

func (t *TrackerGroup) Clone() *TrackerGroup {
	var trackerServers = make([]net.Addr, len(t.TrackerServers))
	for i := 0; i < len(t.TrackerServers); i++ {
		var val = reflect.New(reflect.TypeOf(t.TrackerServers[i]).Elem())
		val.Elem().Set(reflect.ValueOf(t.TrackerServers[i]).Elem())
		trackerServers[i] = val.Interface().(net.Addr)
	}

	return NewTrackerGroup(trackerServers)
}