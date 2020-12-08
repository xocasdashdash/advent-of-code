package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	instructionsRaw, _ := ioutil.ReadFile("input")
	instructions := strings.Split(strings.TrimSpace(string(instructionsRaw)), "\n")

	pointer := 0
	visited := make(map[int]bool, len(instructions))
	accumulator := 0
	visitedInstructions := 0
	for !visited[pointer] {
		instruction := instructions[pointer]
		visited[pointer] = true
		visitedInstructions++

		instructionTokens := strings.Split(instruction, " ")
		switch instructionTokens[0] {
		case "acc":
			increase, _ := strconv.Atoi(instructionTokens[1])
			accumulator = accumulator + increase
			fallthrough
		case "nop":
			pointer++
		case "jmp":
			jump, _ := strconv.Atoi(instructionTokens[1])
			pointer = (pointer + jump) % len(instructions)
		}
	}
	fmt.Printf("Visited instructions: %d\n", visitedInstructions)
	fmt.Printf("Accumulator: %d", accumulator)

}
