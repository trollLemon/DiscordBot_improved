package factories

import (
	"bot/Core/Interfaces"
	audio "bot/Core/Services/Audio"
	"github.com/bwmarrin/discordgo"
)

type ServiceFactory interface {
	CreateVoiceService(dgv *discordgo.VoiceConnection) audio.VoiceService
	CreateNotificationService(s Interfaces.DiscordSession, channel string) audio.NotificationService
}

type DynamicServiceFactory struct{}

func (*DynamicServiceFactory) CreateVoiceService(dgv *discordgo.VoiceConnection) audio.VoiceService {
	return &audio.Voice{
		Vc: dgv,
	}
}

func (*DynamicServiceFactory) CreateNotificationService(s Interfaces.DiscordSession, channel string) audio.NotificationService {
	return &audio.Notifier{
		Session: s,
		Channel: channel,
	}
}
