package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
	tu := NewTestUtil()
	tu.InitInMemStorage()

	allGames, err := Games.GetAll()
	assert.Empty(t, allGames)
	assert.Nil(t, err)

	err = Games.Save(&Game{
		ID:   "temp",
		Name: "temp game",
	}, "temp")
	assert.Nil(t, err)

	allGames, err = Games.GetAll()
	assert.NotEmpty(t, allGames)
	assert.Nil(t, err)

	g := allGames[0]
	assert.NotEmpty(t, g.ID)
	assert.NotEmpty(t, g.Name)

	c := &Character{
		ID:   "tempcar",
		Name: "Saved",
	}

	err = g.Characters().Save(c, c.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.GameID)

}
