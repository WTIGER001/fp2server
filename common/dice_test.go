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

func TestParseDice(t *testing.T) {
	exp1 := "1d10"
	exp2 := "1D10"
	exp3 := "D10"
	exp4 := "d10"

	f1 := parseDiceExpression(exp1).Format()
	f2 := parseDiceExpression(exp2).Format()
	f3 := parseDiceExpression(exp3).Format()
	f4 := parseDiceExpression(exp4).Format()

	assert.Equal(t, "1D10", f1)
	assert.Equal(t, "1D10", f2)
	assert.Equal(t, "1D10", f3)
	assert.Equal(t, "1D10", f4)

	exp5 := "1d10 + 4d6 + 5 - 2d5"
	roll5 := parseDiceExpression(exp5)
	assert.Equal(t, len(roll5.Dice), 3)
	assert.Equal(t, len(roll5.Modifiers), 1)
	f5 := roll5.Format()
	assert.Equal(t, "1D10 + 4D6 - 2D5 + 5", f5)

	exp6 := "-1d10 + 4d6 + 5 - 2d5"
	f6 := parseDiceExpression(exp6).Format()
	assert.Equal(t, "-1D10 + 4D6 - 2D5 + 5", f6)

	exp7 := "D10 + 5"
	f7 := parseDiceExpression(exp7).Format()
	assert.Equal(t, "1D10 + 5", f7)

	exp8 := "D10-5"
	f8 := parseDiceExpression(exp8).Format()
	assert.Equal(t, "1D10 - 5", f8)

	exp9 := "D10x-5"
	f9 := parseDiceExpression(exp9).Format()
	assert.Equal(t, "1D10X - 5", f9)

	exp10 := "D10X-5"
	f10 := parseDiceExpression(exp10).Format()
	assert.Equal(t, "1D10X - 5", f10)
}
