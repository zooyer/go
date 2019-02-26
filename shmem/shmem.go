package shmem

type Shmem interface {
	Attach() (addr uintptr, err error)
	Detach() error
	Close() error
}

const (
	SHMEM_CREAT     = 0x01
	SHMEM_RDONLY    = 0x02
	SHMEM_WRONLY    = 0x04
	SHMEM_EXCL      = 0x08
	SHMEM_RDWR      = SHMEM_RDONLY | SHMEM_WRONLY
)

func Open(name string, size uint64, mode int) (Shmem, error) {
	return OpenSharedMemory(name, size, mode)
}