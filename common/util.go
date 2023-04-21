package common

import (
	"github.com/google/uuid"
)

func GenerateID() string {
	id := uuid.New()
	return id.String()
}

func First(items ...string) string {
	for _, s := range items {
		if s != "" {
			return s
		}
	}
	return ""
}

func FirstInt32(items ...int32) int32 {
	for _, s := range items {
		if s != -999 {
			return s
		}
	}
	return -999
}
