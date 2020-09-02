package localization

import (
	"github.com/bykovme/gotrans"
)

var (
	language string
)

// SetLanguage sets the language for the bot.
func SetLanguage(toSet string) {
	language = toSet
}

// GetString gets the translation given a key.
func GetString(key string) string {
	return gotrans.Tr(language, key)
}
