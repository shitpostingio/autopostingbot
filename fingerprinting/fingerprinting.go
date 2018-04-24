package fingerprinting

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/corona10/goimagehash"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// GetPhotoFingerprint calculates a image hash for a given FileID, and returns it as string
func GetPhotoFingerprint(bot *tgbotapi.BotAPI, fileID string) (fingerprint string, err error) {

	imageDownloadURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return
	}

	filePath := fileID + ".jpg"
	err = downloadFile(filePath, imageDownloadURL)
	if err != nil {
		return
	}

	fingerprint, err = getImageFingerprint(filePath)
	if err != nil {
		return
	}

	err = os.Remove(filePath)
	return
}

//downloadFile downloads a file using a GET http request
func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//getImageFingerprint returns the the pHash fingerprint of a photo given its path
func getImageFingerprint(filepath string) (fingerprint string, err error) {

	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return "", err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}

	// TODO: maybe we could use both PerceptionHash and AverageHash:
	// 1 - calculate PerceptionHash
	// 2 - calculate AverageHash
	// 3 - calculate the SHA-512 of (PerceptionHash CONCAT AverageHash)
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", err
	}

	//We return a substring that removes p: from the hash
	fingerprint = hash.ToString()
	fingerprint = fingerprint[2:len(fingerprint)]
	return
}
