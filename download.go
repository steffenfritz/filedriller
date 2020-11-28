package filedriller

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// Pronomurl holds the location for my copy of pronom signature file taken from Siegfried
const Pronomurl string = "https://codeberg.org/steffenfritz/filedriller/raw/branch/main/cmd/friller/pronom.sig"

// DownloadPronom downloads a pronom signature file
func DownloadPronom() error {
	response, err := http.Get(Pronomurl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	file, err := os.Create("pronom.sig")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
