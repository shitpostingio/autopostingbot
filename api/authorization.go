package api

import (
	"github.com/zelenin/go-tdlib/client"
	"github.com/shitpostingio/autopostingbot/config/structs"
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
		UseFileDatabase:        cfg.UseFileDatabase,
		UseChatInfoDatabase:    cfg.UseChatInfoDatabase,
		UseMessageDatabase:     cfg.UseMessageDatabase,
		UseSecretChats:         cfg.UseSecretChats,
		ApiId:                  cfg.APIId,
		ApiHash:                cfg.APIHash,
		SystemLanguageCode:     cfg.SystemLanguageCode,
		DeviceModel:            cfg.DeviceModel,
		SystemVersion:          cfg.SystemVersion,
		ApplicationVersion:     cfg.ApplicationVersion,
		EnableStorageOptimizer: cfg.EnableStorageOptimizer,
		IgnoreFileNames:        cfg.IgnoreFileNames,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: cfg.LogVerbosityLevel,
	})

	tdlibClient, err = client.NewClient(authorizer, logVerbosity)
	return tdlibClient, err

}
