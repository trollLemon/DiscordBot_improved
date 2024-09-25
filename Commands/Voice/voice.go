package voice

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)



func JoinVoiceChannel(s *discordgo.Session, author string, guild string) (*discordgo.VoiceConnection, error) {
	
	voiceChannel, err := s.State.VoiceState(guild,author)
	
	if err != nil {
		return nil, fmt.Errorf("User is not in a VC, join a VC to use the command")
	}
	
	dgv,err := s.ChannelVoiceJoin(guild,voiceChannel.ChannelID,false,false)
	
	return dgv,err
}


