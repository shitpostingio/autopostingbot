package database

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/shitposting/loglog-ng"

	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"

	"gitlab.com/shitposting/autoposting-bot/types"
)

// AddImage creates an entry in the database for an image
func AddImage(fileID, caption string, user *entities.User, fingerprint *entities.Fingerprint, db *gorm.DB) bool {

	post := entities.Post{
		User:        *user,
		TypeID:      types.Image,
		FileID:      fileID,
		Fingerprint: fingerprint,
		Caption:     caption,
	}

	result := db.Create(&post)
	if result.RowsAffected != 1 {
		loglog.Err(fmt.Sprintf("Unable to add post to the database. Rows affected: %d", result.RowsAffected))
		return false
	}

	return true
}

// AddVideo creates an entry in the database for a video
func AddVideo(fileID, caption string, user *entities.User, fingerprint *entities.Fingerprint, db *gorm.DB) bool {

	post := entities.Post{
		User:    *user,
		TypeID:  types.Video,
		FileID:  fileID,
		Caption: caption,
	}

	if fingerprint != nil {
		post.Fingerprint = fingerprint
	}

	result := db.Create(&post)
	if result.RowsAffected != 1 {
		loglog.Err(fmt.Sprintf("Unable to add post to the database. Rows affected: %d", result.RowsAffected))
		return false
	}

	return true
}

// AddAnimation creates an entry in the database for a animation
func AddAnimation(fileID, caption string, user *entities.User, fingerprint *entities.Fingerprint, db *gorm.DB) bool {

	post := entities.Post{
		User:    *user,
		TypeID:  types.Animation,
		FileID:  fileID,
		Caption: caption,
	}

	if fingerprint != nil {
		post.Fingerprint = fingerprint
	}

	result := db.Create(&post)
	if result.RowsAffected != 1 {
		loglog.Err(fmt.Sprintf("Unable to add post to the database. Rows affected: %d", result.RowsAffected))
		return false
	}

	return true
}

// UpdatePostCaptionByFileID updates the caption of a post given its fileID
func UpdatePostCaptionByFileID(fileID, caption string, db *gorm.DB) bool {

	post := FindPostByFileID(fileID, db)
	if post.ID == 0 {
		return false
	}

	result := db.Model(&post).Update("caption", caption)
	return result.RowsAffected == 1
}

// FindPostByFileID retrieves a post via its fileID
func FindPostByFileID(fileID string, db *gorm.DB) (post entities.Post) {

	if fileID != "" {
		db.Preload("User").Preload("Fingerprint").Where("file_id = ?", fileID).First(&post)
	}

	return
}

// FindPostByID retrieves a post entity via its database id
func FindPostByID(id uint, db *gorm.DB) (post entities.Post) {

	if id > 0 {
		db.Preload("User").Preload("Fingerprint").Where("id = ?", id).First(&post)
	}

	return
}

// DeletePostByFileID deletes a post entity via its fileID
func DeletePostByFileID(fileID string, db *gorm.DB) error {
	var post entities.Post

	resp := db.Where("file_id = ? AND posted_at IS NULL", fileID).First(&post)
	if resp.RowsAffected != 1 {
		return fmt.Errorf("unable to find post, rows affected: %d", resp.RowsAffected)
	}

	resp = db.Unscoped().Where("file_id = ? AND posted_at IS NULL", fileID).Delete(&entities.Post{})
	if resp.RowsAffected != 1 {
		return fmt.Errorf("unable to delete post, rows affected: %d", resp.RowsAffected)
	}

	return nil
}

// GetNextPost retrieves the oldest media in the queue
func GetNextPost(db *gorm.DB) (entities.Post, error) {

	var post entities.Post
	resp := db.Where("posted_at IS NULL AND deleted_at IS NULL").Not("has_error", 1).First(&post)

	if resp.RowsAffected == 0 {
		return post, errors.New("no element to post was found")
	}

	return post, nil
}

// GetQueueLength returns the number of the enqueued posts
func GetQueueLength(db *gorm.DB) (length int) {
	db.Table("posts").Where("posted_at IS NULL AND deleted_at IS NULL AND has_error = 0").Count(&length)
	return
}

// GetQueuePositionByDatabaseID returns the position of the selected post in the queue
func GetQueuePositionByDatabaseID(id uint, db *gorm.DB) (position int) {
	db.Table("posts").Where("posted_at IS NULL AND deleted_at IS NULL AND has_error = 0 AND id <= ?", id).Count(&position)
	return
}

// MarkPostAsPosted marks a post as posted
func MarkPostAsPosted(post entities.Post, messageID int, db *gorm.DB) error {

	if post.ID == 0 {
		return fmt.Errorf("can't update post with ID 0")
	}

	currentTime := time.Now()
	result := db.Model(&post).Updates(entities.Post{PostedAt: &currentTime, MessageID: messageID})
	if result.RowsAffected != 1 {
		return fmt.Errorf("unable to add messageID %d to post with id %d: %s", messageID, post.ID, result.Error)
	}

	return nil
}

// MarkPostAsFailed marks a post as failed
func MarkPostAsFailed(post entities.Post, db *gorm.DB) error {

	if post.ID == 0 {
		return fmt.Errorf("can't update post with ID 0")
	}

	result := db.Model(&post).Updates(entities.Post{HasError: true})
	if result.RowsAffected != 1 {
		return fmt.Errorf("unable to mark post with id %d as failed: %s", post.ID, result.Error)
	}

	return nil
}
