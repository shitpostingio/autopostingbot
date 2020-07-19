package commands

import (
	"github.com/zelenin/go-tdlib/client"
)

type InfoCommandHandler struct {
}

func (InfoCommandHandler) Handle(arguments string, message *client.Message) error {

	//fileID, uniqueFileID := api.GetMediaFileIDs(message)
	//
	//
	//post := database.FindPostByFileID(fileID, repo.Db)
	//if post.ID == 0 {
	//	err = errors.New("post not found")
	//	return
	//}
	//
	//if post.PostedAt != nil {
	//	reply = fmt.Sprintf("Post added by @%s on %s\nPosted on %s\nLink: t.me/%s/%d",
	//		post.User.Handle, utility.FormatDate(post.CreatedAt), utility.FormatDate(*post.PostedAt),
	//		edition.ChannelName, post.MessageID)
	//	return
	//}
	//
	//position := database.GetQueuePositionByDatabaseID(post.ID, repo.Db)
	//timeToPost := manager.GetNextPostTime().Add(algo.EstimatePostTime(position - 1))
	//durationUntilPost := durafmt.Parse(time.Until(timeToPost).Truncate(time.Minute))
	//
	//reply = fmt.Sprintf("📋 The post is number %d in the queue\n👤 Added by @%s on %s\n\n🕜 It should be posted roughly in %s\n📅 On %s",
	//	position, post.User.Handle, utility.FormatDate(post.CreatedAt), durationUntilPost.String(), utility.FormatDate(timeToPost))
	//return

	//TODO: SISTEMARE
	return nil

}
