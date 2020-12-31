package posting

import (
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/dbwrapper"
	"github.com/shitpostingio/autopostingbot/documentstore/entities"
	l "github.com/shitpostingio/autopostingbot/localization"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	minIntervalBetweenPosts  = 5 * time.Minute
	minIntervalBetweenPauses = 5 * time.Minute
)

// tryPosting tries to post on the channel.
func tryPosting(post *entities.Post) error {

	// Check post time
	if time.Since(m.previousPostTime) <= minIntervalBetweenPosts {
		return fmt.Errorf(l.GetString(l.POSTING_POSTING_PREVIOUS_POST_TOO_CLOSE), time.Since(m.previousPostTime))
	}

	// Prepare caption
	caption := post.Caption
	if GetChannelHandle() != "" && !strings.Contains(post.Caption, "@"+GetChannelHandle()) {
		caption = fmt.Sprintf("%s\n\n@%s", caption, GetChannelHandle())
	}

	ft, err := api.GetFormattedText(caption)
	if err != nil {
		return fmt.Errorf(l.GetString(l.POSTING_POSTING_UNABLE_TO_PARSE_CAPTION), err)
	}

	//
	message, err := api.SendMedia(post.Media.Type, m.config.Autoposting.ChannelID, api.NoReply, post.Media.FileID, ft.Text, ft.Entities)
	if err != nil {
		_ = dbwrapper.MarkPostAsFailed(post)
		return err
	}

	//
	err = dbwrapper.MarkPostAsPosted(post, int(message.Id))
	if err != nil {
		log.Error("Unable to mark post ", post.ID, " as posted")
	}

	//
	schedulePosting(time.Now())

	//
	_ = moveToDirectory(post)
	return err

}

// tryPausing tries pausing the posting.
func tryPausing(duration time.Duration) error {

	if time.Since(m.previousPauseTime) <= minIntervalBetweenPauses {
		return fmt.Errorf(l.GetString(l.POSTING_POSTING_PREVIOUS_PAUSE_TOO_CLOSE), time.Since(m.previousPauseTime))
	}

	//
	newTime := m.nextPostScheduled.Add(duration)
	m.nextPostScheduled = newTime

	//
	if !m.timer.Stop() {
		select {
		case <-m.timer.C:
		default:
		}
	}

	m.timer = time.NewTimer(time.Until(newTime))
	return nil

}

// schedulePosting schedules a new post.
func schedulePosting(postTime time.Time) {

	// Stop the timer and drain the channel if need be
	if !m.timer.Stop() {
		select {
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
	if int(queueLength) < m.config.Autoposting.PostAlertThreshold {
		sendLowPostAlerts(int(queueLength))
	}

}

// postScheduled posts a scheduled media.
func postScheduled() error {

	post, err := dbwrapper.GetNextPost()
	if err != nil {
		return err
	}

	return tryPosting(&post)

}
