package Commands

import (
	application "bot/Application"
	"bot/Core/Common"
	"bot/Core/Interfaces"
	"bot/util"
)

func Classify(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	imageClass, err := a.ClassificationApi.ClassifyImage(imgBytes, format)

	Common.ReplyImageClassification(imgBytes, err, imageClass, s, i)

	return err

}
