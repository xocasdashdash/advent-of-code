package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)

	totalYes := 0
	groupAnswers := make(map[byte]int, 26)
	currentGroup := make([]byte, 0, 26)
	groupLength := 0
	for s.Scan() {
		b := s.Bytes()
		if len(b) == 0 {
			currTotal := totalYes
			for _, v := range groupAnswers {
				if v == groupLength {
					totalYes++
				}
			}
			groupAnswers = make(map[byte]int, 26)
			currentGroup = make([]byte, 0, 26)
			groupLength = 0
			continue
		}
		for _, v := range b {
			groupAnswers[v]++
		}
		currentGroup = append(currentGroup, b...)
		groupLength++
	}
	currTotal := totalYes

	for _, v := range groupAnswers {
		if v == groupLength {
			totalYes++
		}
	}
	fmt.Printf("Total yes: %d", totalYes)
}
