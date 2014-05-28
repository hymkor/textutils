package main

import "bufio"
import "io"
import "os"
import "regexp"

func cat1(r io.Reader, ansiStrip *regexp.Regexp) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := ansiStrip.ReplaceAllString(scanner.Text(), "")
		os.Stdout.WriteString(text)
		os.Stdout.WriteString("\n")
	}
}

func main() {
	ansiStrip, err := regexp.Compile("\x1B[^a-zA-Z]*[A-Za-z]")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		return
	}
	for _, arg1 := range os.Args[1:] {
		r, err := os.Open(arg1)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
			return
		}
		defer r.Close()
		cat1(r, ansiStrip)
	}
	if len(os.Args) <= 1 {
		cat1(os.Stdin, ansiStrip)
	}
}
