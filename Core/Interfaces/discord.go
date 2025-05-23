package Interfaces

import "github.com/bwmarrin/discordgo"

/*

Interfaces for DiscordGo library
Idea based off this thread: https://github.com/bwmarrin/discordgo/issues/564

*/

type DiscordSession interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	InteractionResponseEdit(interaction *discordgo.Interaction, newresp *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

type DiscordInteraction interface {
	ApplicationCommandData() discordgo.ApplicationCommandInteractionData
	GetInteractionAuthor() string
	GetChannel() string
	GetInteraction() *discordgo.Interaction
	GetImageURLFromAttachmentID(id string) string
}
