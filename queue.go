package main

import (
	"fmt"
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
	if queue.index >= queue.size {
		// Shift left
		for i := 0; i < queue.index-1; i++ {
			queue.data[i] = queue.data[i+1]
		}
		queue.index--
	}
	// Move to the next available spot
	queue.data[queue.index] = board
	fmt.Printf("==========       Storing at       ==========\n")
	fmt.Printf("[%p] \t = \t %p\n", &queue.data[queue.index], queue.data[queue.index])
	fmt.Printf("=============================================\n")
	queue.index++
	fmt.Printf("==========Board state after inqueue==========\n")
	for i := 0; i < len(queue.data); i++ {
		if queue.data[i] != nil {
			fmt.Printf("[%p] \t = \t%p\n", &queue.data[i], queue.data[i])
			cell := queue.data[i]
			fmt.Println(PrintGrid(false, cell))
		}
	}
	fmt.Printf("=============================================\n")
	return nil
}
