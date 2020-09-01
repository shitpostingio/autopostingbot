package structs

import "go.mongodb.org/mongo-driver/mongo/options"

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

	//TODO: SISTEMARE
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
