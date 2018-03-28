package utility

import (
	"log"

	"github.com/fatih/color"
)

// PrettyError prints an error with a pleasant red text color
func PrettyError(err error) {
	log.Println(color.RedString("ERROR - %s", err.Error()))
}

// PrettyFatal prints an error with a pleasant red text color, then exits
func PrettyFatal(err error) {
	log.Fatal(color.RedString("ERROR - %s", err.Error()))
}
