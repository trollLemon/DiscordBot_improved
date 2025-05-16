/* interface for a notification service
 * This service is a dependency for the audioplayer object, allowing it to notify the user when errors occur during
 * the audio stream, or simpler things like the queue being empty, or what it starts to play
 *
 */

package audio

import (
	"bot/Core/Interfaces"
)

type NotificationService interface {
	SendNotification(content string)
	SendError(error string)
}

type Notifier struct {
	Session Interfaces.DiscordSession
	Channel string
}

func (n *Notifier) SendNotification(content string) {

	n.Session.ChannelMessageSend(n.Channel, content)

}
func (n *Notifier) SendError(content string) {

	n.Session.ChannelMessageSend(n.Channel, content)

}
