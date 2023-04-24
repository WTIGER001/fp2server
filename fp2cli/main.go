package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/wtiger001/fp2server/common"
)

//	func completer(d prompt.Document) []prompt.Suggest {
//		s := []prompt.Suggest{
//			{Text: "users", Description: "Store the username and age"},
//			{Text: "articles", Description: "Store the article text posted by user"},
//			{Text: "comments", Description: "Store the text commented to articles"},
//		}
//		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
//	}
var currentCharacter string
var processors = []ActionProcessor{
	&ListReferenceItemsAP{},
}

func main() {
	connect()
}

func three() {
	for {
		fmt.Println("What Now?")
		action := prompt.Input(currentCharacter+"> ", completeActions,
			prompt.OptionTitle("Fp2 Client"),
			prompt.OptionCompletionOnDown(),
		)

		if strings.HasPrefix(action, "quit") {
			fmt.Println("Goodbye!")
			MessageBus.Quit()
			break
		}

		if strings.HasPrefix(action, "play") {
			args := strings.Split(action, " ")
			currentCharacter = args[1]
			continue
		}

		if strings.HasPrefix(action, "chat ") {
			msg := strings.TrimSpace(action[4:])
			MessageBus.Send(&common.Fp2Message{
				Data: &common.Fp2Message_Chat{
					Chat: &common.Chat{
						Contents: msg,
					},
				},
			})
			continue
		}
		args := strings.Split(action, " ")
		ap := Actions.GetProcessor(args)
		if ap != nil {
			result := ap.Execute(args)
			if !result {
				fmt.Println("Error...")
			}
		} else {
			fmt.Println("I understand. Please Try Again.")

		}

	}
}

func completeActions(d prompt.Document) []prompt.Suggest {
	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) < 2 {
		return Actions.InitialSuggestions(d)
	}
	return Actions.NextSuggestions(d)
}
