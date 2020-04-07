package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	recursiveFlag = flag.Bool("r", false, "recursive search: for directories")
	lineCount     = flag.Bool("n", false, "show line number")
)

type ScanResult struct {
	file       string
	lineNumber int
	line       string
}

func scanFile(fpath, pattern string) ([]ScanResult, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	result := make([]ScanResult, 0)
	var line_number int
	for scanner.Scan() {
		line := scanner.Text()
		line_number++
		if strings.Contains(line, pattern) {
			result = append(result, ScanResult{file: fpath, lineNumber: line_number, line: line})
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

func processFileThread(fpath string, pattern string, channel chan[]ScanResult) {
	res := processFile(fpath, pattern)
	channel <- res
}

func processFile(fpath string, pattern string) []ScanResult {
	res, err := scanFile(fpath, pattern)
	if err != nil {
		exit("Error scanning %s: %s", fpath, err.Error())
	}
	return res
}


func outputCoincidence(coincidences []ScanResult) {
	for _, line := range coincidences {
		if *lineCount == true {
			fmt.Println(line.file+":"+strconv.Itoa(line.lineNumber)+":", line.line)
		} else {
			fmt.Println(line.file+":", line.line)
		}
	}
}

func processDirectory(dir string, pattern string) []ScanResult{
	found := make([]ScanResult, 0)
	channel := make(chan []ScanResult)

	opened := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			opened += 1
			go processFileThread(path, pattern, channel)
		}
		return nil
	})
	if err != nil {
		panic("Files error")
	}



 	for i := 0; i < opened; i++{
		result := <-channel
		found = append(found, result...)
	}
	return found
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
		outputCoincidence(processDirectory(path, pattern))
	} else {
		outputCoincidence(processFile(path, pattern))
	}
}
