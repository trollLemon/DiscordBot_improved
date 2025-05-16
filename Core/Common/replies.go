package Common

import (
	"bot/Core/Interfaces"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Reply(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, text string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	}

	if err := s.InteractionRespond(i.GetInteraction(), response); err != nil {
		log.Error().Err(err).Msg("Interaction Response")
	}

}
