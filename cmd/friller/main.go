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
	"bufio"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"

	fdr "github.com/dla-marbach/filedriller"
	flag "github.com/spf13/pflag"
)

// Version holds the version of filedriller
var Version string

// Build holds the sha1 fingerprint of the build
var Build string

// SigFile holds the download date of the signature file
var SigFile string

// we make this two gloablly available for usage inside imported customized loggers
var logFile = flag.StringP("log", "l", "logs.txt", "Log file")
var errlogFile = flag.StringP("errlog", "w", "errorlogs.txt", "Error log file")

func main() {

	var r fdr.RedisConf

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

	// Create two loggers for two log files, standard and error
	fdr.CreateLogger(*logFile)
	fdr.CreateErrorLogger(*errlogFile)

	if *vers {
		log.Printf("Version: %s. Build: %s. Signature Version: %s", Version, Build, SigFile)
		return
	}

	// Check if single file inspection is not requested
	if len(*iFile) == 0 {
		log.Println("info: friller started")
	}

	// Check if a download is not requested
	if !*sFile {
		if _, err := os.Stat("pronom.sig"); os.IsNotExist(err) {
			log.Println("warning: No pronom.sig file found. Trying to download it.")
			*sFile = true
		}
	}

	// Check if friller should download the pronom.sig file
	if *sFile {

		err := fdr.DownloadPronom()
		if err != nil {
			log.Println(err)
		}
		log.Println("info: Downloaded pronom.sig file")
		log.Println("info: Please restart friller")
		log.Println("info: No log file written")
		log.Println("info: friller ended")
		return
	}

	var nsrlEnabled bool
	var conn redis.Conn
	if *r.Server != "" {
		nsrlEnabled = true
		conn = fdr.RedisConnect(r)
	}

	// This is the single file processing. All we do is to create and pass a list of length 1
	// process item 0 from list and return
	if len(*iFile) != 0 {
		singleResult := fdr.IdentifyFiles([]string{*iFile}, *hashAlg, nsrlEnabled, conn, *entro)
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

	// check if log files are present and if whether they are empty
	if fi, err := os.Stat(*logFile); err == nil && fi.Size() != 0 {
		log.Printf("warning: Log not empty. Quit or append? [Q/a]")
		reader := bufio.NewReader(os.Stdin)
		decision := "Q"
		decision, _ = reader.ReadString('\n')
		decision = strings.TrimSuffix(decision, "\n")

		if decision != "a" {
			log.Println("info: Quitting filedriller")
			return
		}
	}

	// create the custom logger and write startup info
	//fdr.CreateLogger(*logFile)
	fdr.InfoLogger.Println("friller started")
	fdr.InfoLogger.Println("Platform: " + runtime.GOOS + " on " + runtime.GOARCH)
	fdr.InfoLogger.Println("Friller Version: " + Version)
	fdr.InfoLogger.Println("Friller Build: " + Build)
	fdr.InfoLogger.Println("Siegfried signature file: " + SigFile)
	fdr.InfoLogger.Println("Hash algorithm used: " + *hashAlg)

	fdr.InfoLogger.Println("NSRL lookup enabled: " + strconv.FormatBool(nsrlEnabled))
	fdr.InfoLogger.Println("Entropy calculation enabled: " + strconv.FormatBool(*entro))

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

	_, err = fd.WriteString("Filename, SizeInByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, " + strings.ToUpper(*hashAlg) + ", UUID, inNSRL, Entropy\r\n")

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
	log.Println("info: Log file written to " + *logFile)

	// Delete empty error log file
	if fi, err := os.Stat(*errlogFile); err == nil && fi.Size() == 0 {
		err := os.Remove(*errlogFile)
		if err != nil {
			log.Println("error: Could not delete empty error log file.")
			log.Println(err)
		}

	} else {
		log.Println("info: Error log file written to " + *errlogFile)
	}

	log.Println("info: friller ended")
	fdr.InfoLogger.Println("friller stopped")
}
