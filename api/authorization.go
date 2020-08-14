package api

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
)

var (
	tdlibClient *client.Client
)

// Authorize logs the bot into the provided account using tdlib.
func Authorize(botToken string, cfg *structs.TdlibConfiguration) (tClient *client.Client, err error) {

	authorizer := client.BotAuthorizer(botToken)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              cfg.UseTestDc,
		DatabaseDirectory:      cfg.DatabaseDirectory,
		FilesDirectory:         cfg.FilesDirectory,
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  cfg.APIID,
		ApiHash:                cfg.APIHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Other",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        true,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: cfg.LogVerbosityLevel,
	})

	tdlibClient, err = client.NewClient(authorizer, logVerbosity)
	return tdlibClient, err

}
