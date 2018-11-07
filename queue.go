package main

import (
	"strings"
)

type MovementQueue struct {
	data []*Move
	head *Move
}

func MakeMovementQueue() *MovementQueue {
	return &MovementQueue{
		[]*Move{},
		nil,
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
	queue.data = append(queue.data, move)
	queue.head = move
	return nil
}

type BoardQueue struct {
	data []*Grid
	head *Grid
}

func MakeBoardQueue() *BoardQueue {
	return &BoardQueue{
		[]*Grid{},
		nil,
	}
}

func (queue *BoardQueue) Enqueue(board *Grid) error {
	// Make a snapshot of the grid to store as history
	snapshot := board.Clone()
	queue.data = append(queue.data, &snapshot)
	queue.head = &snapshot
	return nil
}

func (queue *BoardQueue) IsKo(board *Grid) bool {
	// Ko is when move n == n-2
	currentBoard := *board
	previousBoard := *queue.data[len(queue.data)-2]
	for i := 0; i < len(currentBoard); i++ {
		for y := 0; y < len(currentBoard); y++ {
			if currentBoard[i][y].piece != previousBoard[i][y].piece {
				return false
			}
		}
	}
	return true
}
