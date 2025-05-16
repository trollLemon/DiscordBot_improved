package audio

import (
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type VoiceService interface {
	PlayAudioFile(url string, Done chan bool)
	Disconnect()
}

/*
Voice service using DiscordGO and dgVoice
*/
type Voice struct {
	Vc *discordgo.VoiceConnection
}

func (v *Voice) PlayAudioFile(url string, Done chan bool) {
	dgvoice.PlayAudioFile(v.Vc, url, Done)
}

func (v *Voice) clean() {
	v.Vc = nil
}

func (v *Voice) Disconnect() {
	v.Vc.Disconnect()
	v.clean()
}
