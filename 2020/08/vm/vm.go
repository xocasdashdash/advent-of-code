package vm

import "../parser"

//Run Runs a virtual machine until halt or exit
func Run(instructions []parser.Instruction) (int, int, []int) {
	runnedInstructions := make([]int, 0, len(instructions))
	visited := make(map[int]bool)
	accumulator := 0
	pointer := 0
	endInstruction := len(instructions)
	maxPointer := 0
	for !visited[pointer] {
		instruction := instructions[pointer]
		visited[pointer] = true
		runnedInstructions = append(runnedInstructions, pointer)
		switch instruction.OpCode {
		case "acc":
			accumulator = accumulator + instruction.Modifier
			pointer++
		case "nop":
			pointer++
		case "jmp":
			pointer = (pointer + instruction.Modifier)
			if pointer == endInstruction {
				return accumulator, pointer, runnedInstructions
			}
		}
		if pointer > maxPointer {
			maxPointer = pointer
		}
	}
	return accumulator, pointer, runnedInstructions
}
