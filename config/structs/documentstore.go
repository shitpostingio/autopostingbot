package structs

import "go.mongodb.org/mongo-driver/mongo/options"

// DocumentStoreConfiguration represents a document store configuration
type DocumentStoreConfiguration struct {
	UseAuthentication bool `type:"optional"`
	UseReplicaSet     bool `type:"optional"`
	DatabaseName      string
	AuthMechanism     string   `type:"optional"`
	Username          string   `type:"optional"`
	Password          string   `type:"optional"`
	AuthSource        string   `type:"optional"`
	ReplicaSetName    string   `type:"optional"`
	Hosts             []string `type:"optional"`
}

// MongoDBConnectionOptions gets the connection options from the DocumentStoreConfiguration
func (c *DocumentStoreConfiguration) MongoDBConnectionOptions() *options.ClientOptions {

	//TODO: CHECK
	//
	clientOptions := options.Client()
	clientOptions.SetHosts(c.Hosts)

	//
	if c.UseAuthentication {
		clientOptions.SetAuth(options.Credential{
			AuthMechanism: c.AuthMechanism,
			AuthSource:    c.AuthSource,
			Username:      c.Username,
			Password:      c.Password,
			PasswordSet:   true,
		})
	}

	//
	if c.UseReplicaSet {
		clientOptions.SetReplicaSet(c.ReplicaSetName)
	}

	return clientOptions

}
