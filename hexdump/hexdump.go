package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	for _, fname := range os.Args[1:] {
		fd, fdErr := os.Open(fname)
		if fdErr != nil {
			fmt.Fprintln(os.Stderr, fdErr.Error())
			return
		}
		defer fd.Close()
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
}
