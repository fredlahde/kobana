package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func main() {
	err := syscall.Chroot("/home/lothar/research/chroot/alpine")
	if err != nil {
		panic(err)
	}
	lscmd := exec.Command("ls", "/home/lothar")
	out, err := lscmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
