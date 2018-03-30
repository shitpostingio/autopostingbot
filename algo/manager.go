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

	"github.com/go-telegram-bot-api/telegram-bot-api"
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
	hourlyPostSignal   <-chan time.Time
	hourlyPostRate     time.Duration
	postSignal         <-chan time.Time
	debug              bool
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
func NewManager(mc ManagerConfig) (m Manager, err error) {
	if mc.ChannelID == 0 {
		err = errors.New("ChannelID is empty")
		return
	}

	m = Manager{
		botAPI:             mc.BotAPIInstance,
		channelID:          mc.ChannelID,
		AddImageChannel:    make(chan MediaPayload),
		AddVideoChannel:    make(chan MediaPayload),
		ModifyMediaChannel: make(chan MediaPayload),
		debug:              mc.Debug,
	}

	// Initialize gorm
	m.db, err = gorm.Open("mysql", mc.DatabaseString)
	if err != nil {
		return
	}

	// Get and initialize the categories
	m.db.Where("name = ?", "image").First(&imageCategory)
	m.db.Where("name = ?", "video").First(&videoCategory)
	if imageCategory.Name != "image" || videoCategory.Name != "video" {
		err = errors.New("cannot load video and/or image categories identities from the database")
		return
	}

	// Calculate the hourly post rate on the current post availability
	m.calculateHourlyPostRate()

	// Print the hourly posting rate in minutes
	utility.YellowLog("Initial hourly posting rate set to " + m.hourlyPostRate.String())

	// Initialize the calculation signal
	m.hourlyPostSignal = time.After(1 * time.Hour)

	// Initialize the postSignal on the hourlyRate
	m.setUpPostSignal()

	// Start the manager lifecycle
	go m.managerLifecycle()

	return
}

// managerLifecycle is the function we run indefinitely in a goroutine.
// It handles incoming updates, and the posting routine.
func (m *Manager) managerLifecycle() {

	// if -debug is specified, immediately send a post and exit
	if m.debug {
		utility.GreenLog("it's time to post!")
		wtp, err := m.whatToPost()
		if err != nil {
			utility.PrettyError(err)

		}

		if err := m.popAndPost(wtp); err != nil {
			utility.PrettyError(err)
			utility.PrettyError(fmt.Errorf("on media with ID %s", wtp.Media))
		} else {
			fmt.Println(wtp)
			utility.GreenLog("all done!")
		}

		os.Exit(0)
	}
	m.setUpPostSignal()

	for {
		select {
		case modifiedPost := <-m.ModifyMediaChannel:
			var entity entities.Post
			id, err := getUserID(m.db, int(modifiedPost.Entity.UserID))
			if err != nil {
				utility.PrettyError(err)
			}
			m.db.Where("media = ? AND user_id = ?", modifiedPost.Entity.Media, id).First(&entity)

			if entity.Media == "" { // an empty media ID means no entity with said ID was found
				utility.PrettyError(fmt.Errorf("someone tried to update the caption for media with id %s, but i don't know any", modifiedPost.Entity.Media))
				continue
			}

			entity.Caption = modifiedPost.Entity.Caption
			m.db.Save(&entity)
			utility.SendTelegramReply(modifiedPost.ChatID, modifiedPost.MessageID, m.botAPI, "Modified!")
		case newPost := <-m.AddVideoChannel:
			utility.GreenLog("got a new video to add!")

			// if we have a duplicate, write a log message and stop
			if m.checkDuplicate(newPost) {
				continue
			}

			userID, err := getUserID(m.db, int(newPost.Entity.UserID))
			if err != nil {
				utility.PrettyError(err)
			}

			newPost.Entity.UserID = uint(userID)
			newPost.Entity.Categories = []entities.Category{videoCategory}

			// add to the database
			m.db.Create(&newPost.Entity)
			utility.SendTelegramReply(newPost.ChatID, newPost.MessageID, m.botAPI, "Video added!")
		case newPost := <-m.AddImageChannel:
			utility.GreenLog("got a new image to add!")

			// if we have a duplicate, write a log message and stop
			if m.checkDuplicate(newPost) {
				continue
			}

			userID, err := getUserID(m.db, int(newPost.Entity.UserID))
			if err != nil {
				utility.PrettyFatal(err)
			}

			newPost.Entity.UserID = uint(userID)
			newPost.Entity.Categories = []entities.Category{imageCategory}

			// add to the database
			m.db.Create(&newPost.Entity)
			utility.SendTelegramReply(newPost.ChatID, newPost.MessageID, m.botAPI, "Image added!")
		case <-m.postSignal:
			// setup the post signal first
			m.setUpPostSignal()

			utility.GreenLog("it's time to post!")

			// could not find anything to post
			wtp, err := m.whatToPost()
			if err != nil {
				utility.PrettyError(err)
				continue
			}

			if err := m.popAndPost(wtp); err != nil {
				// posting did not go well...
				utility.PrettyError(err)
				utility.PrettyError(fmt.Errorf("on media with ID %s", wtp.Media))
				continue
			}

			utility.GreenLog("all done!")
		case <-m.hourlyPostSignal:
			utility.YellowLog("calculating the hourly posting rate...")
			// calculate the new hourly post rate
			m.calculateHourlyPostRate()

			// set up the posting signal, even if we already did that before
			m.setUpPostSignal()

			utility.YellowLog(fmt.Sprintf("new hourly posting rate: %s", m.hourlyPostRate.String()))

			// see you in an hour!
			m.hourlyPostSignal = time.After(1 * time.Hour)
		}
	}
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
	m.db.Find(&postsQueue)

	ppH := postsPerHour(postsQueue)

	if m.debug {
		utility.BlueLog(fmt.Sprintf("posts per hour: %d", ppH))
	}

	if ppH > 0 {
		if 60/ppH <= 0 {
			m.hourlyPostRate = time.Duration(1) * time.Minute
		} else {
			m.hourlyPostRate = time.Duration(60/ppH) * time.Minute
		}

		if m.debug {
			utility.BlueLog(fmt.Sprintf("hourly post rate: %d", m.hourlyPostRate))
		}
		return
	}
	m.hourlyPostRate = 0
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
	m.db.Preload("Categories").Find(&postsQueue)
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
	switch {
	case entity.IsImage(m.db):
		tgImage := tgbotapi.NewPhotoShare(m.channelID, entity.Media)
		tgImage.Caption = caption
		_, err = m.botAPI.Send(tgImage)
	case entity.IsVideo(m.db):
		tgVideo := tgbotapi.NewVideoShare(m.channelID, entity.Media)
		tgVideo.Caption = caption
		_, err = m.botAPI.Send(tgVideo)
	}

	// checking if there's an error here gives us the chance to remove the posted
	// entity if everything was ok - if it wasn't, the error will be handled on the caller
	// level.
	if err == nil {
		m.db.Model(&entity).Association("categories").Delete(entity.Categories)
		m.db.Delete(&entity)
	}
	return err
}

// isDuplicate returns true if post has been already added before
// false otherwise.
func (m Manager) isDuplicate(post entities.Post) bool {
	var duplicate entities.Post
	m.db.Where("media = ?", post.Media).First(&duplicate)

	if duplicate.Media != "" {
		return true
	}

	return false
}

// checkDuplicate checks whether post has already been added to the database,
// and if yes, it will communicate it to the user
func (m Manager) checkDuplicate(post MediaPayload) bool {
	dup := m.isDuplicate(post.Entity)

	if dup {
		msg := fmt.Sprintf("user %d tried to re-add media %s, which is already present in the database", post.Entity.UserID, post.Entity.Media)
		utility.SendTelegramReply(post.ChatID, post.MessageID, m.botAPI, "	Duplicate.")
		utility.YellowLog(msg)
	}

	return dup
}
