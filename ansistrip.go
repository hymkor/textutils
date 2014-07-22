package main

import "bufio"
import "io"
import "os"
import "regexp"
import "fmt"

func cat1(r io.Reader, ansiStrip *regexp.Regexp) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(ansiStrip.ReplaceAllString(scanner.Text(), ""))
	}
}

func main() {
	ansiStrip, err := regexp.Compile("\x1B[^a-zA-Z]*[A-Za-z]")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	for _, arg1 := range os.Args[1:] {
		r, err := os.Open(arg1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		cat1(r, ansiStrip)
		r.Close()
	}
	if len(os.Args) <= 1 {
		cat1(os.Stdin, ansiStrip)
	}
}
