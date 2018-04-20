package fingerprinting

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/corona10/goimagehash"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//handleImageFingerprinting handles the fingerprinting of a photo given its fileID
func handleImageFingerprinting(bot *tgbotapi.BotAPI, fileID string) (fingerprint string, err error) {

	imageDownloadURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Println("Unable to get file direct URL for fileid: ", fileID)
		return
	}

	filePath := fileID + ".jpg"
	err = downloadFile(filePath, imageDownloadURL)
	if err != nil {
		log.Println("Unable to download picture with fileid: ", fileID)
		return
	}

	fingerprint, err = getImageFingerprint(filePath)
	if err != nil {
		log.Println("Unable to get fingerprint for image with fileid: ", fileID)
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
		log.Println("Unable to open file ", filepath, err)
		return "", err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Unable to decode file ", filepath, err)
		return "", err
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		log.Println("Unable to get hash for file: ", filepath, err)
		return "", err
	}

	//We return a substring that removes p: from the hash
	fingerprint = hash.ToString()
	fingerprint = fingerprint[2:len(fingerprint)]
	return
}
