package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	os.Exit(mainImpl())
}

func mainImpl() int {
	prog := parseFilePath(os.Args[0])
	// PATHEXT holds extensins to treat as a program (.com;.exe;...)
	exts := strings.Split(os.Getenv("PATHEXT"), ";")

	sibling := findSibling(prog, exts)
	if sibling == "" {
		log.Print("No file to launch")
		return 99
	}

	cmd := exec.Command(sibling, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		log.Printf("Failed to launch %v: %+v", sibling, err)
		return 99
	}
	return 0
}

type fileInfo struct {
	Dir  string
	Name string
	Ext  string
}

func parseFilePath(path string) *fileInfo {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]

	return &fileInfo{
		Dir:  dir,
		Name: name,
		Ext:  ext,
	}
}

func findSibling(baseFile *fileInfo, exts []string) string {
	for _, ext := range exts {
		if strings.EqualFold(ext, baseFile.Ext) {
			continue
		}
		file := filepath.Join(baseFile.Dir, baseFile.Name+ext)
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}
		return file
	}
	return ""
}
