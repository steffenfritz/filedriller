package filedriller

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/richardlehane/siegfried"
)

// CreateFileList creates a list of file paths
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
func IdentifyFiles(fileList []string, hashDigest string, nsrlEnabled bool, conn redis.Conn, entroEnabled bool) []string {

	var resultList []string
	s, err := siegfried.Load("pronom.sig")
	if err != nil {
		log.Fatal(err)
	}

	var calcNSRL bool
	if hashDigest != "sha1" {
		calcNSRL = true
	}

	var entroFile float64

	for _, filePath := range fileList {
		oneFileResult := siegfriedIdent(s, filePath)
		if oneFileResult == "err" {
			continue
		}

		onefilehash := hex.EncodeToString(Hashit(filePath, hashDigest))
		oneFile := oneFileResult + ",\"" + onefilehash + "\",\"" + CreateUUID() + "\","
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
	}

	return resultList
}
