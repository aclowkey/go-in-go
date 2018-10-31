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
	data  []*Grid
	index int
}

func MakeBoardQueue(size int) *BoardQueue {
	data := make([]*Grid, size)
	return &BoardQueue{
		size,
		data,
		0,
	}
}

func (queue *BoardQueue) Enqueue(board *Grid) error {
	// Make a snapshot of the grid to store as history
	snapshot := make(Grid, len(*board))
	for i := 0; i < len(*board); i++ {
		snapshot[i] = make([]Cell, len(*board))
		for y := 0; y < len(*board); y++ {
			toClone := (*board)[i][y]
			snapshot[i][y] = Cell{
				toClone.piece,
				toClone.liberty,
			}
		}
	}
	// Store only last queue.size moves
	if queue.index >= queue.size {
		// Shift left
		for i := 0; i < queue.index-1; i++ {
			queue.data[i] = queue.data[i+1]
		}
		queue.index--
	}
	// Move to the next available spot
	queue.data[queue.index] = &snapshot
	queue.index++
	return nil
}
