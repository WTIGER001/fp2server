package common

import (
	"path/filepath"
)

/*
Games are the top level for a lot objects
- Characters
- Encounters
- Chat
- Rolls

but they are kept independant.
*/
var ActiveGame *Game

var Games = &ItemManager[*Game]{
	factory: GameFactory,
	path:    filepath.Join("games"),
	cache:   NewCache[*Game](),
}

var gameCharacters = make(map[string]*ItemManager[*Character])
var gameEncounters = make(map[string]*ItemManager[*Encounter])

func (g *Game) Characters() *ItemManager[*Character] {
	im, found := gameCharacters[g.ID]
	if found {
		return im
	}
	im = &ItemManager[*Character]{
		factory: CharacterFactory,
		path:    filepath.Join("games", g.ID, "characters"),
		cache:   NewCache[*Character](),
	}
	gameCharacters[g.ID] = im
	return im
}

func (g *Game) Encounters() *ItemManager[*Encounter] {
	im, found := gameEncounters[g.ID]
	if found {
		return im
	}
	im = &ItemManager[*Encounter]{
		factory: EncounterFactory,
		path:    filepath.Join("games", g.ID, "encounters"),
		cache:   NewCache[*Encounter](),
	}
	gameEncounters[g.ID] = im
	return im
}
