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

func openFile(fromText string) (*os.File, error) {
	f, err := os.Open(fromText)
	// I don't have a good solution for what to return on failure
	// Maybe return both, handle err further down, don't copy if err != nil
	if err != nil {
		return nil, err
	}
	return f, err
}

func initFlags() {
	// Initialize all command line flags to be processed in the go program.
	flag.StringVar(&hashType, "hash", "md5", "hash type of output")
	// 'dir' needs to be stripped of it's trailing slash if it has one
	flag.StringVar(&dir, "dir", ".", "starting directory for traversal")
	flag.Parse()
}

func fileToMD5(text string) ([]byte, error) {
	// Get the MD5 sum from a file
	f, _ := openFile(text)
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil), nil
}

func fileToSHA256(text string) ([]byte, error) {
	// Get the SHA-256 sum from a file
	f, _ := openFile(text)
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil), nil
}

func fileToSHA512(text string) ([]byte, error) {
	// Get the SHA-512 sum from a file
	f, _ := openFile(text)
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
	// Calculates the hash of all files under a directory using md5/sha256/sha512
	// methods
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	var hashSum []byte
	var filePath string
	// fileName is defined here to grab the first character for later use
	var fileName string
	// Variables are defined outside of the for-loop so they can be referenced
	// below

	var hashFunc func(string) ([]byte, error)
	switch hash {

	case "md5":
		hashFunc = fileToMD5
	case "sha256":
		hashFunc = fileToSHA256
	case "sha512":
		hashFunc = fileToSHA512
	default:
		hashFunc = fileToMD5
	}

	// How do I make this range function asynchronous?
	for _, file := range files {
		fileName = file.Name()
		filePath = fmt.Sprintf("%s/%s", dir, fileName)
		// Ignore anything that starts with a '.' in the file name
		if string(fileName[0]) != "." {
			// Files get their hash calculated
			if !file.IsDir() {
				hashSum, _ = hashFunc(filePath)
				fmt.Printf("%s: %x\n", filePath, hashSum)
				// Directories do not get their hash calculated
			} else if file.IsDir() {
				filesInPath(filePath, hash)
			}
		}
	}
}

func main() {
	initFlags()
	filesInPath(dir, hashType)
}
