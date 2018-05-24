package loglogclient

import (
	"encoding/gob"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
)

// LoglogClient is a logging facility, that works on both stdout and
// by logging to a loglog server
type LoglogClient struct {
	config Config
}

// Config define how a LoglogClient behaves
type Config struct {
	SocketPath    string
	ApplicationID string
}

// NewClient returns a LoglogClient with the LoglogClientConfig passed as parameter
func NewClient(llcc Config) *LoglogClient {
	llc := LoglogClient{
		config: llcc,
	}
	return &llc
}

// Log logs a Payload to the database
func (llc LoglogClient) Log(p Payload) error {
	c, err := net.Dial("unix", llc.config.SocketPath)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(c)
	err = encoder.Encode(&p)
	_ = c.Close()

	return err
}

// Info logs something on both stdout and the internal database.
// On stdout, there will be a green string formatted as "TIME INFO: log"
func (llc *LoglogClient) Info(s string) {
	log.Println(color.GreenString("INFO: %s", s))
	err := llc.Log(Payload{
		ApplicationID: llc.config.ApplicationID,
		Type:          INFO,
		Payload:       s,
		CreatedAt:     time.Now(),
	})

	llc.serverNotOn(err)
}

// Warn logs something on both stdout and the internal database.
// On stdout, there will be a yellow string formatted as "TIME WARN: log"
func (llc *LoglogClient) Warn(s string) {
	log.Println(color.YellowString("WARN: %s", s))
	err := llc.Log(Payload{
		ApplicationID: llc.config.ApplicationID,
		Type:          WARNING,
		Payload:       s,
		CreatedAt:     time.Now(),
	})

	llc.serverNotOn(err)
}

// Err logs something on both stdout and the internal database.
// On stdout, there will be a red string formatted as "TIME ERR: log"
func (llc *LoglogClient) Err(s string) {
	log.Println(color.RedString("ERROR: %s", s))
	err := llc.Log(Payload{
		ApplicationID: llc.config.ApplicationID,
		Type:          ERROR,
		Payload:       s,
		CreatedAt:     time.Now(),
	})

	llc.serverNotOn(err)
}

func (llc *LoglogClient) serverNotOn(err error) {
	if err != nil {
		log.Println("SERVER: it seems that the loglog server has not been powered on, future logs will not be retained but I'll still keep on printing colored stdout logging")
	}
}
