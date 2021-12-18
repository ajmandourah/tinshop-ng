package utils

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

type writeCounter struct {
	Total uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc writeCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	log.Printf("\r%s", strings.Repeat(" ", 50)) //nolint:gomnd

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	log.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
// We pass an io.TeeReader into Copy() to report progress on the download.
func DownloadFile(url, filepath string) error {
	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &writeCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	log.Println()

	// Rename the tmp file back to the original file
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
