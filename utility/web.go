package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile saves on disk a given file url
func DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer CloseSafely(out)

	// Get the data
	resp, err := http.Get(url) // nolint: gosec
	if err != nil {
		return err
	}
	defer CloseSafely(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
