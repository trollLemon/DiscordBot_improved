package Interfaces

import "github.com/bwmarrin/discordgo"

/*

Interfaces for DiscordGo library
Idea based off this thread: https://github.com/bwmarrin/discordgo/issues/564

*/

type DiscordSession interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	InteractionResponseEdit(interaction *discordgo.Interaction, newresp *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *discordgo.VoiceConnection, err error)
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	VoiceState(guildID, userID string) (*discordgo.VoiceState, error)
}

type DiscordInteraction interface {
	ApplicationCommandData() *discordgo.ApplicationCommandInteractionData
	GetInteractionAuthor() string
	GetChannel() string
	GetInteraction() *discordgo.Interaction
}
