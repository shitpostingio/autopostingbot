package posting

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
)

func sendLowPostAlerts(postsEnqueued int) {

	//
	cur, err := dbwrapper.GetUsers()
	if err != nil {
		log.Error(err)
		return
	}

	//
	defer func() {
		_ = cur.Close(context.TODO())
	}()

	//
	alert := fmt.Sprintf("🚨 We're running out of posts!\nEnqueued: %d", postsEnqueued)

	//
	for cur.Next(context.TODO()) {

		var user entities.User
		err = cur.Decode(&user)
		if err == nil {
			_, _ = api.SendPlainText(int64(user.TelegramID), alert)
		}

	}

}
