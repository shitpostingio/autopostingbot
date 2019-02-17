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

// GetPhotoFingerprint calculates a image hash for a given FileID, and returns it as string
func GetPhotoFingerprint(bot *tgbotapi.BotAPI, fileID string) (aHash string, pHash string, err error) {

	imageDownloadURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return
	}

	filePath := fileID + ".jpg"
	err = downloadFile(filePath, imageDownloadURL)
	if err != nil {
		return
	}

	aHash, pHash, err = GetImageFingerprint(filePath)
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

//GetImageFingerprint returns the the pHash fingerprint of a photo given its path
func GetImageFingerprint(filepath string) (aHash string, pHash string, err error) {

	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		log.Println("Unable to open file ", filepath, err)
		return "", "", err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Unable to decode file ", filepath, err)
		return "", "", err
	}

	//We want both AverageHash and PerceptionHash
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		log.Println("Unable to get hash for file: ", filepath, err)
		return "", "", err
	}

	//We return a substring that removes p: from the hash
	aHash = hash.ToString()

	hash, err = goimagehash.PerceptionHash(img)
	if err != nil {
		log.Println("Unable to get hash for file: ", filepath, err)
		return "", "", err
	}

	pHash = hash.ToString()
	return
}

// HasSimilarEnoughPhoto returns true if there is a match between photos and photo,
// and returns the MediaID for the duplicate photo if found.
func HasSimilarEnoughPhoto(dataFunc func() (photoHash string, photosHashes []string)) (bool, string) {
	photoHash, photosHashes := dataFunc()
	for _, currentPhoto := range photosHashes {
		if PhotosAreSimilarEnough(photoHash, currentPhoto) {
			return true, currentPhoto
		}
	}

	return false, ""
}

// PhotosAreSimilarEnough returns true if the two PHashes are close enough (cutoff value: 6)
func PhotosAreSimilarEnough(firstPHash string, secondPHash string) bool {

	pHash1, _ := goimagehash.ImageHashFromString(firstPHash)
	pHash2, _ := goimagehash.ImageHashFromString(secondPHash)
	distance, err := pHash1.Distance(pHash2)
	return (distance < 6 && err == nil)
}

// GetMostSimilarPhoto returns the most similar Photo present in the array passed as argument
func GetMostSimilarPhoto(dataFunc func() (photoHash string, photosHashes []string)) (bestMatch string) {

	photoMatcherHash, photosHashes := dataFunc()
	photoHash, _ := goimagehash.ImageHashFromString(photoMatcherHash)
	closestDistance := 100
	for _, currentPhoto := range photosHashes {
		currentPHash, _ := goimagehash.ImageHashFromString(currentPhoto)
		distance, _ := photoHash.Distance(currentPHash)

		if distance < closestDistance {
			closestDistance = distance
			bestMatch = currentPhoto
		}
	}

	return bestMatch
}
