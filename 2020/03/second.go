package main

import (
	"bufio"
	"fmt"
	"os"
)

const xMovement = 3
const yMovement = 1
const treeChar = '#'

type Slope struct {
	dx int
	dy int
}

func main() {

	f, _ := os.Open("input")
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	forestMap := make([][]byte, 0, 10)
	for s.Scan() {
		r := s.Bytes()
		b := make([]byte, len(r), len(r))
		copy(b, r)
		forestMap = append(forestMap, b)
	}
	slopes := []Slope{
		{3, 1},
		{1, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	collisions := make(map[Slope]int, 0)
	total := 1
	for _, s := range slopes {
		posX := 0
		posY := 0
		collisions[s] = 0
		for posY < len(forestMap) {
			posX = (posX + s.dx)
			posX = posX % len(forestMap[posY])
			posY = (posY + s.dy)
			if posY >= len(forestMap) {
				break
			}
			treeCandidate := forestMap[posY][posX]
			if treeCandidate == treeChar {
				collisions[s] = collisions[s] + 1
			}
		}
		total *= collisions[s]
	}
	fmt.Printf("Total %d\n", total)

}
