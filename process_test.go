package filedriller

import (
	"reflect"
	"testing"
)

func TestCreateFileList(t *testing.T) {
	type args struct {
		rootDir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"file input list",args{rootDir: "testdata"},[]string{"testdata/1200px-GPLv3_Logo.svg.png",
			"testdata/emptyfile", "testdata/everywhere.txt", "testdata/test dir/everywhere.txt",
			"testdata/testDir/everywhere.txt", "testdata/testDir/inNSRL/build-classpath",
			"testdata/test_dir/everywhere.txt", "testdata/textfile.asc", "testdata/töstdir/everywhere.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFileList(tt.args.rootDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFileList() = %v, want %v", got, tt.want)
			}
		})
	}
}
