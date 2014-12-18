package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var keyword = map[string]string{}
var rxDQuoteHide = regexp.MustCompile("\"[^\"]*\"")
var rxSQuoteHide = regexp.MustCompile("'[^']*'")
var rxWord = regexp.MustCompile("\\w+")
var rxQuoteShow = regexp.MustCompile("\a")
var empty = []byte{}

func insertYenA(src []byte) []byte {
	dst := make([]byte, 0, len(src)*2)
	for _, ch := range src {
		dst = append(dst, '\a', ch)
	}
	return dst
}

func conv(fname string) error {
	in, inErr := os.Open(fname)
	if inErr != nil {
		return inErr
	}
	outFname := fname + ".tmp"
	out, outErr := os.Create(outFname)
	if outErr != nil {
		in.Close()
		return outErr
	}
	isReplaced := false
	reader := bufio.NewReader(in)
	var err error = nil
	for {
		var line []byte
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		line = rxDQuoteHide.ReplaceAllFunc(line, insertYenA)
		line = rxSQuoteHide.ReplaceAllFunc(line, insertYenA)
		line = rxWord.ReplaceAllFunc(
			line,
			func(src []byte) []byte {
				str := string(src)
				if result, ok := keyword[strings.ToLower(str)]; ok {
					if result != str {
						isReplaced = true
					}
					return []byte(result)
				} else {
					return src
				}
			})
		line = rxQuoteShow.ReplaceAll(line, empty)
		_, err = out.Write(line)
		if err != nil {
			break
		}
	}
	in.Close()
	out.Close()
	if err == nil || err != io.EOF {
		if isReplaced {
			bakfname := fname + "~"
			os.Remove(bakfname)
			os.Rename(fname, bakfname)
			os.Rename(outFname, fname)
			fmt.Println(fname)
		} else {
			os.Remove(outFname)
		}
		return nil
	} else {
		os.Remove(outFname)
		return err
	}
}

func main() {
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		if strings.HasPrefix(line, "@") {
			continue
		}
		keyword[strings.ToLower(line)] = line
	}
	for _, arg1 := range os.Args[1:] {
		matches, matchErr := filepath.Glob(arg1)
		if matchErr != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", arg1, matchErr.Error())
			return
		}
		if matches == nil || len(matches) <= 0 {
			matches = []string{arg1}
		}
		for _, fname := range matches {
			if err := conv(fname); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", fname, err.Error())
				return
			}
		}
	}
}
