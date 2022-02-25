package main

import (
	"io/fs"
	"net/http"
	"strings"
)

func staticHTTPServerWithoutHidingFile() {
	fsys := dotFileHidingFileSystem{http.Dir("/home/go")}
	http.Handle("/", http.FileServer(fsys))
	http.ListenAndServe(":8080", nil)
}

type dotFileHidingFileSystem struct {
	http.FileSystem
}

func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containDotFile(name) {
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, nil
}

func containDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

type dotFileHidingFile struct {
	http.File
}

func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}
