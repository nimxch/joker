package main

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrQueueFull  = errors.New("Queue is full!")
	ErrQueueEmpty = errors.New("Queue is empty!")
)

type Queue[T any] struct {
	data     []T
	head     int
	tail     int
	size     int
	capacity int
}

func NewQueue[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		fmt.Println("Invalid capacity, setting to default")
		capacity = math.MaxInt16
	}

	return &Queue[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

func (q Queue[T]) Enqueue(value T) error {
	if q.size >= q.capacity {
		return ErrQueueFull
	}
	q.data[q.tail] = value
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return nil
}

func (q Queue[T]) Dequeue(value T) (T, error) {
	var zero T
	if q.size <= 0 {
		return zero, ErrQueueEmpty
	}
	value = q.data[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return value, nil
}

func (q *Queue[T]) Peek() (T, error) {
	var zero T

	if q.size == 0 {
		return zero, ErrQueueEmpty
	}

	return q.data[q.head], nil
}

func (q *Queue[T]) Size() int {
	return q.size
}
