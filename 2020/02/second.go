package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	f, _ := os.Open("input")

	s := bufio.NewScanner(f)
	validPassword := 0
	for s.Scan() {
		line := s.Text()
		tokens := strings.Split(line, " ")
		passwordIndexes := strings.Split(tokens[0], "-")
		passwordLow, err := strconv.Atoi(passwordIndexes[0])
		if err != nil {
			panic(err)
		}
		passwordHigh, err := strconv.Atoi(passwordIndexes[1])
		if err != nil {
			panic(err)
		}
		charToCheck := tokens[1][0]
		candidatePassword := tokens[2]
		if passwordHigh > len(candidatePassword) {
			fmt.Println("Line too short")
			continue
		}

		lowChar := tokens[2][passwordLow-1]
		highChar := tokens[2][passwordHigh-1]
		if lowChar == highChar {
			fmt.Println("Both chars match")
			continue
		}
		fmt.Printf("Char to check %c, lowChar %c, highChar %c\n", charToCheck, lowChar, highChar)
		if lowChar == charToCheck || highChar == charToCheck {
			validPassword++
		}
	}
	fmt.Printf("Valid Passwords: %d", validPassword)

}
