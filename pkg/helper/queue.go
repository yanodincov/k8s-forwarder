package helper

type CircularQueueLimited[T any] struct {
	data []T
	cur  int
}

func NewCircularQueue[T any](vals []T, cur int) *CircularQueueLimited[T] {
	return &CircularQueueLimited[T]{data: vals, cur: cur}
}

func (c *CircularQueueLimited[T]) Next() T {
	c.cur = (c.cur + 1) % len(c.data)
	return c.data[c.cur]
}

func (c *CircularQueueLimited[T]) Prev() T {
	c.cur = (c.cur - 1 + len(c.data)) % len(c.data)
	return c.data[c.cur]
}

func (c *CircularQueueLimited[T]) Current() T {
	return c.data[c.cur]
}

func (c *CircularQueueLimited[T]) CurI() int {
	return c.cur
}

func (c *CircularQueueLimited[T]) SetCurrent(cur int) {
	c.cur = cur
}

func (c *CircularQueueLimited[T]) Len() int {
	return len(c.data)
}

func (c *CircularQueueLimited[T]) Data() []T {
	return c.data
}
