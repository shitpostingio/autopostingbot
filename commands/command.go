package commands

import "github.com/zelenin/go-tdlib/client"

// Handler is the base interface for all command handlers.
type Handler interface {

	// Handle is the method that command handlers implement.
	Handle(arguments string, message, replyToMessage *client.Message) error
}
