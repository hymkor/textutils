package main

import "bufio"
import "io"
import "os"
import "regexp"
import "fmt"

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func cat1(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(ansiStrip.ReplaceAllString(scanner.Text(), ""))
	}
}

func main() {
	for _, arg1 := range os.Args[1:] {
		r, err := os.Open(arg1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		cat1(r)
		r.Close()
	}
	if len(os.Args) <= 1 {
		cat1(os.Stdin)
	}
}
