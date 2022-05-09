package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	getSrc := flag.String("src", "", "Sources of the files ")
	if len(os.Args) < 2 {
		fmt.Println("expected src command")
		os.Exit(1)
	}
	flag.Parse()

	directories := strings.Split(*getSrc, ",")

	zipSource(directories, "testFolder.zip")

}

func zipSource(directories []string, target string) {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	//iterating over all directories
	for _, source := range directories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filepath.Walk(home+source, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					panic(err)
				}

				// 3. Create a local file header
				header, err := zip.FileInfoHeader(info)
				if err != nil {
					panic(err)
				}

				// set compression
				header.Method = zip.Deflate

				// 5. Create writer for the file header and save content of the file
				headerWriter, err := writer.CreateHeader(header)
				if err != nil {
					panic(err)
				}

				if info.IsDir() {
					panic(err)
				}

				f, err := os.Open(path)
				if err != nil {
					panic(err)
				}

				_, err = io.Copy(headerWriter, f)
				f.Close()
				return err
			})

		}()
		wg.Wait()
	}

}
