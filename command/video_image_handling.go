package command

// MediaType is the type of media we're dealing with
type MediaType int

//go:generate stringer -type=MediaType
const (
	Photo MediaType = iota
	Video
)

func saveMedia(fileID string, caption string, mediaType MediaType) {
	// TODO: handle video post creation
}
