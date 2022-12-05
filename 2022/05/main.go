package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var part2 = flag.Bool("part2", false, "Set to configure the logic for part2.")

type Crate string
type Stack struct {
	Crates []Crate
}

func parseStacks(stackStr string) map[string]Stack {
	lines := strings.Split(stackStr, "\n")
	numLines := len(lines)
	result := make(map[string]Stack, 10)
	var stackInput [][]string
	stackInput = make([][]string, numLines-1)
	for i := 0; i < numLines-1; i++ {
		stackInput[i] = strings.Split(lines[i], "")
	}
	for k, c := range strings.Split(lines[numLines-1], "") {
		if strings.TrimSpace(c) == "" {
			continue
		}
		crates := make([]Crate, 0, 1)
		for j := 0; j < len(stackInput); j++ {
			v := stackInput[j][k]
			if strings.TrimSpace(v) == "" {
				continue
			}
			crates = append(crates, Crate(v))
		}
		result[c] = Stack{Crates: crates}

	}
	fmt.Printf("Stacks: %+v\n", result)
	return result
}

type Instructions struct {
	Quantity    int
	Origin      string
	Destination string
}

func parseInstructions(instructions string) []Instructions {
	parts := strings.Split(instructions, "\n")
	result := make([]Instructions, len(parts))

	for i, l := range parts {
		n, err := fmt.Sscanf(l, "move %d from %s to %s", &result[i].Quantity, &result[i].Origin, &result[i].Destination)
		if err != nil {
			fmt.Printf("Parsed %d on line %s\n", n, l)
			panic(err)
		}
	}
	fmt.Printf("Instructions: %+v\n", result)
	return result

}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	stacksAndInstructions := strings.Split(string(input), "\n\n")
	s := parseStacks(stacksAndInstructions[0])
	fmt.Println("")
	for n, i := range parseInstructions(stacksAndInstructions[1]) {
		fmt.Printf("instruction(%d):%+v\n", n, i)
		originStack := s[i.Origin]
		cratesToMove := s[i.Origin].Crates[:i.Quantity]
		originStack.Crates = s[i.Origin].Crates[i.Quantity:]
		s[i.Origin] = originStack
		if !*part2 {
			for i, j := 0, len(cratesToMove)-1; i < j; i, j = i+1, j-1 {
				cratesToMove[i], cratesToMove[j] = cratesToMove[j], cratesToMove[i]
			}
		}
		destinationStack := s[i.Destination]
		// TODO: Undertand why using append did not work as expected
		newCrates := make([]Crate, len(cratesToMove)+len(destinationStack.Crates))
		for i := range cratesToMove {
			newCrates[i] = cratesToMove[i]
		}
		for j := range destinationStack.Crates {
			newCrates[j+len(cratesToMove)] = destinationStack.Crates[j]
		}
		destinationStack.Crates = newCrates
		s[i.Destination] = destinationStack

	}
	output := make([]string, len(s)+1)
	for i := range s {
		v, _ := strconv.Atoi(i)
		output[v] = string(s[i].Crates[0])
	}
	fmt.Println("Part1", strings.TrimSpace(strings.Join(output, "")))
}
