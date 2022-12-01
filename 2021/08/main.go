package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func parseInput(input []string)[][][]string {
	result := make([][][]string,len(input), len(input))
	for k, v := range input {
		parts := strings.Split(v, "|")
		result[k] = make([][]string, len(parts))
		for i, p := range parts {
			result[k][i] = strings.Split(strings.TrimSpace(p)," ")
		}
	}
	return result
}
func part2(input [][][]string) int {
	return 0
}
func part1(input [][][]string) int{

	resultMap := make(map[int]int)
	for _, v := range input {
		for _, number := range v[1] {
			resultMap[len(number)] = resultMap[len(number)] + 1
		}
	}
	result := 0
	for k, v := range resultMap {
		switch k {
		case 2,3,4,7:
			result += v
		default:
		} 
		
	}
	return result
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	parsedInput := parseInput(trimmedInput)
	fmt.Println("part1", part1(parsedInput))
}
