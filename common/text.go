package common

import "fmt"

var Text *TextManager

type TextManager struct {
	TextKeys map[string]string
}

const TXT_COMBAT_PARRY = "TXT_COMBAT_PARRY"
const TXT_COMBAT_BLOCK = "TXT_COMBAT_BLOCK"

func NewTextManager() *TextManager {
	t := &TextManager{
		TextKeys: map[string]string{
			TXT_COMBAT_PARRY: "Parry with your %v",
			TXT_COMBAT_BLOCK: "Block with your %v",
		},
	}

	return t
}

func (t *TextManager) Format(key string, args ...interface{}) string {
	fmtStr := t.TextKeys[key]
	return fmt.Sprintf(fmtStr, args...)
}
