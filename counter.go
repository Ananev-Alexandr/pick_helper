package main

import "sync/atomic"

type AtomicCounter struct {
	counter int64
}

// Метод для безопасного увеличения счетчика
func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.counter, 1)
}

// Метод для получения текущего значения счетчика
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.counter)
}

// Функция обнуления счетчика
func (c *AtomicCounter) Reset() {
	atomic.StoreInt64(&c.counter, 0)
}
