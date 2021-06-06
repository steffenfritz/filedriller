package filedriller

import (
	"os"
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