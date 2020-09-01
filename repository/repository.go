package repository

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	Config *structs.Config
	Tdlib  *client.Client
	Me     *client.User
)
