package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/buger/jsonparser"
)

var reader *bufio.Reader

func main() {

	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Not get stat stdout")
		os.Exit(1)
	}

	fmt.Println(info.Mode() & os.ModeCharDevice)

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {

		flag.Parse()
		if flag.NArg() != 1 {
			printUsage()
			os.Exit(1)
		}

		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "LOG OPEN ERROR:", err.Error())
			os.Exit(1)
		}
		defer f.Close()
		reader = bufio.NewReader(f)

	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	buf := bytes.NewBuffer(nil)

	for {
		line, pfx, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, "LOG READ ERROR:", err.Error())
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}

		buf.Write(line)
		if !pfx {
			if buf.Len() > 0 {
				e := new(entry)
				err = jsonparser.ObjectEach(buf.Bytes(), e.parse)
				if err != nil {
					fmt.Fprintln(os.Stderr, "LOG PARSE ERROR:", err.Error(), "\n", string(line))
				} else {
					e.print(os.Stdout)
				}
			}
			buf.Reset()
		}
	}
}

func printUsage() {
	flag.CommandLine.SetOutput(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage: zzz [options] LOG_FILE")
	flag.PrintDefaults()
}
