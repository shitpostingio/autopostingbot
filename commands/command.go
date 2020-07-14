package commands

import "github.com/zelenin/go-tdlib/client"

type Handler interface {
	Handle(message *client.Message) error
}
