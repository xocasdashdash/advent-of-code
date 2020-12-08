package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"./parser"
	"./vm"
)

func runnedInstructionsPrinter(instructions []parser.Instruction, runnedInstructions []int) string {
	instructionsString := make([]string, 0, len(instructions))
	for k, v := range runnedInstructions {
		instructionsString = append(instructionsString, fmt.Sprintf("%d - %s %d", k, instructions[v].OpCode, instructions[v].Modifier))
	}
	return strings.Join(instructionsString, "\n")

}
func runExperiment(instructions []parser.Instruction, instructionToFlipPointer int) (int, error) {

	mutantInstructions := make([]parser.Instruction, len(instructions))
	copy(mutantInstructions, instructions)
	if mutantInstructions[instructionToFlipPointer].OpCode == "nop" {
		mutantInstructions[instructionToFlipPointer] = parser.Instruction{
			OpCode:   "jmp",
			Modifier: instructions[instructionToFlipPointer].Modifier,
		}
	} else if mutantInstructions[instructionToFlipPointer].OpCode == "jmp" {
		mutantInstructions[instructionToFlipPointer] = parser.Instruction{
			OpCode:   "nop",
			Modifier: instructions[instructionToFlipPointer].Modifier,
		}
	} else {
		return 0, fmt.Errorf("Acc is not to experiment with")
	}

	accumulator, pointer, _ := vm.Run(mutantInstructions)
	if pointer == len(instructions) {
		return accumulator, nil
	}
	return 0, fmt.Errorf("Bad experiment")
}

func main() {
	instructionsRaw, _ := ioutil.ReadFile("input")
	instructions := strings.Split(strings.TrimSpace(string(instructionsRaw)), "\n")
	parsedInstructions := parser.ParseInstructions(instructions)
	_, _, runnedInstructions := vm.Run(parsedInstructions)
	experiments := 0
	for k := range runnedInstructions {
		experiments++
		acc, err := runExperiment(parsedInstructions, runnedInstructions[k])
		if err == nil {
			fmt.Printf("SUCCESS on %d/%d experiment. Acumulator: %d\n", experiments, len(runnedInstructions), acc)
			break
		}
	}

}
