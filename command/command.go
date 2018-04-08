package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/shitposting/autoposting-bot/algo"
)

// Handle e` il punto di entrata per il parsing e l'organizzazione dell'azione del bot
// su un messaggio entrante.
func Handle(update tgbotapi.Update, api *tgbotapi.BotAPI, manager *algo.Manager) error {
	if update.Message == nil && update.EditedMessage == nil {
		return errors.New("update Message or EditedMessage body was nil, most likely an error on Telegram side")
	}

	msg := update.Message
	editedMsg := update.EditedMessage

	if editedMsg != nil {
		switch {
		case editedMsg.Video != nil:
			modifyMedia(editedMsg.Video.FileID, editedMsg.Caption, manager, editedMsg.From.ID, editedMsg.MessageID, int(editedMsg.Chat.ID))
		case editedMsg.Photo != nil:
			photos := *editedMsg.Photo
			modifyMedia(photos[len(photos)-1].FileID, editedMsg.Caption, manager, editedMsg.From.ID, editedMsg.MessageID, int(editedMsg.Chat.ID))
		}

		return nil
	}

	switch {
	case msg.Video != nil:
		saveMedia(msg.Video.FileID, msg.Caption, Video, manager, msg.From.ID, msg.MessageID, int(msg.Chat.ID))
	case msg.Photo != nil:
		photos := *msg.Photo
		saveMedia(photos[len(photos)-1].FileID, msg.Caption, Image, manager, msg.From.ID, msg.MessageID, int(msg.Chat.ID))
	case msg.Text != "":
		arg := update.Message.CommandArguments()

		fmt.Println(arg)

		if path, err := getVideo(arg); err == nil {
			fmt.Println("path is " + path)
		}
	}

	return nil
}

func getVideo(url string) (path string, err error) {
	var stdoutStderr []byte
	cmd := exec.Command("youtube-dl", "-f", "mp4", "-o", "%(id)s.%(ext)s", url)
	stdoutStderr, err = cmd.CombinedOutput()

	if err != nil {
		return
	}

	fmt.Printf("%s\n", stdoutStderr)
	temp := string(stdoutStderr)
	log := strings.Split(temp, "\n")

	for _, element := range log {
		if strings.HasPrefix(element, "ERROR: ") {
			temp := strings.Split(element, ". ")
			err = errors.New(temp[1])
		} else if strings.HasPrefix(element, "[download] Destination: ") {
			temp := strings.Trim(element, "[download] Destination: ")
			path = temp
		}
	}
	return
}
