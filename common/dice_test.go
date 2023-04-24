package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiceRoll(t *testing.T) {
	d := &DiceRoll{
		Dice: []*Die{
			{
				Sides:  6,
				Amount: 3,
			},
		},
		Modifiers: []*RollModifier{
			{
				Modifier: 3,
			},
		},
	}

	r := d.Roll()
	f := d.Format()

	assert.Greater(t, r.Total, int32(0))
	assert.NotEmpty(t, f)
}
