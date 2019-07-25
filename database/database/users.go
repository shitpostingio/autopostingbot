package database

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// GetUserByTelegramID retrieves a user via its telegram id
func GetUserByTelegramID(userID int, db *gorm.DB) (user entities.User) {
	db.Where("telegram_id = ?", userID).First(&user)
	return
}

// GetAllUsers returns all the active users in the database
func GetAllUsers(db *gorm.DB) (users []entities.User) {
	db.Where("deleted_at IS NULL").Find(&users)
	return
}

//UserIsAuthorized checks if the user is authorized to perform an action
func UserIsAuthorized(user entities.User) bool {
	return user.ID != 0 && user.DeletedAt == nil
}
