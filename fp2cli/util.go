package main

import (
	"strings"

	"github.com/wtiger001/fp2server/common"
)

var RefNames = map[string]common.ReferenceType{
	"weapon": common.ReferenceType_ReferenceType_Weapon,
	"armor":  common.ReferenceType_ReferenceType_Armor,
}

func NameToRefType(name string) common.ReferenceType {
	return RefNames[strings.ToLower(name)]
}

func RefTypeToName(t common.ReferenceType) string {
	for k, v := range RefNames {
		if v == t {
			return k
		}
	}
	return ""
}
