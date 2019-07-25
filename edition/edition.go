package edition

const (

	//Shitpost is the enum value for the shitpost version
	Shitpost = iota

	//Sushiporn is the enum value for the sushiporn version
	Sushiporn
)

var (

	// ChannelName is the name of the channel
	ChannelName string

	// edition is the enum value of the active version
	edition uint
)

// SetEdition sets the active edition and the
// relative variables
func SetEdition(activeEdition uint) {

	switch activeEdition {
	case Shitpost:
		ChannelName = "shitpost"
	case Sushiporn:
		ChannelName = "sushiporn"
	}

	edition = activeEdition
}

// GetEditionString returns a print friendly version of the edition
func GetEditionString() string {
	switch {
	case IsShitpost():
		return "Shitpost edition"
	case IsSushiporn():
		return "Sushiporn edition"
	default:
		return "ERROR IN EDITION"
	}
}

// IsShitpost returns true if the active edition
// is the shitpost one
func IsShitpost() bool {
	return edition == Shitpost
}

// IsSushiporn returns true if the active edition
// is the sushiporn one
func IsSushiporn() bool {
	return edition == Sushiporn
}
