package structs

// TdlibConfiguration represents the Tdlib configuration.
type TdlibConfiguration struct {

	// LogVerbosityLevel controls log verbosity.
	LogVerbosityLevel int32 `type:"optional"`

	// APIId is the application identifier for Telegram API access, which can be obtained at https://my.telegram.org
	APIId int32

	// APIHash is the application identifier hash for Telegram API access, which can be obtained at https://my.telegram.org
	APIHash string

	// DatabaseDirectory is the path to the directory for the persistent database.
	// If empty, the current working directory will be used.
	DatabaseDirectory string

	// FilesDirectory is the path to the directory for storing files.
	// If empty, DatabaseDirectory will be used.
	FilesDirectory string

	// SystemLanguageCode is the IETF language tag of the user's operating system language.
	SystemLanguageCode string `type:"optional"`

	// DeviceModel is the model of the device the application is being run on.
	DeviceModel string `type:"optional"`

	// SystemVersion is the version of the operating system the application is being run on.
	SystemVersion string `type:"optional"`

	// ApplicationVersion is the application version.
	ApplicationVersion string `type:"optional"`

	// UseFileDatabase, if set to true will save information about downloaded and uploaded files
	// between application restarts.
	UseFileDatabase bool `type:"optional"`

	// UseChatInfoDatabase, if set to true will maintain a cache of users, basic groups,
	// supergroups, channels and secret chats. It implies UseFileDatabase.
	UseChatInfoDatabase bool `type:"optional"`

	// UseMessageDatabase, if set to true will maintain a cache of chats and messages.
	// It implies UseChatInfoDatabase.
	UseMessageDatabase bool `type:"optional"`

	// UseSecretChats, if set to true enables support for secret chats.
	UseSecretChats bool `type:"optional"`

	// UseTestDc, if set to true, will use the Telegram test environment instead of the production one.
	UseTestDc bool `type:"optional"`

	// EnableStorageOptimizer, if set to true will automatically delete old files.
	EnableStorageOptimizer bool `type:"optional"`

	// IgnoreFileNames, if set to true will ignore original file names.
	IgnoreFileNames bool `type:"optional"`
}
