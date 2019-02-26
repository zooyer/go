// +build linux

package shmem

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	/* common mode bits */
	IPC_R          = 000400  /* read permission */
	IPC_W          = 000200  /* write/alter permission */
	IPC_M          = 010000  /* permission to change control info */

	/* SVID required constants (same values as system 5) */
	IPC_CREAT      = 001000  /* create entry if key does not exist */
	IPC_EXCL       = 002000  /* fail if key exists */
	IPC_NOWAIT     = 004000  /* error if request must wait */

	IPC_PRIVATE    = 0       /* private key */

	IPC_RMID       = 0       /* remove identifier */
	IPC_SET        = 1       /* set options */
	IPC_STAT       = 2       /* get options */
)

type IpcPerm struct {
	Key    uint32    /* Key supplied to shmget(2) */
	Uid    uint32    /* Effective UID of owner */
	Gid    uint32    /* Effective GID of owner */
	Cuid   uint32    /* Effective UID of creator */
	Cgid   uint32    /* Effective GID of creator */
	Mode   uint16    /* Permissions + SHM_DEST and SHM_LOCKED flags */
	Seq    uint16    /* Sequence number */
}

type ShmidDs struct {
	ShmPerm      IpcPerm  /* Ownership and permissions */
	ShmSegsz     uint32   /* Size of segment (bytes) */
	ShmAtime     int32    /* Last attach time */
	ShmDtime     int32    /* Last detach time */
	ShmCtime     int32    /* Last change time */
	ShmCpid      int32    /* PID of creator */
	ShmLpid      int32    /* PID of last shmat(2)/shmdt(2) */
	ShmNattch    uint32   /* No. of current attaches */
}

const SHM_KEY = 0x0F

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.EAGAIN:
		return syscall.EAGAIN
	case syscall.EINVAL:
		return syscall.EINVAL
	case syscall.ENOENT:
		return syscall.ENOENT
	}
	return e
}

func Ftok(pathname string, proj_id int) (int, error) {
	info,err := os.Stat(pathname)
	if err != nil {
		return 0, err
	}

	// 8bit
	a := uint8(proj_id)

	stat := info.Sys().(*syscall.Stat_t)

	// 8bit
	b := uint8(stat.Dev)

	// 16bit
	c := uint16(stat.Ino)

	return int(a) << 24 + int(b) << 16 + int(c), nil
}

func Shmget(key int, size uint64, shmflg int) (shmid int, err error) {
	ret,_,e1 := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(shmflg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	shmid = int(ret)

	return
}

func Shmat(shmid int, shmaddr uintptr, shmflg int) (addr uintptr, err error) {
	ret,_,e1 := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), shmaddr, uintptr(shmflg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	addr = ret

	return
}

func Shmdt(shmaddr uintptr) error {
	_,_,e1 := syscall.Syscall(syscall.SYS_SHMDT, shmaddr, 0, 0)
	if e1 != 0 {
		return errnoErr(e1)
	}

	return nil
}

func Shmctl(shmid, cmd int, buf *ShmidDs) error {
	_,_,e1 := syscall.Syscall(syscall.SYS_SHMCTL, uintptr(shmid), uintptr(cmd), uintptr(unsafe.Pointer(buf)))
	if e1 != 0 {
		return errnoErr(e1)
	}

	return nil
}


type sharedMemory struct {
	shmid    int
	addr     uintptr
	mode     int
}

func (s *sharedMemory) Attach() (uintptr, error) {
	addr,err := Shmat(s.shmid, 0, s.mode)
	if err == nil {
		s.addr = addr
	}

	return addr, err
}

func (s *sharedMemory) Detach() error {
	return Shmdt(s.addr)
}

func (s *sharedMemory) Close() error {
	return Shmctl(s.shmid, IPC_RMID, nil)
}

func OpenSharedMemory(name string, size uint64, mode int) (*sharedMemory, error) {
	s := new(sharedMemory)

	var mod int = 0

	if mode & SHMEM_RDWR == SHMEM_RDWR {
		mod |= IPC_R | IPC_W
	}
	if mode & SHMEM_RDONLY == SHMEM_RDONLY {
		mod |= IPC_R
	}
	if mode & SHMEM_WRONLY == SHMEM_WRONLY {
		mod |= IPC_W
	}
	if mode & SHMEM_EXCL == SHMEM_EXCL {
		mod |= IPC_EXCL
	}
	if mode & SHMEM_CREAT == SHMEM_CREAT {
		mod |= IPC_CREAT
	}

	key,err := Ftok(name, SHM_KEY)
	if err != nil {
		return nil, err
	}

	s.mode = mod

	shmid,err := Shmget(key, size, mod | 0666)
	if err != nil {
		return nil, err
	}

	s.shmid = shmid

	return s, nil
}
