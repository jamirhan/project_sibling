package tclient

import (
	"strings"

	"golang.org/x/exp/slices"
)

var commandType = "bot_command"

type ParsedCommand struct {
	name   string
	offset int
	length int
}

func ParseCommandsFromMessage(message Message, allowedSuffixes []string) []ParsedCommand {
	var res []ParsedCommand
	runeText := []rune(message.Text)
	for _, entity := range message.MessageEntities {
		if entity.Type != commandType {
			continue
		}
		command := runeText[entity.Offset : entity.Offset+entity.Length]
		before, after, found := strings.Cut(string(command), "@")
		if !found || slices.Contains(allowedSuffixes, after) {
			res = append(res, ParsedCommand{
				name:   before,
				offset: entity.Offset,
				length: entity.Length,
			})
			continue
		}
	}
	return res
}