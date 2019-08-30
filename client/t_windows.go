// +build !linux

package client

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	dll            = syscall.MustLoadDLL("kernel32")
	setConsoleMode = dll.MustFindProc("SetConsoleMode")
)

func SetInputConsoleMode(h syscall.Handle, m uint32) error {
	r, _, err := setConsoleMode.Call(uintptr(h), uintptr(m))
	if r == 0 {
		return err
	}
	return nil
}

func main() {
	h := syscall.Handle(os.Stdin.Fd())
	var m uint32
	if err := syscall.GetConsoleMode(h, &m); err != nil {
		log.Fatal(err)
	}
	if err := SetInputConsoleMode(h, 0); err != nil {
		log.Fatal(err)
	}
	defer SetInputConsoleMode(h, m)

	fmt.Printf("press any key to exit ...")

	b := make([]byte, 10)
	if _, err := os.Stdin.Read(b); err != nil {
		log.Fatal(err)
	}
}
