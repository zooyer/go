package shmem

import (
	"testing"
	"fmt"
)

// test open or create
func TestOpen(t *testing.T) {
	const name = "test_shared_memory"
	shmem,err := Open(name, 1024, SHMEM_RDWR | SHMEM_CREAT)
	if err != nil {
		panic(err)
	}
	addr,err := shmem.Attach()
	if err != nil {
		panic(err)
	}
	fmt.Println("addr:", addr)
	if err = shmem.Detach(); err != nil {
		panic(err)
	}
	if err = shmem.Close(); err != nil {
		panic(err)
	}
}

// test open only
func TestOpen2(t *testing.T) {
	const name = "test_shared_memory"
	shmem,err := Open(name, 1024, SHMEM_RDWR)
	if err != nil {
		fmt.Println("success, not have this memory.")
		return
	}
	fmt.Println("failed, should not open successfully.")
	if err = shmem.Close(); err != nil {
		panic(err)
	}
}

// test create and open
func TestOpen3(t *testing.T) {
	const name = "test_shared_memory"
	shmem,err := Open(name, 1024, SHMEM_RDWR | SHMEM_CREAT)
	if err != nil {
		panic(err)
	}
	defer shmem.Close()
	shmem2,err := Open(name, 0, SHMEM_RDWR)
	if err != nil {
		panic(err)
	}
	addr,err := shmem2.Attach()
	if err != nil {
		panic(err)
	}
	fmt.Println("addr2:", addr)
	if err = shmem2.Detach(); err != nil {
		panic(err)
	}
	if err = shmem2.Close(); err != nil {
		panic(err)
	}
}

// test create and create only
func TestOpen4(t *testing.T) {
	const name = "test_shared_memory"
	shmem,err := Open(name, 1024, SHMEM_RDWR | SHMEM_CREAT)
	if err != nil {
		panic(err)
	}
	defer shmem.Close()
	shmem2,err := Open(name, 0, SHMEM_RDWR | SHMEM_CREAT | SHMEM_EXCL)
	if err != nil {
		fmt.Println("success, this memory already exists.")
		return
	}
	fmt.Println("failed, should not create successfully.")
	if err = shmem2.Close(); err != nil {
		panic(err)
	}
}
