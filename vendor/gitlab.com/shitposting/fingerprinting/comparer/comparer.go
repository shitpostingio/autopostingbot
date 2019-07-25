package fpcompare

import (
	"github.com/corona10/goimagehash"
)

// HasSimilarEnoughPhoto returns true if there is a match between photos and photo,
// and returns the MediaID for the duplicate photo if found.
func HasSimilarEnoughPhoto(referencePHash string, pHashes []string) bool {

	for _, currentPHash := range pHashes {
		if PhotosAreSimilarEnough(referencePHash, currentPHash) {
			return true
		}
	}

	return false
}

// PhotosAreSimilarEnough returns true if the two PHashes are close enough (cutoff value: 6)
func PhotosAreSimilarEnough(firstPHash string, secondPHash string) bool {

	pHash1, err := goimagehash.ImageHashFromString(firstPHash)
	if err != nil {
		return false
	}

	pHash2, err := goimagehash.ImageHashFromString(secondPHash)
	if err != nil {
		return false
	}

	distance, err := pHash1.Distance(pHash2)
	return distance < 6 && err == nil
}

// GetMostSimilarPhoto returns the most similar Photo present in the array passed as argument
// or the photo itself if there are no satisfactory matches
func GetMostSimilarPhoto(referencePHashString string, pHashesStrings []string) (bestMatch string) {

	bestMatch = referencePHashString
	closestDistance := 100
	referencePHash, err := goimagehash.ImageHashFromString(referencePHashString)
	if err != nil {
		return
	}

	for _, currentPHashString := range pHashesStrings {
		currentPHash, err := goimagehash.ImageHashFromString(currentPHashString)
		if err != nil {
			continue
		}

		distance, err := referencePHash.Distance(currentPHash)
		if err != nil {
			continue
		}

		if distance < closestDistance {
			closestDistance = distance
			bestMatch = currentPHashString
		}
	}

	return bestMatch
}
