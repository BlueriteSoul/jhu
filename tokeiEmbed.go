package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed libs/tokei
var embeddedTokei embed.FS

// ensureTokei extracts the embedded binary once per user session
func ensureTokei() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	target := filepath.Join(cacheDir, "jhu_tokei")

	// Already extracted?
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	data, err := fs.ReadFile(embeddedTokei, "libs/tokei")
	if err != nil {
		return "", fmt.Errorf("cannot read embedded tokei: %w", err)
	}
	if err := os.WriteFile(target, data, 0755); err != nil {
		return "", fmt.Errorf("cannot write tokei binary: %w", err)
	}
	return target, nil
}

// runEmbeddedTokei forwards args to the extracted binary
func runEmbeddedTokei(extraArgs []string) error {
	bin, err := ensureTokei()
	if err != nil {
		return err
	}

	cmd := exec.Command(bin, extraArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
