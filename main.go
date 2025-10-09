package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define flags
	locFlag := flag.Bool("loc", false, "Count lines of code in current dir and nested dirs")
	flag.Parse()

	if *locFlag {
		totalLines, err := countLOC(".")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Total LOC: %d\n", totalLines)
	} else {
		fmt.Println("No flag provided. Use -loc to count lines of code.")
	}
}

// countLOC recursively counts lines in all text files under rootDir
func countLOC(rootDir string) (int, error) {
	total := 0

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Skip hidden files or folders (optional)
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		isText, err := isTextFile(path)
		if err != nil {
			// Ignore unreadable files
			return nil
		}
		if !isText {
			return nil
		}

		lines, err := countFileLines(path)
		if err != nil {
			return nil
		}
		total += lines
		return nil
	})

	return total, err
}

// Simple text file heuristic: no null bytes in first 1KB
func isTextFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	// nasty hack to get closer to desired output
	// TBH, this entire thing seems like a nasty hack
	// detecting text files based on null in the first 1kb *facepalm*
	// for "prod" I would need to come up with more robust way of identifying text/code files
	if strings.HasPrefix(path, ".git/") {
		//I don't wanna count git bollocks, that's definetely not part of the codebase
		return false, nil
	}
	defer f.Close()
	data, _ := os.ReadFile(path)
	mime := http.DetectContentType(data)
	if !strings.HasPrefix(mime, "text/") {
		return false, nil
	}

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			fmt.Println(path, false)
			return false, nil
		}
	}
	fmt.Println(path, true)
	return true, nil
}

// Count lines in a single file
func countFileLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}
