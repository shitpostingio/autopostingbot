package commands

import "github.com/zelenin/go-tdlib/client"

type Handler interface {
	Handle(arguments string, message *client.Message) error
}
