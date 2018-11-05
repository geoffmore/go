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
	//"strings"
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
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
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

func filesInPath(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	var fileName string
	var hashSum []byte
	var filePath string

	for _, file := range files {
		// completePath := join(filePath, filename, '/')
		//filePath := dir
		fileName = file.Name()
		filePath = dir
		if !file.IsDir() {
			hashSum, _ = fileToMD5(fileName)
			// Printf here doesn't allow for recursion
			fmt.Printf("%s/%s: %x\n",
				filePath,
				fileName,
				hashSum,
			)
		}
	}
}
func main() {
	initFlags()
	fmt.Printf("hash: %s, dir:'%s'\n", hashType, dir)

	startingDir := dir
	filesInPath(startingDir)
}
