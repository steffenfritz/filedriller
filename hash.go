package filedriller

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/blake2b"
)

// Hashit hashes a file using the provided hash algorithm
func Hashit(inFile string, hashalg string) []byte {
	fd, err := os.Open(inFile)
	e(err)
	defer fd.Close()

	var hasher hash.Hash

	if hashalg == "sha256" {
		hasher = sha256.New()

	} else if hashalg == "md5" {
		hasher = md5.New()

	} else if hashalg == "sha1" {
		hasher = sha1.New()

	} else if hashalg == "sha512" {
		hasher = sha512.New()

	} else if hashalg == "blake2b-512" {
		hasher, err = blake2b.New512(nil)

	} else {
		log.Println("Hash not implemented")
		os.Exit(1)
	}

	_, err = io.Copy(hasher, fd)
	io.Copy(hasher, fd)
	e(err)

	checksum := hasher.Sum(nil)

	return checksum
}
