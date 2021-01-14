/*
Copyright (C) 2020 Steffen Fritz <steffen@fritz.wtf>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <https://www.gnu.org/licenses/>.

*/
package main

import (
	"github.com/dla-marbach/filedriller"
	"github.com/gomodule/redigo/redis"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"strconv"
	"strings"
)

// Version holds the version of filedriller
var Version string

// Build holds the sha1 fingerprint of the build
var Build string

func main() {
	var r filedriller.RedisConf


	rootDir := flag.StringP("in", "i", "", "Root directory to work on")
	hashAlg := flag.StringP("algo", "a", "sha256", "The hash algorithm to use: md5, sha1, sha256, sha512, blake2b-512")
	r.Server = flag.StringP("redisserv", "s", "", "Redis server address for a NSRL database")
	r.Port = flag.StringP("redisport", "p", "6379", "Redis port number for a NSRL database")
	sFile := flag.BoolP("download", "d", false, "Download siegfried's signature file")
	oFile := flag.StringP("output", "o", "info.csv", "Output file")
    iFile := flag.StringP("file", "f", "", "Inspect single file")
	entro := flag.BoolP("entropy", "e", false, "Calculate the entropy of files. Limited to file sizes up to 1GB")
	vers := flag.BoolP("version", "v", false, "Print version and build info")


	flag.Parse()

	if *vers {
		log.Printf("Version: %s. Build: %s", Version, Build)
		return
	}

	if len(*iFile) == 0 {
		log.Println("info: friller started")
	}

	if !*sFile {
		if _, err := os.Stat("pronom.sig"); os.IsNotExist(err) {
			log.Println("warning: No pronom.sig file found. Trying to download it.")
			*sFile = true
		}
	}

	if *sFile {

		err := filedriller.DownloadPronom()
		if err != nil {
			log.Println(err)
		}
		log.Println("info: Downloaded pronom.sig file")
		log.Println("info: Please start friller again.")
		log.Println("info: friller ended")
		return
	}

	var nsrlEnabled bool
	var conn redis.Conn
	if *r.Server != "" {
		nsrlEnabled = true
		conn = filedriller.RedisConnect(r)
	}

	// This is the single file processing. All we do is to create and pass a list of length 1
	if len(*iFile) != 0 {
		singleResult := filedriller.IdentifyFiles([]string{*iFile}, *hashAlg, nsrlEnabled, conn, *entro)
		println(singleResult[0])
		return
	}

	if len(*rootDir) == 0 {
		log.Println("error: -in is a mandatory flag")
		return
	}

	if !strings.HasSuffix(*rootDir, "/") {
		*rootDir = *rootDir + "/"
	}

	fileList := filedriller.CreateFileList(*rootDir)
	log.Println("info: Created file list. Found " + strconv.Itoa(len(fileList)) + " files.")
	log.Println("info: Started file format identification")
	resultList := filedriller.IdentifyFiles(fileList, *hashAlg, nsrlEnabled, conn, *entro)
	log.Println("info: Inspected " + strconv.Itoa(len(resultList)) + " files.")

	log.Println("info: Creating output file")
	fd, err := os.Create(*oFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	log.Println("info: Writing output to " + *oFile)

	_, err = fd.WriteString("Filename, SizeInByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, " + strings.ToUpper(*hashAlg) + ", UUID, inNSRL, Entropy\r\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range resultList {
		_, err := fd.WriteString(result + "\r\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("info: Output written to " + *oFile)

	log.Println("info: friller ended")
}
