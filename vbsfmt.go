package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	for scanner := bufio.NewScanner(in); scanner.Scan(); {
		text := rxQuoteHide.ReplaceAllStringFunc(
			scanner.Text(),
			func(str string) string {
				var result bytes.Buffer
				for _, ch := range str {
					result.WriteRune('\a')
					result.WriteRune(ch)
				}
				return result.String()
			})
		text = rxWord.ReplaceAllStringFunc(
			text,
			func(str string) string {
				if result, ok := keyword[strings.ToLower(str)]; ok {
					if result != str {
						isReplaced = true
					}
					return result
				} else {
					return str
				}
			})
		fmt.Fprintln(out, strings.Replace(text, "\a", "", -1))
	}
	in.Close()
	out.Close()
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
}

func main() {
	for _, fname := range os.Args[1:] {
		if err := conv(fname); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", fname, err.Error())
			return
		}
	}
}
