package Wrappers

import "github.com/bwmarrin/discordgo"

type DiscordSessionWrapper struct {
	*discordgo.Session
}

type InteractionCreateWrapper struct {
	*discordgo.InteractionCreate
}

func (i InteractionCreateWrapper) GetInteractionAuthor() string {
	return i.Member.User.ID
}

func (i InteractionCreateWrapper) GetChannel() string {
	return i.ChannelID
}

func (i InteractionCreateWrapper) GetInteraction() *discordgo.Interaction {
	return i.InteractionCreate.Interaction
}

func (i InteractionCreateWrapper) GetImageURLFromAttachmentID(id string) string {
	return i.ApplicationCommandData().Resolved.Attachments[id].URL
}
