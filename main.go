package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
)

const SourceFileName = "tpl"

func main() {

	filepath := os.Args[1]
	archive, err := zip.OpenReader(filepath)

	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	for _, f := range archive.File {

		if f.Name != SourceFileName {
			continue
		}

		dist, err := os.Create(fmt.Sprintf("%s.xlsx", path.Base(filepath)))
		if err != nil {
			log.Fatal(err)
		}

		distPath := path.Join(path.Dir(filepath), dist.Name())

		if err != nil {
			log.Fatal(err)
		}

		archiveFile, err := f.Open()

		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(dist, archiveFile); err != nil {
			log.Fatal(err)
		}

		if err := archiveFile.Close(); err != nil {
			log.Fatal(err)
		}

		if err := dist.Close(); err != nil {
			log.Fatal(err)
		}

		exec.Command("explorer.exe", distPath)

		wg.Add(1)

		go func() {
			<-time.After(2 * time.Second)
			err := os.RemoveAll(distPath)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()

	}

	wg.Wait()
	<-time.After(10 * time.Second)
}
