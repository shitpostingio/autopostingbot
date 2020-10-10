package repository

import (
	"github.com/shitpostingio/autopostingbot/config/structs"
	"github.com/zelenin/go-tdlib/client"
)

var (
	// Config contains all the configuration structures.
	Config *structs.Config

	// Tdlib is the Telegram client instance.
	Tdlib *client.Client

	// Me represents the current bot as a Telegram client.User.
	Me *client.User
)
