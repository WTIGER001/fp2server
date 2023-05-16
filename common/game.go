package common

import (
	"path/filepath"
	"time"
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
		onSave: func(item *Character) error {
			item.GameID = g.ID
			return nil
		},
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

func (g *Game) Time() time.Time {
	if g.GameTime == 0 {
		g.GameTime = time.Now().Unix()
	}
	return time.Unix(g.GameTime, 0)
}

// Calculate the time in "n" rounds
func (g *Game) TimeInRounds(n int) time.Time {
	return g.TimeIn(time.Minute * time.Duration(SecondsPerRound*n))
}

// Calculate the game time in
func (g *Game) TimeIn(d time.Duration) time.Time {
	return g.Time().Add(d)
}

// Advance the Game time by a single round
func (g *Game) TimeAdvanceRounds(n int) time.Time {
	return g.TimeAdvance(time.Minute * time.Duration(SecondsPerRound*n))
}

// Advance the game time arbitarily
func (g *Game) TimeAdvance(d time.Duration) time.Time {
	g.GameTime = g.Time().Add(d).Unix()
	return time.Unix(g.GameTime, 0)
}
