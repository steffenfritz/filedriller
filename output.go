package filedriller

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

// WriteCSV writes each entry in the result list to a single line of a csv file
func WriteCSV(oFile *string, hashAlg *string, resultList []string) error {
	fd, err := os.Create(*oFile)
	if err != nil {
		return err
	}
	defer fd.Close()


	_, err = fd.WriteString("Filename, SizeInByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, " + strings.ToUpper(*hashAlg) + ", UUID, inNSRL, Entropy\r\n")

	if err != nil {
		return err
	}
	for _, result := range resultList {
		_, err := fd.WriteString(result + "\r\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteLogfile creates a summary log file after the identification run
func WriteLogfile(Version string, Build string, SigFile string, hashAlg string, nsrlEnabled bool, entro bool, fileList []string, resultList []string){
	InfoLogger.Println("friller started")
	InfoLogger.Println("Platform: " + runtime.GOOS + " on " + runtime.GOARCH)
	InfoLogger.Println("Friller Version: " + Version)
	InfoLogger.Println("Friller Build: " + Build)
	InfoLogger.Println("Siegfried signature file: " + SigFile)
	InfoLogger.Println("Hash algorithm used: " + hashAlg)
	InfoLogger.Println("NSRL lookup enabled: " + strconv.FormatBool(nsrlEnabled))
	InfoLogger.Println("Entropy calculation enabled: " + strconv.FormatBool(entro))
	InfoLogger.Println("Inspecting " + strconv.Itoa(len(fileList)) + " files")
	InfoLogger.Println("Inspected " + strconv.Itoa(len(resultList)) + " files.")
}