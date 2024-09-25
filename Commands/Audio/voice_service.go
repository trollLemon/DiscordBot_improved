package audio

import (


	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)



/*  Dependency Injection Implementations
 * 
 *  We abstract the DiscordGO voice connection, and put 
 *  The audio playing, and voice connection stuff behind an interface.
 *  This makes testing easier to do, and also if we decide to switch the bot library
 *  we only need to change the interfaces implementation, not the audio players
 */
type VoiceService interface {
	PlayAudioFile(url string, Done chan bool)
	Disconnect()
}


/* 
 * Voice service using DiscordGO and dgVoice
 */
type Voice struct {
	Vc *discordgo.VoiceConnection
}

func (v *Voice) PlayAudioFile(url string, Done chan bool) {
	dgvoice.PlayAudioFile(v.Vc, url, Done)
}


func (v *Voice) clean() {
	v.Vc=nil
}

func (v *Voice) Disconnect(){
	v.Vc.Disconnect()
	v.clean()
}


