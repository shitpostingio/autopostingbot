package posting

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/shitposting/autoposting-bot/api"
	"gitlab.com/shitposting/autoposting-bot/documentstore/dbwrapper"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"time"
)

const (
	minIntervalBetweenPosts  = 5 * time.Minute
	minIntervalBetweenPauses = 5 * time.Minute
)

func tryPosting(post *entities.Post) error {

	//

	// Check post time
	if time.Since(m.previousPostTime) <= minIntervalBetweenPosts {
		return fmt.Errorf("only %s has passed since the last post", time.Since(m.previousPostTime))
	}

	// Prepare caption
	caption := fmt.Sprintf("%s\n\n@%s", post.Caption, m.e.GetEditionName())
	ft, err := api.GetFormattedText(caption)
	if err != nil {
		return fmt.Errorf("unable to parse caption: %s", err)
	}

	message, err := api.SendMedia(post.Media.Type, m.config.ChannelID, api.NoReply, post.Media.FileID, ft.Text, ft.Entities)
	if err != nil {
		_ = dbwrapper.MarkPostAsFailed(post)
		return err
	}

	//New PostTime

	//set messageid etc
	err = dbwrapper.MarkPostAsPosted(post, int(message.Id))

	//TODO: CONTROLLARE IL SALVATAGGIO DEI MEME: saranno da spostare

	// update tickers

	// reschedule
	schedulePosting(time.Now())


	//
	err = moveToDirectory(post)


	return err

}

func tryPausing(duration time.Duration) error {

	if time.Since(m.previousPauseTime) <= minIntervalBetweenPauses {
		return fmt.Errorf("only %s has passed since the last pause", time.Since(m.previousPauseTime))
	}

	//
	newTime := m.nextPostScheduled.Add(duration)
	m.nextPostScheduled = newTime

	//
	if !m.timer.Stop() {
		select{
		case <-m.timer.C:
		default:
		}
	}

	m.timer = time.NewTimer(time.Until(newTime))
	return nil

}

func schedulePosting(postTime time.Time) {

	// Stop the timer and drain the channel if need be
	if !m.timer.Stop() {
		select{
		case <-m.timer.C:
		default:
		}
	}

	//
	queueLength := dbwrapper.GetQueueLength()
	newRate := m.e.GetNewPostingRate(int(queueLength))
	m.postingRate = newRate
	m.timer = time.NewTimer(newRate)

	//
	m.previousPostTime = postTime
	m.nextPostScheduled = time.Now().Add(newRate)

	// Send alerts if there are less than X amount of posts enqueued
	if int(queueLength) < m.config.PostAlertThreshold {
		sendLowPostAlerts(int(queueLength))
	}

}

func postScheduled() error {

	post, err := dbwrapper.GetNextPost()
	if err != nil {
		log.Error("postScheduled:", err)
		return err
	}

	return tryPosting(&post)

}
