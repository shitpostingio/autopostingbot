package telegram

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

// GetNameFromUser returns the user's first and last name, if available.
func GetNameFromUser(user *client.User) string {

	if user.LastName == "" {
		return user.FirstName
	}

	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)

}
