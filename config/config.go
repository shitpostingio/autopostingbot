package configuration

import (
	"log"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

//TODO: ASSOLUTAMENTE SISTEMARE STO PORCAIO GALATTICO
const (
	defaultFileSizeThreshold  = 20971520 //20MB
	defaultPostAlertThreshold = 10
	defaultSocketPath         = "/tmp/loglog.socket"
	defaultDBAddress          = "localhost:3306"
)

// Load reads a configuration file and returns its config instance
func Load(path string, useWebhook bool) (cfg Config, err error) {

	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	checkMandatoryFields(cfg)
	checkOptionalFields(&cfg)

	if useWebhook {
		checkWebhookConfig(&cfg.Server)
	}

	return
}

// checkMandatoryFields uses reflection to see if there are
// mandatory fields with zero value
func checkMandatoryFields(config Config) {
	checkStruct(reflect.TypeOf(config), reflect.ValueOf(config))
}

// checkOptionalFields sets default values in case optional fields
// have a zero value
func checkOptionalFields(config *Config) {

	if config.FileSizeThreshold == 0 {
		config.FileSizeThreshold = defaultFileSizeThreshold //20MB
	}

	if config.PostAlertThreshold == 0 {
		config.PostAlertThreshold = defaultPostAlertThreshold
	}

	if config.LogLog.SocketPath == "" {
		config.LogLog.SocketPath = defaultSocketPath
	}

	if config.DB.Address == "" {
		config.DB.Address = defaultDBAddress
	}
}

// checkWebhookConfig perform checks and sets default values for
// webhook configuration details
func checkWebhookConfig(config *ServerDetails) {

	if config.ReverseProxy {

		/* LOCALHOST FOR REVERSE PROXY */
		config.IP = "127.0.0.1"

		if !isStandardPort(config.Port) {
			log.Fatal("cannot use non-standard port when ReverseProxy is disabled")
		}
	}

	if config.Domain == "" { // Domain not set
		log.Fatal("Domain not set")
	} else if strings.HasPrefix(config.Domain, "http://") || strings.HasPrefix(config.Domain, "https://") {
		log.Fatal("Domain must not contain http:// or https://")
	}

	if config.TLS {
		if config.TLSCertPath == "" {
			log.Fatal("missing TLS certificate path")
		} else if config.TLSKeyPath == "" {
			log.Fatal("missing TLS key path")
		}
	}
}

// checkStruct explores structures recursively and checks if
// struct fields have a zero value
func checkStruct(typeToCheck reflect.Type, valueToCheck reflect.Value) {

	for i := 0; i < typeToCheck.NumField(); i++ {

		currentField := typeToCheck.Field(i)
		currentValue := valueToCheck.Field(i)

		if currentField.Type.Kind() == reflect.Struct {
			checkStruct(currentField.Type, currentValue)
		} else {
			checkField(currentField, currentValue)
		}
	}
}

// checkField checks if a field is optional or a webhook field
// if it isn't, it checks if the field has a zero value
func checkField(typeToCheck reflect.StructField, valueToCheck reflect.Value) {

	typeTagValue := typeToCheck.Tag.Get("type")

	if typeTagValue == "optional" || typeTagValue == "webhook" {
		return
	}

	zeroValue := reflect.Zero(typeToCheck.Type)

	if valueToCheck.Interface() == zeroValue.Interface() {
		log.Fatalf("non optional field %s had zero value", typeToCheck.Name)
	}
}
