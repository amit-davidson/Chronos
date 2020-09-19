package utils

import (
	"github.com/satori/go.uuid"
)

func GetUUID() string {
	uuidRes := uuid.Must(uuid.NewV4())
	return uuidRes.String()
}

type Counter struct {
	count int
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) GetNext() int {
	c.count += 1
	return c.count
}