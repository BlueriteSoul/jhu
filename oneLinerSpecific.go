package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var jhuConfPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't open home DIR")
		return
	}
	jhuConfPath = filepath.Join(home, ".config", "jhu.conf")
}

func CopySpecificIntoClipboard() error {
	projectPath, files, err := parseJHUConf(jhuConfPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no files listed in config")
	}

	cumulativeString := ""

	for _, file := range files {
		fullPath := filepath.Join(projectPath, file)
		isText, content, err := isTextFileOfInterest(fullPath)
		if err != nil {
			// skip unreadable files
			fmt.Println("Skipping", fullPath, "due to error:", err)
			continue
		}
		if !isText {
			continue
		}
		fileHeader := fmt.Sprintf("----------------------------\n%s\n----------------------------\n", fullPath)
		cumulativeString += fileHeader + content
	}

	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(cumulativeString)
	return cmd.Run()
}

func parseJHUConf(path string) (string, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	lines := strings.Split(string(data), "\n")

	activeBlock := false
	projectPath := ""
	var files []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "#if 1" {
			activeBlock = true
			continue
		} else if line == "#endif" {
			activeBlock = false
			continue
		}
		if !activeBlock || line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "PROJECT_PATH=") {
			projectPath = strings.Trim(line[len("PROJECT_PATH="):], `"`)
		} else if strings.HasPrefix(line, "FILES=") {
			filesStr := strings.Trim(line[len("FILES="):], `"`)
			for _, f := range strings.Split(filesStr, ",") {
				files = append(files, strings.TrimSpace(f))
			}
		}
	}

	if projectPath == "" {
		return "", nil, fmt.Errorf("PROJECT_PATH not found in config")
	}

	return projectPath, files, nil
}
