package utils

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