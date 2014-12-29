package main

import "bufio"
import "fmt"
import "io"
import "os"
import "regexp"
import "strconv"
import "time"

import "github.com/zetamatta/nyagos/Src/dos"

func output(line []byte, bom bool) int {
	var line1 string
	if bom {
		line1 = string(line)
	} else {
		line1, _ = dos.AtoU(line)
	}
	fmt.Print(line1)
	return len(line)
}

func tail(reader io.Reader, count int) (int64, bool) {
	var fileSize int64 = 0
	br := bufio.NewReader(reader)
	tailbuf := make([][]byte, count, count)
	i := 0
	bom := false
	for {
		line, err := br.ReadBytes('\n')
		fileSize += int64(len(line))
		if len(line) >= 3 && line[0] == 0xEF && line[1] == 0xBB && line[2] == 0xBF {
			bom = true
			line = line[3:]
		}
		if err != nil {
			j := i - count
			if j < 0 {
				j = 0
			}
			for ; j < i; j++ {
				output(tailbuf[j%count], bom)
			}
			return fileSize, bom
		}
		tailbuf[i%count] = line
		i++
	}
}

func main() {
	rxCounter := regexp.MustCompile("^-(\\d+)")
	tail_f := false
	linecount := 10
	var lastFileName string
	var lastFileSize int64
	bom := false
	for _, arg := range os.Args[1:] {
		if m := rxCounter.FindStringSubmatch(arg); m != nil {
			count1, countErr := strconv.Atoi(m[1])
			if countErr != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, countErr)
				return
			}
			linecount = count1
		} else if arg == "-f" {
			tail_f = true
		} else {
			reader, readerErr := os.Open(arg)
			if readerErr != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, readerErr)
				continue
			}
			lastFileSize, bom = tail(reader, linecount)
			lastFileName = arg
			reader.Close()
		}
	}
	if tail_f {
		for {
			time.Sleep(3)
			stat, err := os.Stat(lastFileName)
			if err != nil {
				switch err.(type) {
				default:
					fmt.Fprintf(os.Stderr, "%s: %s(%T)\n",
						lastFileName,
						err.Error(),
						err)
					return
				case *os.PathError:
					continue
				}
			}
			newSize := stat.Size()
			if newSize == lastFileSize {
				continue
			}
			if newSize < lastFileSize {
				fmt.Fprintf(os.Stderr, "%s was truncated.\n", lastFileName)
				lastFileSize = newSize
				continue
			}
			reader, readerErr := os.Open(lastFileName)
			if readerErr != nil {
				fmt.Fprintf(os.Stderr, "%s: %s(%T)\n",
					lastFileName,
					readerErr.Error(),
					readerErr)
				return
			}
			_, seekErr := reader.Seek(lastFileSize, 0)
			if seekErr != nil {
				fmt.Fprintf(os.Stderr, "%s: Seek error: %s\n",
					lastFileName,
					seekErr.Error())
				reader.Close()
				return
			}
			br := bufio.NewReader(reader)
			for {
				line, err := br.ReadBytes('\n')
				if err != nil {
					break
				}
				lastFileSize += int64(output(line, bom))
			}
			reader.Close()
		}
	}
}
