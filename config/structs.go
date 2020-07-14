package configuration

import "go.mongodb.org/mongo-driver/mongo/options"

// Config is a structure containing information used
// to set up the bot
type Config struct {
	BotToken           string
	ChannelID          int64
	FileSizeThreshold  int `type:"optional"`
	MemePath           string
	PostAlertThreshold int `type:"optional"`
	LogLog             LoglogConfig
	Server             ServerDetails
	DB                 DBConfig
	Fpserver           FpServerConfig
	Tdlib              TdlibConfiguration
	DocumentStore      DocumentStoreConfiguration
}

// LoglogConfig contains the configuration for Loglog
type LoglogConfig struct {
	SocketPath    string `type:"optional"`
	ApplicationID string
}

//FpServerConfig contains the configuration for FPServer
type FpServerConfig struct {
	Address                  string
	ImageEndpoint            string
	VideoEndpoint            string
	AuthorizationHeaderName  string
	AuthorizationHeaderValue string
	CallerAPIKeyHeaderName   string
	FilePathHeaderName       string
}

//DBConfig contains the configuration for a database
type DBConfig struct {
	Name     string
	Username string
	Password string
	Address  string `type:"optional"`
}

//ServerDetails contains the details for webhook updates
type ServerDetails struct {
	Port             int    `type:"webhook"`
	ReverseProxy     bool   `type:"webhook"`
	ReverseProxyPort int    `type:"webhook"`
	IP               string `type:"webhook"`
	Domain           string `type:"webhook"`
	TLS              bool   `type:"webhook"`
	TLSCertPath      string `type:"webhook"`
	TLSKeyPath       string `type:"webhook"`
}

// DocumentStoreConfiguration represents a document store configuration
type DocumentStoreConfiguration struct {
	DatabaseName   string
	Username       string
	Password       string
	AuthSource     string
	CollectionName string
	ReplicaSetName string
	Hosts          []string `type:"optional"` //TODO: rimuovere
}

// MongoDBConnectionOptions gets the connection options from the DocumentStoreConfiguration
func (c *DocumentStoreConfiguration) MongoDBConnectionOptions() *options.ClientOptions {

	clientOptions := options.Client()
	//clientOptions.SetAuth(options.Credential{
	//	AuthMechanism: "SCRAM-SHA-1",
	//	AuthSource:    c.AuthSource,
	//	Username:      c.Username,
	//	Password:      c.Password,
	//	PasswordSet:   true,
	//})

	clientOptions.SetHosts(c.Hosts)
	//clientOptions.SetReplicaSet(c.ReplicaSetName) //TODO: SISTEMARE
	return clientOptions

}
