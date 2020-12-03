package main

import (
	"bufio"
	"fmt"
	"os"
)

const xMovement = 3
const yMovement = 1
const treeChar = '#'

func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)
	s.Scan()
	initialIndex := 0
	trees := 0
	for s.Scan() {
		initialIndex = initialIndex + xMovement
		//read following line
		line := s.Text()
		initialIndex = initialIndex % len(line)
		//Read new index at current position
		treeCandidate := line[initialIndex]
		//isTree ?
		if treeCandidate == treeChar {
			trees++
		}
	}
	fmt.Printf("Trees %d", trees)

}
