package filedriller

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"

	"github.com/richardlehane/siegfried"
)

// IdentifyFiles creates metadata with siegfried and some hashing
func CreateFileList(rootDir string) []string {
	var fileList []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return fileList
}

// IdentifyFiles creates metadata with siegfried and some hashing
func IdentifyFiles(fileList []string, hashDigest string, redisconf RedisConf) []string {
	var redisenabled bool

	if len(*redisconf.Server) != 0 {
		redisenabled = true
	}

	var resultList []string
	s, err := siegfried.Load("pronom.sig")
	if err != nil {
		log.Fatal(err)
	}

	for _, filePath := range fileList {
		oneFileResult := siegfriedIdent(s, filePath)
		onefilehash := Hashit(filePath, hashDigest)
		oneFile := oneFileResult + "," + hex.EncodeToString(onefilehash)

		if redisenabled {
			// ToDo: send to redis server
		} else {
			resultList = append(resultList, oneFile)
		}
	}

	return resultList
}
