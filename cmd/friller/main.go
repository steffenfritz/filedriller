package main

import (
	"codeberg.org/steffenfritz/filedriller"
	"flag"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"strconv"
)

func main() {
	var r filedriller.RedisConf
	rootDir := flag.String("in", "", "Root directory to work on")
	hashAlg := flag.String("hash", "sha256", "The hash algorithm to use: md5, sha1, sha256, sha512, blake2b-512")
	r.Server = flag.String("redisserv", "", "Redis server address for a NSRL database")
	r.Port = flag.String("redisport", "6379", "Redis port number for a NSRL database")
	sFile := flag.Bool("download", false, "Download siegfried's signature file")
	oFile := flag.String("output", "info.txt", "Output file")

	flag.Parse()

	log.Println("info: friller started")
	if *sFile {

		// ToDo: Download file
		return
	}
	if len(*rootDir) == 0 {
		log.Println("error: -in is a mandatory flag")
		return
	}

	var nsrlEnabled bool
	var conn redis.Conn
	if *r.Server != "" {
		nsrlEnabled = true
		conn = filedriller.RedisConnect(r)
	}

	log.Println("info: Creating output file")
	fd, err := os.Create(*oFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	fileList := filedriller.CreateFileList(*rootDir)
	log.Println("info: Created file list. Found " + strconv.Itoa(len(fileList)) + " files.")
	log.Println("info: Started file format identification")
	resultList := filedriller.IdentifyFiles(fileList, *hashAlg, nsrlEnabled, conn)
	log.Println("info: Inspected " + strconv.Itoa(len(resultList)) + " files.")
	log.Println("info: Writing output to " + *oFile)

	fd.WriteString("Filename, SizeinByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, HashSum, UUID, inNSRL\r\n")
	for _, result := range resultList {
		_, err := fd.WriteString(result + "\r\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("info: Output written to " + *oFile)

	log.Println("info: friller ended")
}
