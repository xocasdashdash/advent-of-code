package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
const (
	UP = "up"
	DOWN = "down"
	FORWARD = "forward"
)
type Instruction struct {
	Direction string
	Value int
}

func parseInstruction(input string) *Instruction {
	parts := strings.Split(strings.ToLower(input), " ")
	if len(parts) != 2 {
		return nil
	}
	value, _ := strconv.Atoi(parts[1])
	return &Instruction{
		Direction: parts[0],
		Value: value,
	}
}
func parseInstructions(input []string) []Instruction{
	result := make([]Instruction, len(input))
	for k := range input{
		result[k] = *parseInstruction(input[k])
	}
	return result
}
func calcDirections(instructions []Instruction, useAim bool) int {
	horizontal := 0
	vertical := 0
	aim := 0
	for _, instruction := range instructions {
		direction := 1
		switch instruction.Direction {
		case UP:
			direction = -1
			fallthrough
		case DOWN:
			if useAim {
				aim = aim + instruction.Value*direction
			}else {
				vertical = vertical+instruction.Value*direction
			}
		case FORWARD:
			horizontal = horizontal + instruction.Value 
			if useAim {
				vertical = vertical + instruction.Value*aim
			}
		}
	}
	return horizontal*vertical
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	
	instructions := parseInstructions(trimmedInput)
	
	fmt.Println(calcDirections(	instructions,false))
	fmt.Println(calcDirections(	instructions,true))
	
}
