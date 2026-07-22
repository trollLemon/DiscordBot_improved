package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

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
		slog.Error("failed to download attachment", "error", err)
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	imageClass, err := a.Classification.ClassifyImage(imgBytes, format)

	if err != nil {
		common.ClassificationError(s, i, "Classification failed", err)
	} else {
		common.ReplyImageClassification(imgBytes, imageClass, s, i)
	}

	return err

}
