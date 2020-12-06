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
	groupAnswers := make(map[byte]bool, 26)
	currentGroup := make([]byte, 0, 26)
	for s.Scan() {
		b := s.Bytes()
		if len(b) == 0 {
			currTotal := totalYes
			for _, v := range groupAnswers {
				if v {
					totalYes++
				}
			}
			fmt.Printf("Group yes: %d\n", totalYes-currTotal)
			fmt.Printf("Current group %s\n", currentGroup)
			groupAnswers = make(map[byte]bool, 26)
			currentGroup = make([]byte, 0, 26)
			continue
		}
		for _, v := range b {
			groupAnswers[v] = true
		}
		currentGroup = append(currentGroup, b...)
	}
	currTotal := totalYes

	for _, v := range groupAnswers {
		if v {
			totalYes++
		}
	}
	fmt.Printf("Last group yes: %d\n", totalYes-currTotal)
	fmt.Printf("Last group %s\n", currentGroup)
	fmt.Printf("Total answers: %d", totalYes)
}
