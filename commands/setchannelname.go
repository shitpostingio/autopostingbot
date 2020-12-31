package commands

import (
	"errors"
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/config"
	l "github.com/shitpostingio/autopostingbot/localization"
	"github.com/zelenin/go-tdlib/client"
	"strings"
)

// SetChannelNameCommandHandler represents the handler of the /setchannelname command.
type SetChannelNameCommandHandler struct{}

// Handle handles the /setchannelname command.
// /setchannelname sets the channel name in the configuration.
func (SetChannelNameCommandHandler) Handle(arguments string, message, _ *client.Message) error {

	if len(arguments) == 0 {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_CHANNEL_HANDLE_MISSING))
		return errors.New("channel handle missing")
	}

	channelName := strings.ReplaceAll(arguments, "@", "")
	if len(channelName) < 4 {
		_, _ = api.SendPlainReplyText(message.ChatId, message.Id, l.GetString(l.COMMANDS_CHANNEL_HANDLE_TOO_SHORT))
		return errors.New("channel handle must be at least 4 characters long")
	}

	confirmation := fmt.Sprintf(l.GetString(l.COMMANDS_CHANNEL_HANDLE_UPDATED), channelName)
	_, _ = api.SendPlainReplyText(message.ChatId, message.Id, confirmation)
	return config.UpdateChannelHandle(channelName)

}
