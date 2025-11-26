package Commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/common"
	"github.com/trollLemon/DiscordBot/internal/util"
)

func Classify(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		log.Err(err).Msg("failed to download attachment")
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	imageClass, err := a.Classification.ClassifyImage(imgBytes, format)

	if err != nil {
		Common.ClassificationError(s, i, "Classification failed", err.Error())
	} else {
		Common.ReplyImageClassification(imgBytes, imageClass, s, i)
	}

	return err

}
