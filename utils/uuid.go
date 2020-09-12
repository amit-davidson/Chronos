package utils

import (
	"github.com/satori/go.uuid"
)

func GetUUID() string {
	uuidRes := uuid.Must(uuid.NewV4())
	return uuidRes.String()
}
