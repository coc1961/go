package main

import (
	"os"
	"regexp"
	"testing"
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
		// TODO: Add test cases.
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
		})
	}
}
