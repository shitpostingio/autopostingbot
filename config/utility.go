package configuration

import (
	"fmt"
	"strconv"
)

// DatabaseConnectionString returns a connection string for MariaDB
func (c DBConfig) DatabaseConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=True", c.Username, c.Password, c.Address, c.Name)
}

// BindString returns IP+Port, in a suitable syntax for http.ListenAndServe
func (c ServerDetails) BindString() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

// WebHookURL returns the URL to listen on for WebHooks
func (c Config) WebHookURL() string {
	port := strconv.Itoa(c.Server.Port)

	if c.Server.ReverseProxy {
		port = strconv.Itoa(c.Server.ReverseProxyPort)
	}
	return "https://" + c.Server.Domain + ":" + port + "/" + c.BotToken + "/updates"
}

// WebHookPath returns only the relative path where Telegram will send updates
func (c Config) WebHookPath() string {
	return "/" + c.BotToken + "/updates"
}

// isStandardPort returns if the specified port is suitable
// for a webhook connection without reverse proxy
func isStandardPort(port int) bool {
	switch port {
	case 443, 80, 88, 8443:
		return true
	default:
		return false
	}
}
