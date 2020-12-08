package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"./parser"
	"./vm"
)

func main() {
	instructionsRaw, _ := ioutil.ReadFile("input")
	instructions := strings.Split(strings.TrimSpace(string(instructionsRaw)), "\n")

	
	parsedInstructions := parser.ParseInstructions(instructions)
	accumulator, _, visitedInstructions := vm.Run(parsedInstructions)
	fmt.Printf("Visited instructions: %d\n", len(visitedInstructions))
	fmt.Printf("Accumulator: %d", accumulator)

}
