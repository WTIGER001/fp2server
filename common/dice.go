package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Properly checks and records
func Roll(dice *DiceRoll) *DiceRollResults {
	return PlayerInteractions.Roll(dice)
}

func ParseDiceRoll(roll string) *DiceRoll {
	return nil
}

func NewDieRoll(sides int32, exploding bool, reason DiceRollReason) *DiceRoll {
	d := &DiceRoll{
		Reason: reason,
		Dice: []*Die{
			{
				Sides:   sides,
				Amount:  1,
				Explode: exploding,
			},
		},
	}

	return d
}

func RollD10() *DiceRollResults {
	d := &DiceRoll{
		Dice: []*Die{
			{
				Sides:   10,
				Amount:  1,
				Explode: true,
			},
		},
	}

	return d.Roll()
}

func RollOnce(sides int32, exploding bool) int32 {
	v := Random1(sides)
	if v == sides && exploding {
		v = +RollOnce(sides, exploding)
	}
	return v
}

func (d *DiceRoll) Roll() *DiceRollResults {
	// Roll The Dice
	dieTotal := int32(0)
	results := make([]*DiceRollResult, len(d.Dice))
	for i, d := range d.Dice {
		rolls := make([]*DieRollResult, d.Amount)
		for n := 0; n < int(d.Amount); n++ {
			v := RollOnce(d.Sides, d.Explode)
			rolls[n] = &DieRollResult{
				Value: v,
			}
			dieTotal += v
		}
		results[i] = &DiceRollResult{
			Dice:    d,
			Results: rolls,
		}
	}

	// Sum the modifiers
	mods := int32(0)
	for _, m := range d.Modifiers {
		mods += m.Modifier
	}
	total := dieTotal + mods

	return &DiceRollResults{
		Rolls:     results,
		Total:     int32(total),
		Modifiers: d.Modifiers,
	}
}

func (d *Die) Format() string {
	format := fmt.Sprintf("%vD%v", d.Amount, d.Sides)
	if d.Explode {
		format += "X"
	}
	return format
}

func (d *DiceRoll) Format() string {
	// Roll The Dice
	var results []string
	for _, d := range d.Dice {
		results = append(results, d.Format())
	}

	// Sum the modifiers
	mods := int32(0)
	for _, m := range d.Modifiers {
		mods += m.Modifier
	}

	diceFmt := ""
	for i, dFormat := range results {
		if i > 0 && !d.Dice[i].Negative {
			diceFmt += " + "
		} else if i > 0 && d.Dice[i].Negative {
			diceFmt += " - "
		} else if d.Dice[i].Negative {
			diceFmt += "-"
		}
		diceFmt += dFormat
	}
	if mods > 0 {
		diceFmt = fmt.Sprintf("%v + %v", diceFmt, mods)
	} else if mods < 0 {
		absmods := -mods
		diceFmt = fmt.Sprintf("%v - %v", diceFmt, absmods)
	}
	return diceFmt
}

func (d *DiceRollResult) Format() string {
	return d.GetDice().Format()
}

func (d *DiceRollResults) Format() string {
	var results []string
	for _, d := range d.Rolls {
		results = append(results, d.Format())
	}

	// Sum the modifiers
	mods := int32(0)
	for _, m := range d.Modifiers {
		mods += m.Modifier
	}

	diceFmt := strings.Join(results, " + ")
	if mods > 0 {
		diceFmt = fmt.Sprintf("%v + %v", diceFmt, mods)
	} else if mods < 0 {
		absmods := -mods
		diceFmt = fmt.Sprintf("%v - %v", diceFmt, absmods)
	}
	return diceFmt
}

func (d *DiceRoll) AppendMod(tag string, mod int32) {
	d.Modifiers = append(d.Modifiers, &RollModifier{
		Modifier: mod,
		Tag:      tag,
	})
}

func (d *DiceRoll) FindModifier(tag string) int32 {
	for _, m := range d.Modifiers {
		if m.Tag == tag {
			return m.Modifier
		}
	}
	return 0
}

func (d *DiceRollResults) FindModifier(tag string) int32 {
	for _, m := range d.Modifiers {
		if m.Tag == tag {
			return m.Modifier
		}
	}
	return 0
}

func (d *DiceRoll) addModifier(tag string, mod int32) {
	d.Modifiers = append(d.Modifiers, &RollModifier{
		Tag:      tag,
		Modifier: mod,
	})
}

var diceregex = regexp.MustCompile(`(?i)([-+]?)[\s]*(\d*)([d]?)(\d{1,10})([x]?)`)

func parseDiceExpression(exp string) *DiceRoll {
	roll := &DiceRoll{}

	matches := diceregex.FindAllStringSubmatch(exp, -1)
	for _, match := range matches {
		all := match[0]
		plusMinus := match[1]
		quantStr := match[2]
		dStr := match[3]
		facesStr := match[4]
		xStr := match[5]

		if dStr == "" {
			toParse := strings.ReplaceAll(strings.ReplaceAll(all, " ", ""), "+", "")
			mod, _ := strconv.Atoi(toParse)
			roll.addModifier("", int32(mod))
		} else {
			numberOfDice := 1
			if len(quantStr) > 0 {
				numberOfDice, _ = strconv.Atoi(quantStr)
			}
			faces, _ := strconv.Atoi(facesStr)
			isNegative := plusMinus == "-"
			isExploding := len(xStr) > 0
			roll.Dice = append(roll.Dice, &Die{
				Amount:   int32(numberOfDice),
				Sides:    int32(faces),
				Negative: isNegative,
				Explode:  isExploding,
			})
		}
	}
	return roll
}

/*
parse(exp: string): DiceRoll {
        let roll = new DiceRoll()
        roll.expression = exp

        // Expression to identify dice rolls or modifiers
        // const regex = /([-+]?[\s]*\d*[d]?\d{1,3})/gi
        const regex = /([-+]?)[\s]*(\d*)([d]?)(\d{1,10})/gi

        let match;
        while ((match = regex.exec(exp)) !== null) {
            // This is necessary to avoid infinite loops with zero-width matches
            if (match.index === regex.lastIndex) {
                regex.lastIndex++;
            }
            let all = match[0]

            // Group 1 is the + or - sign
            let plusMinus = match[1]

            // Group 2 is the number in front of the dice
            let quantStr = match[2]

            // Group 3 is the 'd' indicator, which signifies a die
            let dStr = match[3]

            // Group 4 is the modifier or dice classifier
            let facesStr = match[4]

            if (dStr.length == 0) {
                // This is a modifier
                let mod = parseInt(all)
                // if (plusMinus == '-') {
                //     mod = -mod
                // }
                roll.addModifier(mod)
            } else {
                let numberOfDice = quantStr.length == 0 ? 1 : parseInt(quantStr)
                roll.addDice(numberOfDice, parseInt(facesStr), plusMinus == '-')
            }
        }
        return roll
    }

*/
