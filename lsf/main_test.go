package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_listDir_process(t *testing.T) {
	type fields struct {
		root          string
		filter        *regexp.Regexp
		prefix        string
		posfix        string
		printFile     bool
		printDirs     bool
		readDirFn     func(dirname string) ([]os.FileInfo, error)
		flushFn       func() error
		writeStringfn func(s string) (int, error)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"TestDirs",
			fields{
				root:          "/",
				filter:        nil,
				prefix:        "",
				posfix:        "",
				printFile:     true,
				printDirs:     true,
				readDirFn:     readDirFn,
				flushFn:       flushFn,
				writeStringfn: writeStringfn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			par := &listDir{
				root:          tt.fields.root,
				filter:        tt.fields.filter,
				prefix:        tt.fields.prefix,
				posfix:        tt.fields.posfix,
				printFile:     tt.fields.printFile,
				printDirs:     tt.fields.printDirs,
				readDirFn:     tt.fields.readDirFn,
				flushFn:       tt.fields.flushFn,
				writeStringfn: tt.fields.writeStringfn,
			}
			par.process()
			if result != "/|/file1|/dir1|/dir1/file2|" {
				t.Fatal("Error " + result)
			}
		})
	}
}

var result string

type mokFileInfo struct {
	name  string
	isdir bool
}

func (f *mokFileInfo) Name() string {
	return f.name
}
func (f *mokFileInfo) Size() int64 {
	return 0
}
func (f *mokFileInfo) Mode() os.FileMode {
	return 0
}
func (f *mokFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f *mokFileInfo) IsDir() bool {
	return f.isdir
}
func (f *mokFileInfo) Sys() interface{} {
	return nil
}

func readDirFn(dirname string) (info []os.FileInfo, err error) {
	info = make([]os.FileInfo, 0)
	if dirname == "/dir1" {
		info = append(info, &mokFileInfo{"file2", false})
	} else {
		info = append(info, &mokFileInfo{"file1", false})
		info = append(info, &mokFileInfo{"dir1", true})
	}
	err = nil
	return
}
func flushFn() error {
	return nil
}

func writeStringfn(s string) (int, error) {
	result += strings.Replace(s, "\n", "", 1) + "|"
	// fmt.Print(s)
	// fmt.Println(result)
	return 0, nil
}
