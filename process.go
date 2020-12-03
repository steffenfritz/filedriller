package filedriller

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gomodule/redigo/redis"
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

// IdentifyFiles creates metadata with siegfried and hashsum
func IdentifyFiles(fileList []string, hashDigest string, nsrlEnabled bool, conn redis.Conn) []string {
	var wg sync.WaitGroup
	var m sync.RWMutex

	wg.Add(len(fileList))

	var resultList []string
	s, err := siegfried.Load("pronom.sig")
	if err != nil {
		log.Fatal(err)
	}

	var calcNSRL bool
	if hashDigest != "sha1" {
		calcNSRL = true
	}

	for _, filePath := range fileList {
		go func(filePath string) {
			oneFileResult := siegfriedIdent(s, filePath)
			if oneFileResult == "err" {
				log.Println(err)
			}

			onefilehash := hex.EncodeToString(Hashit(filePath, hashDigest))
			oneFile := oneFileResult + ",\"" + onefilehash + "\","
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
			m.Lock()
			resultList = append(resultList, oneFile)
			m.Unlock()
		}(filePath)
	}

	wg.Wait()
	return resultList
}
