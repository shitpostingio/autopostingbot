package media

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fpcompare "gitlab.com/shitposting/fingerprinting/comparer"
	"gitlab.com/shitposting/telegram-markdown-processor/dbCaption"

	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/edition"
	"gitlab.com/shitposting/autoposting-bot/repository"
	"gitlab.com/shitposting/autoposting-bot/utility"

	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// HandleNewPhoto handles a new photo and performs duplicate checks if the flag is active
func HandleNewPhoto(msg *tgbotapi.Message, user *entities.User, repo *repository.Repository, checkDuplicates bool) (reply string, duplicatePost entities.Post) {

	/* GET THE FILEID */
	fileID := msg.Photo[len(msg.Photo)-1].FileID

	/* CHECK BY FILEID */
	matchingPost := database.FindPostByFileID(fileID, repo.Db)
	if matchingPost.ID != 0 {
		reply = "The same fileID was already in the database"
		duplicatePost = matchingPost
		return
	}

	/* GET TELEGRAM FILE PATH */
	file, err := repo.Bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		reply = "Unable to get file path from Telegram"
		return
	}

	/* GET PHOTO FINGERPRINT */
	aHash, pHash, err := getPhotoFingerprint(fileID, file.FilePath, repo.Config.Fpserver, repo.Config.BotToken, repo.Log)
	if err != nil {
		reply = "Unable to get photo fingerprint"
		return
	}

	var matchingEntity entities.Fingerprint

	if checkDuplicates {
		/* CHECK FOR DUPLICATES */
		matchingFingerprints := database.GetFingerprintsByAHash(aHash, repo.Db)

		/* GET PHASH ARRAY */
		pHashes := make([]string, len(matchingFingerprints))
		for index, elem := range matchingFingerprints {
			pHashes[index] = elem.PHash
		}

		/* CHECK IF THERE IS A SIMILAR ENOUGH PHOTO */
		bestPHashMatch := fpcompare.GetMostSimilarPhoto(pHash, pHashes)
		if fpcompare.PhotosAreSimilarEnough(pHash, bestPHashMatch) {
			matchingEntity = findPHash(bestPHashMatch, matchingFingerprints)
		}
	}

	/* A MATCHING IMAGE IS ALREADY IN THE DB */
	if matchingEntity.ID != 0 {
		post := database.FindPostByID(matchingEntity.PostID, repo.Db)
		reply = "Found a matching photo!"
		duplicatePost = post
		return
	}

	/* ADD IMAGE TO DATABASE */
	fingerprintEntity := entities.Fingerprint{AHash: aHash, PHash: pHash}
	fixedCaption := dbCaption.PrepareCaptionForDB(msg.Caption, edition.ChannelName, utility.GetMessageEntities(msg), 0)
	success := database.AddImage(fileID, fixedCaption, user, &fingerprintEntity, repo.Db, repo.Log)
	if success {
		reply = "Image added correctly!"
	} else {
		reply = "Unable to add image to database"
	}

	return
}

// HandleEditedPhoto handles a photo whose post has been edited
func HandleEditedPhoto(msg *tgbotapi.Message, repo *repository.Repository) string {

	/* GET THE FILEID */
	fileID := msg.Photo[len(msg.Photo)-1].FileID

	/* UPDATE THE CAPTION IN THE DB */
	success := database.UpdatePostCaptionByFileID(fileID, dbCaption.PrepareCaptionForDB(msg.Caption, edition.ChannelName, utility.GetMessageEntities(msg), 0), repo.Db)

	if success {
		return "Image edited correctly!"
	}
	return "Unable to edited image to database"

}
