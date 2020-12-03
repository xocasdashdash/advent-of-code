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
		passwordLimits := strings.Split(tokens[0], "-")
		passwordMin, err := strconv.Atoi(passwordLimits[0])
		if err != nil {
			panic(err)
		}
		passwordMax, err := strconv.Atoi(passwordLimits[1])
		if err != nil {
			panic(err)
		}
		charCheck := tokens[1][0]
		c := strings.Count(tokens[2], string(charCheck))
		if !(c < passwordMin || c > passwordMax) {
			validPassword++
		}
	}
	fmt.Printf("Valid Passwords: %d", validPassword)

}
