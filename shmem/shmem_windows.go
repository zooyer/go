// +build windows

package shmem

import (
	"syscall"
	"unsafe"
)

const TRUE    BOOL = 1
const FALSE   BOOL = 0
const NULL    uintptr = 0

type BOOL		int32

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")

func typeToBool(b BOOL) bool {
	if b != FALSE {
		return true
	}

	return false
}
func typeToBOOL(b bool) BOOL {
	if b {
		return TRUE
	}

	return FALSE
}

// errno to error
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.ERROR_IO_PENDING:
		return syscall.ERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

func OpenFileMapping(access uint32, possession bool, name *uint16) (handle syscall.Handle, err error) {
	ret,_,errno := procOpenFileMapping.Call(uintptr(access), uintptr(typeToBOOL(possession)), uintptr(unsafe.Pointer(name)))
	handle = syscall.Handle(ret)
	if handle == 0 {
		if errno.(syscall.Errno) != 0 {
			err = errnoErr(errno.(syscall.Errno))
		} else {
			err = syscall.EINVAL
		}
	}

	return
}


type sharedMemory struct {
	mode      uint32
	addr      uintptr
	handle    syscall.Handle
}

func (s *sharedMemory) Attach() (uintptr, error) {
	addr,err := syscall.MapViewOfFile(s.handle, s.mode, 0, 0, 0)
	if err == nil {
		s.addr = addr
	}

	return addr, err
}

func (s *sharedMemory) Detach() error {
	return syscall.UnmapViewOfFile(s.addr)
}

func (s *sharedMemory) Close() error {
	return syscall.CloseHandle(s.handle)
}

func OpenSharedMemory(name string, size uint64, mode int) (*sharedMemory, error) {
	s := new(sharedMemory)

	var mod uint32 = 0
	var excl, create = false, false

	if mode & SHMEM_RDWR == SHMEM_RDWR {
		mod |= syscall.FILE_MAP_READ | syscall.FILE_MAP_WRITE
	}
	if mode & SHMEM_RDONLY == SHMEM_RDONLY {
		mod |= syscall.FILE_MAP_READ
	}
	if mode & SHMEM_WRONLY == SHMEM_WRONLY {
		mod |= syscall.FILE_MAP_WRITE
	}
	if mode & SHMEM_EXCL == SHMEM_EXCL {
		excl = true
	}
	if mode & SHMEM_CREAT == SHMEM_CREAT {
		create = true
	}

	s.mode = mod

	var err error
	var handle syscall.Handle

	switch {
	case create && excl: {
		var sizeHigh uint32 = 0
		var sizeLow  uint32 = 0
		if size > uint64(^uint32(0)) {
			sizeHigh = uint32(size >> 32)
		}
		sizeLow = uint32(size)
		handle,err = syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_READWRITE, sizeHigh, sizeLow, syscall.StringToUTF16Ptr(name))
	}
	case !create: {
		handle,err = OpenFileMapping(mod, true, syscall.StringToUTF16Ptr(name))
	}
	case create && !excl: {
		handle,err = OpenFileMapping(mod, true, syscall.StringToUTF16Ptr(name))
		if err != nil {
			var sizeHigh uint32 = 0
			var sizeLow  uint32 = 0
			if size > uint64(^uint32(0)) {
				sizeHigh = uint32(size >> 32)
			}
			sizeLow = uint32(size)
			handle,err = syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_READWRITE, sizeHigh, sizeLow, syscall.StringToUTF16Ptr(name))
		}
	}
	}
	if err != nil {
		return nil, err
	}

	s.handle = handle

	return s, nil
}
