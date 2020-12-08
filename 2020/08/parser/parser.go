package parser

import (
	"strconv"
	"strings"
)

type Instruction struct {
	OpCode   string
	Modifier int
}

func ParseInstructions(input []string) []Instruction {
	parsedInstructions := make([]Instruction, len(input))
	for k, instruction := range input {
		instructionTokens := strings.Split(instruction, " ")
		modifier, _ := strconv.Atoi(instructionTokens[1])
		parsedInstructions[k] = Instruction{
			OpCode:   instructionTokens[0],
			Modifier: modifier,
		}
	}
	return parsedInstructions
}
