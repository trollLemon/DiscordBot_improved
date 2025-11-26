package Commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/common"
)

func Add(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.RandomWords.Insert(term); err != nil {
		Common.Reply(s, i, fmt.Sprintf("error inserting word: %s. %s", term, err.Error()))
		return err
	}

	Common.Reply(s, i, fmt.Sprintf("Added %s to database", term))
	return nil
}
func Remove(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.RandomWords.Delete(term); err != nil {
		Common.Reply(s, i, fmt.Sprintf("error removing word: %s. %s", term, err.Error()))
		return err
	}

	Common.Reply(s, i, fmt.Sprintf("Added %s to database", term))
	return nil
}
func Show(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	terms, err := a.RandomWords.GetAll()

	if err != nil {
		Common.Reply(s, i, "error getting all words in database")
		return err
	}
	wordsString := strings.Join(terms, "\n")

	Common.Reply(s, i, wordsString)

	return nil
}
