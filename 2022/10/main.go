package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string

//go:embed smallInput
var smallInput string
var testMode = flag.Bool("test", false, "Set to run using the testInput")
var smallMode = flag.Bool("small", false, "Test with the small input")

type Instruction struct {
	Kind           string
	Value          int
	TotalCycles    int
	NumberOfCycles int
}

func parseInstructions(lines []string) []Instruction {

	result := make([]Instruction, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		ins := Instruction{Kind: "noop", TotalCycles: 0}
		if line != "noop" {
			fmt.Sscanf(line, "%s %d", &ins.Kind, &ins.Value)
			ins.TotalCycles = 1
		}
		result = append(result, ins)
	}
	return result
}
func part1(instructions []Instruction) []int {
	result := make([]int, 0, 1)
	register := 1
	pc := 0
	var increment int
	for pc < len(instructions) {
		register = register + increment
		result = append(result, register)
		instructions[pc].NumberOfCycles++
		increment = 0
		if instructions[pc].NumberOfCycles > instructions[pc].TotalCycles {
			// We increase the value of current cycleCounter
			increment = instructions[pc].Value
			// We start with the next instruction
			pc++
		}

	}
	return result
}
func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		if *smallMode {
			input = smallInput
		} else {
			input = testInput
		}
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	instructions := parseInstructions(trimmedInput)
	p1 := part1(instructions)

	signal := 0
	for i := 19; i < len(p1); i += 40 {
		signal += ((i + 1) * p1[i])
	}
	fmt.Println("part1", signal)
	fmt.Println("part2")
	spritePos := 1
	index := 1
	char := "."
	for i := 1; i < len(p1); i++ {
		char = "."
		if index >= spritePos && index < (spritePos+3) {
			char = "#"
		}
		fmt.Printf("%s", char)
		if index > 39 {
			fmt.Printf("\n")
			index = 0
		}
		spritePos = p1[i]
		index++
	}
	fmt.Printf("%s\n", char)
	fmt.Println("Took", time.Since(start))
}
