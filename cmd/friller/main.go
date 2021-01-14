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
	"flag"
	fdr "github.com/dla-marbach/filedriller"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Version holds the version of filedriller
var Version string

// Build holds the sha1 fingerprint of the build
var Build string

// SigFile holds the download date of the signature file
var SigFile string

func main() {
	var r fdr.RedisConf
	rootDir := flag.String("in", "", "Root directory to work on")
	hashAlg := flag.String("hash", "sha256", "The hash algorithm to use: md5, sha1, sha256, sha512, blake2b-512")
	r.Server = flag.String("redisserv", "", "Redis server address for a NSRL database")
	r.Port = flag.String("redisport", "6379", "Redis port number for a NSRL database")
	sFile := flag.Bool("download", false, "Download siegfried's signature file")
	oFile := flag.String("output", "info.csv", "Output file")
	logFile := flag.String("log", "logs.txt", "Log file")
	entro := flag.Bool("entropy", false, "Calculate the entropy of files. Limited to file sizes up to 1GB")
	vers := flag.Bool("version", false, "Print version and build info")

	flag.Parse()

	if *vers {
		log.Printf("Version: %s. Build: %s. Signature Version: %s", Version, Build, SigFile)
		return
	}

	log.Println("info: friller started")

	if _, err := os.Stat("pronom.sig"); os.IsNotExist(err) {
		log.Println("warning: No pronom.sig file found. Trying to download it.")
		*sFile = true
	}

	if *sFile {

		err := fdr.DownloadPronom()
		if err != nil {
			log.Println(err)
		}
		log.Println("info: Downloaded pronom.sig file")
		log.Println("info: Please start friller again")
		log.Println("info: No log file written")
		log.Println("info: friller ended")
		return
	}
	if len(*rootDir) == 0 {
		log.Println("error: -in is a mandatory flag")
		return
	}

	if !strings.HasSuffix(*rootDir, "/") {
		*rootDir = *rootDir + "/"
	}

	// create the custom logger and write startup info
	fdr.CreateLogger(*logFile)
	fdr.InfoLogger.Println("friller started")
	fdr.InfoLogger.Println("Platform: " + runtime.GOOS + " on " + runtime.GOARCH)
	fdr.InfoLogger.Println("Friller Version: " + Version)
	fdr.InfoLogger.Println("Friller Build: " + Build)
	fdr.InfoLogger.Println("Siegfried signature file: " + SigFile)
	fdr.InfoLogger.Println("Hash algorithm used: " + *hashAlg)

	var nsrlEnabled bool
	var conn redis.Conn
	if *r.Server != "" {
		nsrlEnabled = true
		conn = fdr.RedisConnect(r)
	}
	fdr.InfoLogger.Println("NSRL enabled: " + strconv.FormatBool(nsrlEnabled))

	fileList := fdr.CreateFileList(*rootDir)
	log.Println("info: Created file list. Found " + strconv.Itoa(len(fileList)) + " files.")
	fdr.InfoLogger.Println("Inspecting " + strconv.Itoa(len(fileList)) + " files")
	log.Println("info: Started file format identification")
	resultList := fdr.IdentifyFiles(fileList, *hashAlg, nsrlEnabled, conn, *entro)
	log.Println("info: Inspected " + strconv.Itoa(len(resultList)) + " files.")
	fdr.InfoLogger.Println("Inspected " + strconv.Itoa(len(resultList)) + " files.")

	log.Println("info: Creating output file")
	fd, err := os.Create(*oFile)
	if err != nil {
		fdr.ErrorLogger.Println(err)
		log.Fatal(err)
	}
	defer fd.Close()

	log.Println("info: Writing output to " + *oFile)

	_, err = fd.WriteString("Filename, SizeInByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, HashSum, UUID, inNSRL, Entropy\r\n")
	if err != nil {
		fdr.ErrorLogger.Println(err)
		log.Fatal(err)
	}
	for _, result := range resultList {
		_, err := fd.WriteString(result + "\r\n")
		if err != nil {
			fdr.ErrorLogger.Println(err)
			log.Fatal(err)
		}
	}

	log.Println("info: Output written to " + *oFile)
	fdr.InfoLogger.Println("Output written to " + *oFile)

	log.Println("info: friller ended")
	fdr.InfoLogger.Println("friller stopped")
}
