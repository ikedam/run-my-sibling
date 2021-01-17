package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func fixFromWindowsPath(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}
	return strings.ReplaceAll(path, "\\", "/")
}
func TestParseFilePath(t *testing.T) {
	tests := []struct {
		Name    string
		Path    string
		Expects *fileInfo
	}{
		{
			Name: "dir, name, ext",
			Path: fixFromWindowsPath("dir\\name.ext"),
			Expects: &fileInfo{
				Dir:  fixFromWindowsPath("dir\\"),
				Name: "name",
				Ext:  ".ext",
			},
		},
		{
			Name: "subdir, name, ext",
			Path: fixFromWindowsPath(".\\dir\\name.ext"),
			Expects: &fileInfo{
				Dir:  fixFromWindowsPath(".\\dir\\"),
				Name: "name",
				Ext:  ".ext",
			},
		},
		{
			Name: "nodir, name, ext",
			Path: fixFromWindowsPath("name.ext"),
			Expects: &fileInfo{
				Dir:  "",
				Name: "name",
				Ext:  ".ext",
			},
		},
		{
			Name: "dir, name, noext",
			Path: fixFromWindowsPath("dir\\name"),
			Expects: &fileInfo{
				Dir:  fixFromWindowsPath("dir\\"),
				Name: "name",
				Ext:  "",
			},
		},
		{
			Name: "dir, name, multiext",
			Path: fixFromWindowsPath("dir\\name.ext1.ext2"),
			Expects: &fileInfo{
				Dir:  fixFromWindowsPath("dir\\"),
				Name: "name.ext1",
				Ext:  ".ext2",
			},
		},
		{
			Name: "only name",
			Path: "name",
			Expects: &fileInfo{
				Dir:  "",
				Name: "name",
				Ext:  "",
			},
		},
		{
			Name: "blank",
			Path: "",
			Expects: &fileInfo{
				Dir:  "",
				Name: "",
				Ext:  "",
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.Name,
			func(t *testing.T) {
				actual := parseFilePath(test.Path)
				if !reflect.DeepEqual(test.Expects, actual) {
					t.Errorf("Expects: %+v\nActual : %+v", test.Expects, actual)
				}
			},
		)
	}
}

func runInTemporary(name string, f func()) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get the current working directory: %w", err)
	}

	dir, err := ioutil.TempDir("", name)
	if err != nil {
		return fmt.Errorf("Failed to create a temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to chdir to the temporary directory %v: %w", dir, err)
	}
	defer os.Chdir(pwd)

	f()

	return nil
}

func createFiles(files []string) error {
	for _, file := range files {
		dir := filepath.Dir(file)
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("Failed to create %v: %w", dir, err)
		}
		if _, err := os.Stat(file); os.IsNotExist(err) {
			file, err := os.Create(file)
			if err != nil {
				return fmt.Errorf("Failed to create %v: %w", file, err)
			}
			defer file.Close()
		}
	}
	return nil
}

func TestFindSibling(t *testing.T) {
	tests := []struct {
		Name     string
		BaseFile *fileInfo
		Exts     []string
		Files    []string
		Expects  string
	}{
		{
			Name: "normal",
			BaseFile: &fileInfo{
				Dir:  fixFromWindowsPath("dir\\"),
				Name: "name",
				Ext:  ".exe",
			},
			Exts: []string{
				".com",
				".exe",
				".bat",
			},
			Files: []string{
				fixFromWindowsPath("dir\\name.exe"),
				fixFromWindowsPath("dir\\name.bat"),
			},
			Expects: fixFromWindowsPath("dir\\name.bat"),
		},
		{
			Name: "not found",
			BaseFile: &fileInfo{
				Dir:  fixFromWindowsPath("dir\\"),
				Name: "name",
				Ext:  ".exe",
			},
			Exts: []string{
				".com",
				".exe",
				".bat",
			},
			Files: []string{
				fixFromWindowsPath("dir\\name.exe"),
			},
			Expects: "",
		},
	}

	for _, test := range tests {
		t.Run(
			test.Name,
			func(t *testing.T) {
				runInTemporary(
					"TestFindSibling",
					func() {
						if err := createFiles(test.Files); err != nil {
							t.Fatalf("Failed to create files: %+v", err)
						}
						actual := findSibling(test.BaseFile, test.Exts)
						if test.Expects != actual {
							t.Errorf("Expects: %+v\nActual : %+v", test.Expects, actual)
						}
					},
				)
			},
		)
	}
}
