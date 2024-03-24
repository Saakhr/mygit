package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"

	// Uncomment this block to pass the first stage!
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	gitDir := ".git"
	// You can use print statements as follows for debugging, they'll be visible when running tests.

	// Uncomment this block to pass the first stage!
	//
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	//
	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}
		//
		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}
		//
		fmt.Println("Initialized git directory")
	//
	case "cat-file":
		if len(os.Args) < 4 {

			PrintErrAndExit("Error: illegal git cat-file args\n")

		}
		if len(os.Args[3]) != 40 {
			PrintErrAndExit("Error: illegal git cat-file args\n")

		}

		objectPath := filepath.Join(gitDir, "objects", os.Args[3][:2], os.Args[3][2:])
		switch command := os.Args[2]; command {
		case "-p":
			content := ReadObject(objectPath)
			start := MustIndexByte(content, 0)
			log.Printf("con: %b\nstart:%d", content, start)

			fmt.Printf(string(content[start+1:]))

		case "-t":
			content := ReadObject(objectPath)
			end := MustIndexByte(content, []byte(" ")[0])
			fmt.Printf(string(content[:end]))
		case "-s":
			content := ReadObject(objectPath)
			start := MustIndexByte(content, []byte(" ")[0])
			end := MustIndexByte(content, 0)
			fmt.Printf(string(content[start+1 : end]))

		default:
			fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
			os.Exit(1)
		}
	case "hash-object":

		switch command := os.Args[2]; command {
		case "-w":
			input := os.Args[3]
			file, err := os.Open(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "No file: %s\n", err)
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Println("Error getting file info:", err)
				return
			}
			fileSize := fileInfo.Size()
			data := make([]byte, fileSize)
			_, err = file.Read(data)
			if err != nil {
				PrintErrAndExit("File Not Valid\n")
			}
			sizeString := strconv.FormatInt(fileSize, 10)

			// Create a new byte slice with the prefix "blob [size]\x00" followed by the file data
			prefix := []byte("blob " + sizeString + "\x00")
			finalData := append(prefix, data...)
			hash := sha1.New()
			hash.Write(finalData)
			hashBytes := hash.Sum(nil)

			// Convert the hash to a hexadecimal string
			hashString := hex.EncodeToString(hashBytes)
			fmt.Println(hashString)
			WriteObject(hashString[2:], hashString[:2], finalData)

		default:
			PrintErrAndExit("Command not found!\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
func PrintErrAndExit(err string) {

	fmt.Fprintf(os.Stderr, err)
	os.Exit(1)
}
func ReadObject(objectPath string) []byte {
	file, err := os.Open(objectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}
	return data
}
func WriteObject(objectname, dirv string, data []byte) error {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	err := os.Mkdir(".git/objects/"+dirv, 0755)
	if err != nil {
		PrintErrAndExit("Error: couln't write")
	}
	err = os.WriteFile(".git/objects/"+dirv+"/"+objectname, b.Bytes(), 0644)
	if err != nil {
		PrintErrAndExit("Error: couln't write")
	}

	return nil
}
func MustIndexByte(contents []byte, b byte) int {

	start := bytes.IndexByte(contents, b)

	if start == -1 {

		PrintErrAndExit("Error: illegal object content")

	}

	return start

}
