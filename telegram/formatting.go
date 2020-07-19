package telegram

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

func GetNameFromUser(user *client.User) string {

	if user.LastName == "" {
		return user.FirstName
	}

	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)

}
