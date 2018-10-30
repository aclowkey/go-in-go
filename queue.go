package main

import (
	"strings"
)

type MovementQueue struct {
	size  int
	data  []*Move
	index int
}

func MakeMovementQueue(size int) *MovementQueue {
	data := make([]*Move, size)

	return &MovementQueue{
		size,
		data,
		0,
	}
}

func (queue *MovementQueue) String() string {
	var str strings.Builder

	for i := 0; i < len(queue.data); i++ {
		str.WriteString(queue.data[i].String() + "\n")
	}

	return str.String()
}

// Queue is FIFO
func (queue *MovementQueue) Enqueue(move *Move) error {
	if queue.index >= queue.size {
		// Shift left
		for i := 0; i < queue.index-1; i++ {
			queue.data[i] = queue.data[i+1]
		}
		queue.index--
	}
	// Move to the next available spot
	queue.data[queue.index] = move
	queue.index++
	return nil
}

type BoardQueue struct {
	size  int
	data  []*[][]Cell
	index int
}

func MakeBoardQueue(size int) *BoardQueue {
	data := make([]*[][]Cell, size)
	return &BoardQueue{
		size,
		data,
		0,
	}
}

func (queue *BoardQueue) Enqueue(board *[][]Cell) error {
	if queue.index >= queue.size {
		// Shift left
		for i := 0; i < queue.index-1; i++ {
			queue.data[i] = queue.data[i+1]
		}
		queue.index--
	}
	// Move to the next available spot
	queue.data[queue.index] = board
	queue.index++

	return nil
}
