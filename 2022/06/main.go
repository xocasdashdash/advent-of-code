package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func findUniqueSequenceOfLength(input string, length int) int {

	chars := strings.Split(input, "")
	for i := length; i < len(chars); i++ {
		uniqueMap := make(map[string]bool)
		// TODO: This is too brute force, need a datastruct that
		// keeps uniqueness and also can tell you how many unique
		// characters it's keeping.
		for _, c := range chars[i : i+length] {
			uniqueMap[c] = true
		}
		if len(uniqueMap) == length {
			return i + length
		}
	}
	return -1
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	for _, l := range trimmedInput {
		fmt.Println("Part1", findUniqueSequenceOfLength(l, 4))
	}
	for _, l := range trimmedInput {
		fmt.Println("Part2", findUniqueSequenceOfLength(l, 14))
	}
}
