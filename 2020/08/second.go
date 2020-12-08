package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Instruction struct {
	OpCode   string
	Modifier int
}
type RunnedInstruction struct {
	Pointer     int
	Instruction Instruction
}

func simulate(instructions map[int]Instruction) (int, int, []RunnedInstruction) {
	runnedInstructions := make([]RunnedInstruction, 0, len(instructions))
	visited := make(map[int]bool)
	accumulator := 0
	pointer := 0
	endInstruction := len(instructions)
	maxPointer := 0
	for !visited[pointer] {
		instruction := instructions[pointer]
		visited[pointer] = true
		runnedInstructions = append(runnedInstructions, RunnedInstruction{pointer, instruction})
		switch instruction.OpCode {
		case "acc":
			accumulator = accumulator + instruction.Modifier
			pointer++
		case "nop":
			pointer++
		case "jmp":
			nextPointer := (pointer + instruction.Modifier)
			if nextPointer == endInstruction {
				return accumulator, nextPointer, runnedInstructions
			}
			pointer = nextPointer % len(instructions)
		}
		if pointer > maxPointer {
			maxPointer = pointer
		}
	}
	return accumulator, pointer, runnedInstructions
}
func runnedInstructionsPrinter(instructions []RunnedInstruction) string {
	instructionsString := make([]string, 0, len(instructions))
	for k := range instructions {
		instructionsString = append(instructionsString, fmt.Sprintf("%d - %s %d", k, instructions[k].Instruction.OpCode, instructions[k].Instruction.Modifier))
	}
	return strings.Join(instructionsString, "\n")

}
func runExperiment(instructions map[int]Instruction, instructionToFlipPointer int) (int, error) {

	mutantInstructions := make(map[int]Instruction, len(instructions))
	for k := range instructions {
		mutantInstructions[k] = instructions[k]
	}
	if mutantInstructions[instructionToFlipPointer].OpCode == "nop" {
		mutantInstructions[instructionToFlipPointer] = Instruction{
			OpCode:   "jmp",
			Modifier: instructions[instructionToFlipPointer].Modifier,
		}
	} else if mutantInstructions[instructionToFlipPointer].OpCode == "jmp" {
		mutantInstructions[instructionToFlipPointer] = Instruction{
			OpCode:   "nop",
			Modifier: instructions[instructionToFlipPointer].Modifier,
		}
	} else {
		return 0, fmt.Errorf("Acc is not to experiment with")
	}

	accumulator, pointer, _ := simulate(mutantInstructions)
	if pointer == len(instructions) {
		return accumulator, nil
	}
	return 0, fmt.Errorf("Bad experiment")
}
func main() {
	instructionsRaw, _ := ioutil.ReadFile("input")
	instructions := strings.Split(strings.TrimSpace(string(instructionsRaw)), "\n")

	parsedInstructions := make(map[int]Instruction, len(instructions))
	for k, instruction := range instructions {
		instructionTokens := strings.Split(instruction, " ")
		modifier, _ := strconv.Atoi(instructionTokens[1])
		parsedInstructions[k] = Instruction{
			OpCode:   instructionTokens[0],
			Modifier: modifier,
		}
	}
	accumulator, pointer, runnedInstructions := simulate(parsedInstructions)

	fmt.Printf("Original run:\n Accumulator: %d, pointer: %d, runnedInstructions: \n%s\n", accumulator, pointer, runnedInstructionsPrinter(runnedInstructions))

	experiments := 0
	for k := range runnedInstructions {
		fmt.Printf("Running experiment %d/%d\n", experiments, len(runnedInstructions))
		experiments++
		acc, err := runExperiment(parsedInstructions, runnedInstructions[k].Pointer)
		if err == nil {
			fmt.Printf("SUCCESS. Acumulator: %d\n", acc)
			break
		}
	}

}
