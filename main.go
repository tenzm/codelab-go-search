package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	recursiveFlag = flag.Bool("r", false, "recursive search: for directories")
)

type ScanResult struct {
	file       string
	lineNumber int
	line       string
}

func scanFile(fpath, pattern string) ([]string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	result := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			result = append(result, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func exit(format string, val ...interface{}) {
	if len(val) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format, val)
		fmt.Println()
	}
	os.Exit(1)
}

func processFile(fpath string, pattern string) {
	res, err := scanFile(fpath, pattern)
	if err != nil {
		exit("Error scanning %s: %s", fpath, err.Error())
	}
	for _, line := range res {
		fmt.Println(fpath+":", line)
	}
}

func processDirectory(dir string, pattern string)  {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir(){
			processFile(path, pattern)
		}
		return nil
	})
	if err != nil{
		panic("Files error")
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		exit("usage: go-search <path> <pattern> to search")
	}

	path := flag.Arg(0)
	pattern := flag.Arg(1)

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	recursive := *recursiveFlag
	if info.IsDir() && !recursive {
		exit("%s: is a directory", info.Name())
	}

	if info.IsDir() && recursive {
		processDirectory(path, pattern)
	} else {
		processFile(path, pattern)
	}
}
