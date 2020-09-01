package structs

// TdlibConfiguration represents the tdlib configuration
type TdlibConfiguration struct {
	UseTestDc              bool `type:"optional"`
	DatabaseDirectory      string
	FilesDirectory         string
	UseFileDatabase        bool `type:"optional"`
	UseChatInfoDatabase    bool `type:"optional"`
	UseMessageDatabase     bool `type:"optional"`
	UseSecretChats         bool `type:"optional"`
	APIId                  int32
	APIHash                string
	SystemLanguageCode     string `type:"optional"`
	DeviceModel            string `type:"optional"`
	SystemVersion          string `type:"optional"`
	ApplicationVersion     string `type:"optional"`
	EnableStorageOptimizer bool   `type:"optional"`
	IgnoreFileNames        bool   `type:"optional"`
	LogVerbosityLevel      int32  `type:"optional"`
}
