package media

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
	fpcompare "gitlab.com/shitposting/fingerprinting/comparer"
	"gitlab.com/shitposting/telegram-markdown-processor/dbCaption"

	"gitlab.com/shitposting/autoposting-bot/database/database"
	"gitlab.com/shitposting/autoposting-bot/edition"
	"gitlab.com/shitposting/autoposting-bot/repository"
	"gitlab.com/shitposting/autoposting-bot/utility"
)

// HandleNewVideo handles a new video and performs duplicate checks if the flag is active
func HandleNewVideo(msg *tgbotapi.Message, user *entities.User, repo *repository.Repository, checkDuplicates bool) (reply string, duplicatePost entities.Post) {

	/* GET THE FILEID */
	fileID := msg.Video.FileID

	/* TRY TO FIND THE POST BY EXACT FILEID MATCHING */
	post := database.FindPostByFileID(fileID, repo.Db)
	if post.ID != 0 {
		reply = "Same exact fileID already present in the database"
		duplicatePost = post
		return
	}

	/* GET TELEGRAM FILE PATH */
	file, err := repo.Bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		reply = "Unable to get file path from Telegram"
		return
	}

	/* VIDEO FINGERPRINTS ONLY UNTIL A CERTAIN THRESHOLD */
	var fingerprintEntity entities.Fingerprint
	if file.FileSize < repo.Config.FileSizeThreshold {

		/* GET VIDEO FINGERPRINT */
		_, aHash, pHash, err := getVideoFingerprint(&file, repo)
		if err != nil {
			reply = "Unable to get video fingerprint"
			return
		}

		if checkDuplicates { /* CHECK FOR DUPLICATES */
			matchingFingerprints := database.GetFingerprintsByAHash(aHash, repo.Db)

			/* GET PHASH ARRAY */
			pHashes := make([]string, len(matchingFingerprints))
			for index, elem := range matchingFingerprints {
				pHashes[index] = elem.PHash
			}

			/* CHECK IF THERE IS A SIMILAR ENOUGH PHOTO */
			bestPHashMatch := fpcompare.GetMostSimilarPhoto(pHash, pHashes)
			if fpcompare.PhotosAreSimilarEnough(pHash, bestPHashMatch) {
				fingerprintEntity = findPHash(bestPHashMatch, matchingFingerprints)
			}
		}

		/* A MATCHING IMAGE IS ALREADY IN THE DB */
		if fingerprintEntity.ID != 0 {
			post := database.FindPostByID(fingerprintEntity.PostID, repo.Db)
			reply = "Found a matching video!"
			duplicatePost = post
			return
		}

		/* SET FINGERPRINT */
		fingerprintEntity.AHash = aHash
		fingerprintEntity.PHash = pHash
	}

	/* ADD VIDEO TO DATABASE */
	var success bool
	fixedCaption := dbCaption.PrepareCaptionForDB(msg.Caption, edition.ChannelName, utility.GetMessageEntities(msg), 0)
	if fingerprintEntity.AHash != "" && fingerprintEntity.PHash != "" {
		success = database.AddVideo(fileID, fixedCaption, user, &fingerprintEntity, repo.Db, repo.Log)
	} else {
		success = database.AddVideo(fileID, fixedCaption, user, nil, repo.Db, repo.Log)
	}

	if success {
		reply = "Video added correctly!"
	} else {
		reply = "Unable to add video to database"
	}

	return
}

// HandleEditedVideo handles a video whose post has been edited
func HandleEditedVideo(msg *tgbotapi.Message, repo *repository.Repository) string {

	/* GET THE FILEID */
	fileID := msg.Video.FileID

	/* UPDATE THE CAPTION IN THE DB */
	success := database.UpdatePostCaptionByFileID(fileID, dbCaption.PrepareCaptionForDB(msg.Caption, edition.ChannelName, utility.GetMessageEntities(msg), 0), repo.Db)

	if success {
		return "Video edited correctly!"
	}
	return "Unable to edit video in the database"

}
