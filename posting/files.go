package posting

import (
	"fmt"
	"github.com/shitpostingio/autopostingbot/api"
	"github.com/shitpostingio/autopostingbot/documentstore/entities"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// moveToDirectory moves a file from the Tdlib directory to the
// persistent directory specified in the configuration file.
func moveToDirectory(post *entities.Post) error {

	//
	file, err := api.DownloadFile(post.Media.TdlibID)
	if err != nil {
		return err
	}

	pieces := strings.Split(file.Local.Path, "/")
	log.Debugln("Pieces: ", pieces)
	fileName := pieces[len(pieces)-1]
	log.Debugln("Filename: ", fileName)
	extension := fileName[strings.LastIndex(fileName, ".")+1:]
	log.Debugln("Extension: ", extension)

	err = os.Rename(file.Local.Path, fmt.Sprintf("%s/%s.%s", m.config.Autoposting.MediaPath, post.Media.FileID, extension))
	return err

}
