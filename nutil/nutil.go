package nutil

import "github.com/google/uuid"

func GetGUUID() string {
	return uuid.New().String()
}
