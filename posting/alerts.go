package posting

import (
	"context"
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	"github.com/shitpostingio/autopostingbot/documentstore/entities"
	l "github.com/shitpostingio/autopostingbot/localization"
	log "github.com/sirupsen/logrus"
)

// sendLowPostAlerts sends a message to all the authorized users,
// notifying them of a low number of posts enqueued.
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
	alert := fmt.Sprintf(l.GetString(l.POSTING_ALERTS_LOW_POSTS), postsEnqueued)

	//
	for cur.Next(context.TODO()) {

		var user entities.User
		err = cur.Decode(&user)
		if err == nil {
			_, _ = api.SendPlainText(int64(user.TelegramID), alert)
		}

	}

}
