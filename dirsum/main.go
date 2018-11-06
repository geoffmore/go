package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var hashType string
var dir string

func initFlags() {
	// Initialize all command line flags to be processed in the go program.
	flag.StringVar(&hashType, "hash", "md5", "hash type of output")
	flag.StringVar(&dir, "dir", ".", "starting directory for traversal")
	flag.Parse()
}

func fileToMD5(text string) ([]byte, error) {
	// Get the MD5 sum from a file
	f, err := os.Open(text)
	// I don't have a good solution for what to return on failure
	// Maybe return both, handle err further down, don't copy if err != nil
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil), nil
}

func fileToSHA256(text string) ([]byte, error) {
	// Get the SHA-256 sum from a file
	f, err := os.Open(text)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil), nil
}

func fileToSHA512(text string) ([]byte, error) {
	// Get the SHA-512 sum from a file
	f, err := os.Open(text)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil), nil
	// Pretty sure crypto/sha512.Sum() isn't documented via godoc
	//return h.Sum512(nil), nil
}

func filesInPath(dir, hash string) {
	// Calculates the hash of all files under a directory
	// Currently fails on dotfiles
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	var hashSum []byte
	var filePath string

	for _, file := range files {
		filePath = fmt.Sprintf("%s/%s", dir, file.Name())
		if !file.IsDir() {
			switch hash {
			case "md5":
				hashSum, _ = fileToMD5(filePath)
			case "sha256":
				hashSum, _ = fileToSHA256(filePath)
			case "sha512":
				hashSum, _ = fileToSHA512(filePath)
			default:
				hashSum, _ = fileToMD5(filePath)
			}
			fmt.Printf("%s: %x\n", filePath, hashSum)
		} else if file.IsDir() {
			filesInPath(filePath, hash)
		}
	}
}

func main() {
	initFlags()
	filesInPath(dir, hashType)
}
