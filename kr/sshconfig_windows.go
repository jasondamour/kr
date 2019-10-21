package main

import (
	"github.com/kryptco/kr"
	"os"
	"path/filepath"
)

func getPrefix() (string, error) {
	if ex, err := os.Executable(); err == nil {
		return filepath.Dir(ex), nil
	} else {
		PrintErr(os.Stderr, kr.Red("Krypton ▶ Problem getting path of kr.exe"))
		return "", err
	}
}
