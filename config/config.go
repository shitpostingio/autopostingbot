package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the bot configuration representation, read
// from a configuration file.
type Config struct {
	BotToken         string
	Port             int
	ReverseProxy     bool
	ReverseProxyPort int
	IP               string
	Domain           string
	TLS              bool
	TLSCertPath      string
	TLSKeyPath       string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseAddress  string
	ChannelID        int
}

// DatabaseConnectionString returns a well-formatted database connection string for MySQL
func (c Config) DatabaseConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=True", c.DatabaseUsername, c.DatabasePassword, c.DatabaseAddress, c.DatabaseName)
}

// BindString returns IP+Port, in a suitable syntax for http.ListenAndServe
func (c Config) BindString() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

// WebHookURL returns the URL to listen on for WebHooks.
// It is useful since we have to tell Telegram where we want our updates.
func (c Config) WebHookURL() string {
	port := strconv.Itoa(c.Port)

	if c.ReverseProxy {
		port = strconv.Itoa(c.ReverseProxyPort)
	}
	return "https://" + c.Domain + ":" + port + "/" + c.BotToken + "/updates"
}

// WebHookPath returns only the relative path where Telegram will send
// updates
func (c Config) WebHookPath() string {
	return "/" + c.BotToken + "/updates"
}

// ReadConfigFile reads a configuration file and returns its Config instance
func ReadConfigFile(path string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return Config{}, err
	}

	if conf.BotToken == "" { // Missing bot token
		return buildErrorMessage("missing Bot token")
	} else if conf.DatabaseUsername == "" {
		return buildErrorMessage("missing database username")
	} else if conf.DatabasePassword == "" {
		return buildErrorMessage("missing database password")
	} else if conf.DatabaseName == "" {
		return buildErrorMessage("missing database name")
	} else if conf.ChannelID == 0 {
		return buildErrorMessage("missing Telegram channel identifier")
	} else if !conf.ReverseProxy && !isStandardPort(conf.Port) { // Not running behind a reverse proxy, and using non-standard port
		return buildErrorMessage("cannot use non-standard port when ReverseProxy is disabled")
	} else if conf.ReverseProxy && conf.ReverseProxyPort == 0 { // Running behind a reverse proxy, but its port has not been defined
		return buildErrorMessage("running behind a reverse proxy, but no reverse proxy port has been defined")
	} else if conf.Domain == "" { // Domain not set
		return buildErrorMessage("Domain not set")
	} else if strings.HasPrefix(conf.Domain, "http://") || strings.HasPrefix(conf.Domain, "https://") {
		return buildErrorMessage("Domain must not contain http:// or https://")
	} else if conf.TLS { // If we must use TLS...
		if conf.TLSCertPath == "" { // ...check if we have a TLS certificate path...
			return buildErrorMessage("missing TLS certificate path")
		} else if conf.TLSKeyPath == "" { // ...and a TLS key path
			return buildErrorMessage("missing TLS key path")
		}
	}

	// If we run behind a reverse proxy, bind to localhost
	if conf.ReverseProxy {
		conf.IP = "127.0.0.1"
	}

	// If we don't have a DatabaseAddress, set it to localhost
	if conf.DatabaseAddress == "" {
		conf.DatabaseAddress = "127.0.0.1:3306"
	}

	return conf, nil
}

func buildErrorMessage(message string) (Config, error) {
	return Config{}, errors.New(message)
}

func isStandardPort(port int) bool {
	switch port {
	case 443, 80, 88, 8443:
		return true
	default:
		return false
	}
}
