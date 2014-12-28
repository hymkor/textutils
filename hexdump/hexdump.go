package main

import (
	"fmt"
	"io"
	"os"
)

func dump(fd io.Reader) {
	var buffer [16]byte
	for {
		n, nErr := fd.Read(buffer[:])
		if nErr != nil {
			if nErr != io.EOF {
				fmt.Fprintln(os.Stderr, nErr.Error())
			}
			break
		}
		for i := 0; i < n; i++ {
			fmt.Printf("%02X ", buffer[i])
		}
		fmt.Println()
	}
}

func main() {
	count := 0
	for _, fname := range os.Args[1:] {
		fd, fdErr := os.Open(fname)
		if fdErr != nil {
			fmt.Fprintln(os.Stderr, fdErr.Error())
			return
		}
		dump(fd)
		fd.Close()
		count++
	}
	if count <= 0 {
		dump(os.Stdin)
	}
}
