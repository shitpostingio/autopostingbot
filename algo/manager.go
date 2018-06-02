package algo

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"gitlab.com/shitposting/autoposting-bot/database/entities"
	"gitlab.com/shitposting/autoposting-bot/utility"
	"gitlab.com/shitposting/fingerprinting"
	"gitlab.com/shitposting/loglog/loglogclient"

	"github.com/empetrone/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

// Manager is the central point of input/output for @AntonioBusillo's algorithm.
// It handles:
//  - channel updates
//  - database updates
//  - algorithm lifecycle
type Manager struct {
	botAPI             *tgbotapi.BotAPI
	channelID          int64
	db                 *gorm.DB
	AddImageChannel    chan MediaPayload
	AddVideoChannel    chan MediaPayload
	ModifyMediaChannel chan MediaPayload
	DeleteMediaChannel chan MediaPayload
	StatusChannel      chan MediaPayload
	hourlyPostSignal   <-chan time.Time
	hourlyPostRate     time.Duration
	postSignal         <-chan time.Time
	debug              bool
	log                *loglogclient.LoglogClient
}

// StatusInfo holds informations about the bot's work.
// It helps to monitor the current status of the bot returning informations
// like the number of posts or the posts' rate in an hour.
type StatusInfo struct {
	PostNumber  int64
	PostPerHour string
}

// ManagerConfig is the configuration wanted for a given Manager instance.
// While BotAPIInstance is necessary, DatabasePath is not: if not present,
// Manager will try to load an existing database from ./autopostingbot.db,
// as per config.go.
type ManagerConfig struct {
	DatabaseString string
	BotAPIInstance *tgbotapi.BotAPI
	ChannelID      int64
	Debug          bool
	Log            *loglogclient.LoglogClient
}

// MediaPayload holds informations about who sent an entity, and what was
// the message id.
type MediaPayload struct {
	ChatID    int
	MessageID int
	Entity    entities.Post
}

// Variables holding the two categories we're using, to distinguish
// between images and video media types.
var (
	videoCategory entities.Category
	imageCategory entities.Category
)

// NewManager returns a new Manager instance
func NewManager(mc ManagerConfig) (m *Manager, err error) {
	if mc.ChannelID == 0 {
		err = errors.New("ChannelID is empty")
		return
	}

	mm := Manager{
		botAPI:             mc.BotAPIInstance,
		channelID:          mc.ChannelID,
		AddImageChannel:    make(chan MediaPayload),
		AddVideoChannel:    make(chan MediaPayload),
		ModifyMediaChannel: make(chan MediaPayload),
		DeleteMediaChannel: make(chan MediaPayload),
		StatusChannel:      make(chan MediaPayload),
		debug:              mc.Debug,
		log:                mc.Log,
	}

	// Initialize gorm
	mm.db, err = gorm.Open("mysql", mc.DatabaseString)
	if err != nil {
		return
	}

	// Get and initialize the categories
	mm.db.Where("name = ?", "image").First(&imageCategory)
	mm.db.Where("name = ?", "video").First(&videoCategory)
	if imageCategory.Name != "image" || videoCategory.Name != "video" {
		err = errors.New("cannot load video and/or image categories identities from the database")
		return
	}

	// Calculate the hourly post rate on the current post availability
	mm.calculateHourlyPostRate()

	// Print the hourly posting rate in minutes
	mm.log.Info("Initial hourly posting rate set to " + mm.hourlyPostRate.String())

	// Initialize the calculation signal
	mm.hourlyPostSignal = time.After(1 * time.Hour)

	// Initialize the postSignal on the hourlyRate
	mm.setUpPostSignal()

	// Start the manager lifecycle
	go mm.managerLifecycle()

	m = &mm
	return
}

// managerLifecycle is the function we run indefinitely in a goroutine.
// It handles incoming updates, and the posting routine.
func (m *Manager) managerLifecycle() {

	// if -debug is specified, immediately send a post and exit
	if m.debug {
		err := m.doPost()
		if err != nil {
			m.log.Err(err.Error())
		}
		os.Exit(0)
	}

	for {
		select {
		case deletedPost := <-m.DeleteMediaChannel:

			var post entities.Post

			m.db.Where("media = ? AND isnull(posted_at)", deletedPost.Entity.Media).First(&post)

			// Since I have the var post, it has ID 0 if the query doesn't return any value, so I check the ID.
			if post.ID == 0 {
				m.log.Warn("Can't delete post. Probably File ID is invalid or it was already posted on channel")
				utility.SendTelegramReply(int(deletedPost.ChatID), deletedPost.MessageID, m.botAPI, "I can't delete it")
			} else {
				m.db.Delete(&post)
				m.log.Info(fmt.Sprintf("Deleted post with ID: %d, deleted by: %d", int(post.ID), deletedPost.Entity.UserID))
				utility.SendTelegramReply(int(deletedPost.ChatID), deletedPost.MessageID, m.botAPI, "Deleted! \xF0\x9F\x9A\xAE")
			}
		case status := <-m.StatusChannel:

			s := m.GetStatus()

			msgText := fmt.Sprintf("\xF0\x9F\x95\x9C Post rate: %s \n\xF0\x9F\x93\x8B Memes enqueued: %d \n \n \n\xE2\x9E\xA1 You're Welcome my ni\xF0\x9F\x85\xB1\xF0\x9F\x85\xB1a", s.PostPerHour, s.PostNumber)

			utility.SendTelegramReply(status.ChatID, status.MessageID, m.botAPI, msgText)

		case modifiedPost := <-m.ModifyMediaChannel:
			var entity entities.Post
			id, err := getUserID(m.db, int(modifiedPost.Entity.UserID))
			if err != nil {
				m.log.Err(err.Error())
				continue
			}
			k := m.db.Where("media = ? AND user_id = ?", modifiedPost.Entity.Media, id).First(&entity)

			if entity.Media == "" { // an empty media ID means no entity with said ID was found
				m.log.Warn(fmt.Sprintf("someone tried to update the caption for media with id %s, but i don't know any", modifiedPost.Entity.Media))
				continue
			}

			k.Update("caption", modifiedPost.Entity.Caption)
			utility.SendTelegramReply(modifiedPost.ChatID, modifiedPost.MessageID, m.botAPI, "Modified!")
		case newPost := <-m.AddVideoChannel:
			m.log.Info("got a new video to add!")

			// if we have a duplicate, write a log message and stop
			if m.checkDuplicate(newPost) {
				continue
			}

			userID, err := getUserID(m.db, int(newPost.Entity.UserID))
			if err != nil {
				m.log.Warn(err.Error())
				continue
			}

			newPost.Entity.UserID = uint(userID)
			newPost.Entity.Categories = []entities.Category{videoCategory}

			// add to the database
			m.db.Create(&newPost.Entity)
			utility.SendTelegramReply(newPost.ChatID, newPost.MessageID, m.botAPI, "Video added!")
		case newPost := <-m.AddImageChannel:
			m.log.Info("got a new image to add!")

			// if we have a duplicate, write a log message and stop
			if m.checkDuplicate(newPost) {
				continue
			}

			userID, err := getUserID(m.db, int(newPost.Entity.UserID))
			if err != nil {
				m.log.Warn(err.Error())
				continue
			}

			newPost.Entity.UserID = uint(userID)
			newPost.Entity.Categories = []entities.Category{imageCategory}

			aHash, pHash, err := fingerprinting.GetPhotoFingerprint(m.botAPI, newPost.Entity.Media)
			if err != nil {
				m.log.Warn(fmt.Sprintf("cannot calculate hash for image with ID %s, proceeding without one", newPost.Entity.Media))
			}

			newPost.Entity.PHash = pHash
			newPost.Entity.AHash = aHash

			// add to the database
			m.db.Create(&newPost.Entity)
			utility.SendTelegramReply(newPost.ChatID, newPost.MessageID, m.botAPI, "Image added!")
		case <-m.postSignal:
			err := m.doPost()
			if err != nil {
				m.log.Err(err.Error())
			}
		case <-m.hourlyPostSignal:
			m.log.Info("calculating the hourly posting rate...")
			lastPostingRate := m.hourlyPostRate
			// calculate the new hourly post rate
			m.calculateHourlyPostRate()

			// set up the post signal if the last hourlyPostSignal was zero
			// and only if lastPostingRate is not 1
			if lastPostingRate <= 0 {
				m.setUpPostSignal()
			}

			m.log.Info(fmt.Sprintf("new hourly posting rate: %s", m.hourlyPostRate.String()))

			// see you in an hour!
			m.hourlyPostSignal = time.After(1 * time.Hour)
		}
	}
}

func (m *Manager) doPost() error {
	// setup the post signal first
	defer m.setUpPostSignal()

	m.log.Info("it's time to post!")

	// could not find anything to post
	wtp, err := m.whatToPost()
	if err != nil {
		return err
	}

	if err := m.popAndPost(wtp); err != nil {
		// posting did not go well...
		// mark that media with the error flag
		wtp.HasError = true
		m.db.Save(&wtp)

		return fmt.Errorf("%s on media with ID %s", err, wtp.Media)
	}

	m.log.Info("all done!")

	return nil
}

// getUserID gets the database user ID for each Telegram user
func getUserID(db *gorm.DB, id int) (int, error) {
	// find the user id based on the telegram id
	var user entities.User
	db.Where("telegram_id = ?", id).First(&user)

	// if no user with said telegram_id was found, throw an error
	if user.ID == 0 {
		return 0, fmt.Errorf("cannot find user_id for telegram_id %d", id)
	}

	// return the correct user ID
	return int(user.ID), nil
}

// calculateHourlyPostRate calculate the hourly post rate, and saves it in the Manager
// instance.
func (m *Manager) calculateHourlyPostRate() {
	var postsQueue []entities.Post
	m.db.Not("has_error", 1).Find(&postsQueue)
	postsQueue = cleanFromPosted(postsQueue)

	ppH := postsPerHour(postsQueue)

	if m.debug {
		m.log.Info(fmt.Sprintf("posts per hour: %d", ppH))
	}

	if ppH > 0 {
		hourlyRate := 60 / ppH
		if hourlyRate <= 0 {
			m.hourlyPostRate = time.Duration(1) * time.Minute
		} else {
			m.hourlyPostRate = time.Duration(hourlyRate) * time.Minute
		}

		if m.debug {
			m.log.Info(fmt.Sprintf("hourly post rate: %s", m.hourlyPostRate))
		}
		return
	}
	m.hourlyPostRate = time.Duration(0) * time.Minute
}

// setUpPostSignal sets up the posting signal if there's something to post
func (m *Manager) setUpPostSignal() {
	if m.hourlyPostRate != 0 {
		m.postSignal = time.After(m.hourlyPostRate)
	}
}

// whatToPost returns the oldest media in the queue
func (m *Manager) whatToPost() (entities.Post, error) {
	var postsQueue []entities.Post
	m.db.Preload("Categories").Not("has_error", 1).Find(&postsQueue)
	postsQueue = cleanFromPosted(postsQueue)

	sort.Sort(entities.Posts(postsQueue))

	if len(postsQueue) <= 0 {
		return entities.Post{}, errors.New("no element to post has been found")
	}

	return postsQueue[0], nil
}

// popAndPost removes entity from the database and post it to the Telegram channel
// identified by Manager.channelID
func (m *Manager) popAndPost(entity entities.Post) error {
	caption := "@shitpost"
	if entity.Caption != "" {
		entity.Caption = strings.TrimSpace(entity.Caption)
		entity.Caption = strings.Replace(entity.Caption, "@Shitpost", "", -1)
		entity.Caption = strings.Replace(entity.Caption, "@shitpost", "", -1)
		caption = fmt.Sprintf("%s\n\n@shitpost", entity.Caption)
	}

	var err error
	var sentMessage tgbotapi.Message
	switch {
	case entity.IsImage(m.db):
		tgImage := tgbotapi.NewPhotoShare(m.channelID, entity.Media)
		tgImage.Caption = caption
		sentMessage, err = m.botAPI.Send(tgImage)
	case entity.IsVideo(m.db):
		tgVideo := tgbotapi.NewVideoShare(m.channelID, entity.Media)
		tgVideo.Caption = caption
		sentMessage, err = m.botAPI.Send(tgVideo)
	}

	// checking if there's an error here gives us the chance to remove the posted
	// entity if everything was ok - if it wasn't, the error will be handled on the caller
	// level.
	if err == nil {
		entity.PostedAt = time.Now()
		entity.MessageID = sentMessage.MessageID
		m.db.Save(&entity)
	}
	return err
}

// isDuplicate returns true if post has been already added before
// false otherwise.
func (m Manager) isDuplicate(post entities.Post) (bool, string, error) {

	var videoDuplicate entities.Post
	var hasSimilar bool
	var dupePhoto string

	newPostAhash, newPostPhash, err := fingerprinting.GetPhotoFingerprint(m.botAPI, post.Media)
	if err != nil {
		return false, dupePhoto, err
	}

	var photos []entities.Post
	m.db.Select("id, media, p_hash").Where("a_hash = ?", newPostAhash).Find(&photos)
	photosPhash := make([]string, len(photos))
	for index, elem := range photos {
		photosPhash[index] = elem.PHash
	}

	// populate the post with all the data we have in the database, if any
	if post.IsImage(m.db) {
		hasSimilar, dupePhoto = fingerprinting.HasSimilarEnoughPhoto(func() (string, []string) {
			return newPostPhash, photosPhash
		})
	} else {
		m.db.Where("media = ?", post.Media).First(&videoDuplicate)
	}

	for _, elem := range photos {
		if elem.PHash == dupePhoto {
			dupePhoto = elem.Media
		}
	}

	if videoDuplicate.Media != "" || hasSimilar {
		return true, dupePhoto, nil
	}

	return false, dupePhoto, nil
}

// checkDuplicate checks whether post has already been added to the database,
// and if yes, it will communicate it to the user
func (m Manager) checkDuplicate(post MediaPayload) bool {
	dup, dupImg, err := m.isDuplicate(post.Entity)

	if err != nil {
		m.log.Err(fmt.Sprintf("error while trying to calculate hash for image with ID %s: %s", post.Entity.Media, err.Error()))
		return dup
	}

	if dup {
		m.log.Warn(fmt.Sprintf("user %d tried to re-add media %s, which is already present in the database", post.Entity.UserID, post.Entity.Media))

		msg := tgbotapi.NewPhotoShare(int64(post.ChatID), dupImg)
		msg.Caption = "🚨 Duplicate detected! 🚨\nOriginal photo has been attached."
		msg.BaseChat.ReplyToMessageID = post.MessageID

		_, err := m.botAPI.Send(msg)
		if err != nil {
			m.log.Err(fmt.Sprintf("error while sending duplicate image report: %s", err.Error()))
		}

	}

	return dup
}

// cleanFromPosted removes all the posted entities.Post from a given array of such
// elements.
func cleanFromPosted(e []entities.Post) []entities.Post {
	t := []entities.Post{}
	for _, element := range e {
		if (time.Time{}).Equal(element.PostedAt) {
			t = append(t, element)
		}
	}

	return t
}

// GetStatus returns informations about bot's current inner status
func (m *Manager) GetStatus() (s StatusInfo) {
	var postsQueue []entities.Post
	m.db.Not("has_error", 1).Find(&postsQueue)
	postsQueue = cleanFromPosted(postsQueue)

	s.PostNumber = int64(len(postsQueue))
	s.PostPerHour = m.hourlyPostRate.String()

	return
}
