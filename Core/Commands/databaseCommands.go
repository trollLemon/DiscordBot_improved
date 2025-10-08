package Commands

import (
	"bot/Application"
	"bot/Core/Common"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)
func Add(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.WordDatabase.Insert(term); err != nil {
		Common.Reply(s, i, fmt.Sprintf("error inserting word: %s. %s", term, err.Error()))
		return err
	}

	Common.Reply(s, i, fmt.Sprintf("Added %s to database", term))
	return nil
}
func Remove(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.WordDatabase.Delete(term); err != nil {
		Common.Reply(s, i, fmt.Sprintf("error removing word: %s. %s", term, err.Error()))
		return err
	}

	Common.Reply(s, i, fmt.Sprintf("Added %s to database", term))
	return nil
}
func Show(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	terms, err := a.WordDatabase.GetAll()

	if err != nil {
		Common.Reply(s, i, "error getting all words in database")
		return err
	}
	wordsString := strings.Join(terms, "\n")

	Common.Reply(s, i, wordsString)

	return nil
}
