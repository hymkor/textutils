package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var keyword = map[string]string{
	"arguments":    "Arguments",
	"count":        "Count",
	"createobject": "CreateObject",
	"dim":          "Dim",
	"else":         "Else",
	"end":          "End",
	"explicit":     "Explicit",
	"if":           "If",
	"is":           "Is",
	"item":         "Item",
	"left":         "Left",
	"nothing":      "Nothing",
	"option":       "Option",
	"right":        "Right",
	"set":          "Set",
	"then":         "Then",
	"wscript":      "WScript",
}

var rxQuoteHide = regexp.MustCompile("\"[^\"]*\"")
var rxWord = regexp.MustCompile("\\w+")
var rxQuoteShow = regexp.MustCompile("\a")
var empty = []byte{}

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
		line = rxQuoteHide.ReplaceAllFunc(
			line,
			func(src []byte) []byte {
				dst := make([]byte, 0, len(src)*2)
				for _, ch := range src {
					dst = append(dst, '\a', ch)
				}
				return dst
			})
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
	for _, fname := range os.Args[1:] {
		if err := conv(fname); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", fname, err.Error())
			return
		}
	}
}
