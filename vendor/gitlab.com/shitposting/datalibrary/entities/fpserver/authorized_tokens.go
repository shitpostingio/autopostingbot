package fpserver

import "github.com/jinzhu/gorm"

// AuthorizedToken is a single token that represents an authorized fpserver user
type AuthorizedToken struct {
	gorm.Model
	Token string
}
