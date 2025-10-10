package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CopyIntoClipboard() error {
	/*fmt.Println("This is stub and it's compiling.")
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString("I am a string literal")
	return cmd.Run()*/
	cumulativeString := ""

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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

		isOfInterest, content, err := isTextFileOfInterest(path)
		if err != nil {
			// Ignore unreadable files
			return nil
		}
		if !isOfInterest {
			return nil
		}
		fileHeader := fmt.Sprintf("----------------------------\n%s\n----------------------------\n", path)

		cumulativeString += fileHeader + content
		return nil
	})
	if err != nil {
		return err
	}
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(cumulativeString)
	return cmd.Run()
}

func isTextFileOfInterest(path string) (bool, string, error) {
	// Skip .git folder
	if strings.HasPrefix(path, ".git") {
		fmt.Println(path, "file NOT included (git check)")
		return false, "", nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return false, "", err
	}

	// Detect MIME type
	mime := http.DetectContentType(data)
	if !strings.HasPrefix(mime, "text/") {
		fmt.Println(path, "file NOT included (http/MIME check)")
		return false, "", nil
	}

	// Check first 1KB for null bytes
	for i := 0; i < len(data) && i < 1024; i++ {
		if data[i] == 0 {
			fmt.Println(path, "file NOT included (defensive check)")
			return false, "", nil
		}
	}

	fmt.Println(path, "file IS included")
	return true, string(data), nil
}
