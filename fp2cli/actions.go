package main

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/c-bata/go-prompt"
	"github.com/wtiger001/fp2server/common"
)

var Actions = &ActionProcessManager{
	processors: []ActionProcessor{
		&ListReferenceItemsAP{},
	},
}

type ActionProcessManager struct {
	processors []ActionProcessor
}

func (apm *ActionProcessManager) Suggest(d prompt.Document) []prompt.Suggest {
	text := d.CurrentLine()
	args := strings.Split(text, " ")

	if len(args) == 1 {
		return apm.InitialSuggestions(d)
	}
	return nil
}

func (apm *ActionProcessManager) GetProcessor(args []string) ActionProcessor {
	for _, p := range apm.processors {
		initial := p.InitialSuggestions()
		for _, s := range initial {
			if strings.EqualFold(s.Text, args[0]) {
				return p
			}
		}
	}
	return nil
}

func (apm *ActionProcessManager) NextSuggestions(d prompt.Document) []prompt.Suggest {
	text := d.CurrentLine()
	args := strings.Split(text, " ")
	p := apm.GetProcessor(args)
	if p != nil {
		_, s := p.Build(text)
		return s
	}
	return nil
}

func (apm *ActionProcessManager) InitialSuggestions(d prompt.Document) []prompt.Suggest {
	var rtn []prompt.Suggest
	for _, p := range apm.processors {
		rtn = append(rtn, p.InitialSuggestions()...)
	}
	return prompt.FilterHasPrefix(rtn, d.CurrentLineBeforeCursor(), true)
}

type ActionProcessor interface {
	InitialSuggestions() []prompt.Suggest
	Build(text string) (bool, []prompt.Suggest)
	Execute(args []string) bool
}

type ListReferenceItemsAP struct {
}

func (ap *ListReferenceItemsAP) InitialSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "list", Description: "List something <type> <name or id>"},
	}
}

// list weapons
func (ap *ListReferenceItemsAP) Build(text string) (bool, []prompt.Suggest) {
	args := strings.Split(text, " ")

	types := []prompt.Suggest{
		{Text: "weapon", Description: "List all weapons"},
		{Text: "term", Description: "Describe a Game Term"},
		// {Text: "character", Description: "Describe a Character"},
		// {Text: "character", Description: "Describe a Character"},
	}

	if len(args) == 1 {
		return false, types
	}

	if len(args) == 2 {
		return false, prompt.FilterContains(types, args[1], true)
	}

	if len(args) == 3 {
		if len(args) == 3 {
			switch args[2] {
			case "weapon":
				return true, nil
			}
		}

	}
	// Look up the type
	return false, nil
}

func (ap *ListReferenceItemsAP) Execute(args []string) bool {
	if len(args) == 2 {
		switch args[1] {
		case "weapon":
			refs, _ := ModelGetAll(common.ModelType_ModelType_RefWeapon)
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: "ID"},
					{Align: simpletable.AlignLeft, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "Damage1H"},
				},
			}

			for _, model := range refs {
				row := model.GetRefWeapon()
				r := []*simpletable.Cell{
					{Text: row.ID},
					{Text: row.Name},
					{Text: row.Damage1H, Align: simpletable.AlignCenter},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}

			// table.SetStyle(simpletable.StyleCompactLite)
			fmt.Println(table.String())

			return true
		}

	}
	return false
}
