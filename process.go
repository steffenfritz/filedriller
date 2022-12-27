package filedriller

import (
	"encoding/hex"
	"github.com/djherbis/times"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/richardlehane/siegfried"
)

// CreateFileList creates a list of file paths and a directory listing
func CreateFileList(rootDir string) ([]string, []string) {
	var fileList []string
	var dirList []string
	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		} else if info.IsDir() {
			dirList = append(dirList, path)
		}
		return nil
	})
	if err != nil {
		ErrorLogger.Println(err)
	}
	return fileList, dirList
}

// IdentifyDirs reads metadata from the filesystem
func IdentifyFSInfo(entryList []string) {
	for _, entry := range entryList {
		fdinfo, err := os.Stat(entry)
		if err != nil {
			log.Println(err)
		}
		// debug
		println(fdinfo.Mode().Perm())
		//println(fdinfo.ModTime())
		// end debug
	}
}

// IdentifyFiles creates metadata with siegfried and hashsum
func IdentifyFiles(fileList []string, hashDigest string, nsrlEnabled bool, conn redis.Conn, entroEnabled bool) []string {
	var resultList []string
	s, err := siegfried.Load("pronom.sig")
	if err != nil {
		e(err)
	}

	var calcNSRL bool
	if hashDigest != "sha1" {
		calcNSRL = true
	}

	var entroFile float64

	showBar := true

	if len(fileList) < 2 {
		showBar = false
	}

	bar := pb.New(len(fileList))

	if showBar {
		bar.Start()
	}

	for _, filePath := range fileList {
		successful, oneFileResult := siegfriedIdent(s, filePath)
		if !successful {
			ErrorLogger.Println(oneFileResult)
		}

		onefilehash := hex.EncodeToString(Hashit(filePath, hashDigest))
		oneFile := oneFileResult + ",\"" + onefilehash + "\",\"" + CreateUUID() + "\","

		// get atime,mtime and ctime from files
		t, err := times.Stat(filePath)
		if err != nil {
			ErrorLogger.Println(oneFileResult)
		}

		oneFile = oneFile + t.AccessTime().String() + ","
		oneFile = oneFile + t.ModTime().String() + ","

		if t.HasChangeTime() {
			oneFile = oneFile + t.ChangeTime().String() + ","
		} else {
			oneFile = oneFile + ","
		}

		if t.HasBirthTime() {
			oneFile = oneFile + t.BirthTime().String() + ","
		} else {
			oneFile = oneFile + ","
		}

		// we need a sha1 for redis. if sha1 is not used in this run we
		// need to calculate sha1 for redis if nsrl is enabled
		if nsrlEnabled {
			var nsrlHash string
			if calcNSRL {
				nsrlHash = hex.EncodeToString(Hashit(filePath, "sha1"))
			} else {
				nsrlHash = onefilehash
			}

			inNSRL := RedisGet(conn, strings.ToUpper(nsrlHash))
			oneFile = oneFile + inNSRL
		}

		if entroEnabled {
			entroFile, err = entropy(filePath)
			if err == nil {
				oneFile = oneFile + ",\"" + strconv.FormatFloat(entroFile, 'E', -1, 32) + "\""
			} else {
				oneFile = oneFile + ",\"" + "ERROR calculating entropy" + "\""
			}
		} else {
			oneFile = oneFile + ","

		}

		resultList = append(resultList, oneFile)

		if showBar {
			bar.Increment()
		}
	}

	bar.Finish()

	return resultList
}

// IdentifyFilesGUI creates metadata with siegfried and hashsum
// ToDo: This is not the best solution as a lot of code is duplicated.
// IdentifyFiles() should be refactored and split
func IdentifyFilesGUI(fileList []string, nsrlEnabled bool, conf Config, progress *float64) []string {
	var resultList []string

	s, err := siegfried.Load("/Users/steffen/pronom.sig")
	if err != nil {
		e(err)
	}

	var conn redis.Conn
	if conf.RedisServer != "" {
		nsrlEnabled = true
		r := RedisConf{Server: &conf.RedisServer, Port: &conf.RedisPort}
		conn = RedisConnect(r)
	}

	var calcNSRL bool
	if conf.HashAlg != "sha1" {
		calcNSRL = true
	}

	var entroFile float64

	for _, filePath := range fileList {
		successful, oneFileResult := siegfriedIdent(s, filePath)
		if !successful {
			ErrorLogger.Println(oneFileResult)
		}

		onefilehash := hex.EncodeToString(Hashit(filePath, conf.HashAlg))
		oneFile := oneFileResult + ",\"" + onefilehash + "\",\"" + CreateUUID() + "\","

		// we need a sha1 for redis. if sha1 is not used in this run we
		// need to calculate sha1 for redis if nsrl is enabled
		if nsrlEnabled {
			var nsrlHash string
			if calcNSRL {
				nsrlHash = hex.EncodeToString(Hashit(filePath, "sha1"))
			} else {
				nsrlHash = onefilehash
			}

			inNSRL := RedisGet(conn, strings.ToUpper(nsrlHash))
			oneFile = oneFile + inNSRL
		}

		if conf.Entro {
			entroFile, err = entropy(filePath)
			if err == nil {
				oneFile = oneFile + ",\"" + strconv.FormatFloat(entroFile, 'E', -1, 32) + "\""
			} else {
				oneFile = oneFile + ",\"" + "ERROR calculating entropy" + "\""
			}
		} else {
			oneFile = oneFile + ","

		}

		*progress += 1.0
		resultList = append(resultList, oneFile)

	}

	return resultList
}
