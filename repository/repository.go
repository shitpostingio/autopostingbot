package repository

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	// Config contains all the configuration structures.
	Config *structs.Config

	// Tdlib is the Telegram client instance.
	Tdlib *client.Client

	// Me represents the current bot as a Telegram client.User.
	Me *client.User
)
