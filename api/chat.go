package api

import "github.com/zelenin/go-tdlib/client"

// GetChat retrieves a chat by its chatID.
// It may be required before the bot is able to send messages
// to a certain chat, if it hasn't received updates from it.
func GetChat(chatID int64) (*client.Chat, error) {
	return tdlibClient.GetChat(&client.GetChatRequest{ChatId: chatID})
}
