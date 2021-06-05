package filedriller

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
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
		ErrorLogger.Println(err)
	}
	return fileList
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
