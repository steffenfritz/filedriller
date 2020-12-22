package filedriller

import (
	"github.com/gomodule/redigo/redis"
	"reflect"
	"testing"
)

func BenchmarkCreateFileList(b *testing.B) {
	type args struct {
		rootDir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"File Input List", args{rootDir: "testdata"}, []string{"testdata/1200px-GPLv3_Logo.svg.png", "testdata/emptyfile", "testdata/everywhere.txt", "testdata/test dir/everywhere.txt", "testdata/testDir/everywhere.txt", "testdata/test_dir/everywhere.txt", "testdata/testdir/everywhere.txt", "testdata/testdir/inNSRL/build-classpath", "testdata/textfile.asc", "testdata/töstdir/everywhere.txt"}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			if got := CreateFileList(tt.args.rootDir); !reflect.DeepEqual(got, tt.want) {
				b.Errorf("CreateFileList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIdentifyFiles(b *testing.B) {
	type args struct {
		fileList    []string
		hashDigest  string
		nsrlEnabled bool
		conn        redis.Conn
	}

	conn, _ := redis.Dial("tcp", "127.0.0.1")
	wantString := `"testdata/töstdir/everywhere.txt","2","pronom","x-fmt/111","Plain Text File","","text/plain","text match ASCII","match on text only","3abb6677af34ac57c0ca5828fd94f9d886c26ce59a8ce60ecf6778079423dccff1d6f19cb655805d56098e6d38a1a710dee59523eed7511e5a9e4b8ccb3a4686","0000-0000-0000-0000",`

	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Identify Files", args{fileList: []string{"testdata/töstdir/everywhere.txt"},
			hashDigest:  "sha512",
			nsrlEnabled: false, conn: conn},
			[]string{wantString}}}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			got := IdentifyFiles(tt.args.fileList, tt.args.hashDigest, tt.args.nsrlEnabled, tt.args.conn)
			gotmodified := got[0]
			gotmodin := []string{gotmodified[:264] + ",\"0000-0000-0000-0000\","}
			if !reflect.DeepEqual(gotmodin, tt.want) {
				b.Errorf("IdentifyFiles() = %v, want %v", gotmodin, tt.want)
			}
		})
	}
}
