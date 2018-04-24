package fingerprinting

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	blockhash "github.com/dsoprea/go-perceptualhash"
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

	bHash := blockhash.NewBlockhash(img, 16)
	fingerprint = bHash.Hexdigest()
	return
}
