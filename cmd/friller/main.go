package main

import (
	"codeberg.org/steffenfritz/filedriller"
	"flag"
	"log"
	"strconv"
)

func main() {
	var r filedriller.RedisConf
	rootDir := flag.String("in", "", "Root directory to work on")
	hashAlg := flag.String("hash", "sha256", "The hash algorithm to use: md5, sha1, sha256, sha512, blake2b-512")
	r.Server = flag.String("redisserv", "", "Redis server address")
	r.Port = flag.String("redisport", "", "Redis port number")
	sFile := flag.Bool("download", false, "Download siegfried's signature file")

	flag.Parse()

	log.Println("info: friller started")
	if *sFile {

		// ToDo: Download file
		return
	}
	if len(*rootDir) == 0 {
		log.Println("error: in is a mandatory flag")
		return
	}

	fileList := filedriller.CreateFileList(*rootDir)
	log.Println("Created file list. Found " + strconv.Itoa(len(fileList)) + " files.")

	resultList := filedriller.IdentifyFiles(fileList, *hashAlg, r)
	log.Println("info: Inspected " + strconv.Itoa(len(resultList)) + " files.")
	log.Println("info: Writing output")
	println(resultList[0])
	log.Println("info: Output written to FILENAME or redis")

	log.Println("info: friller ended")
}
