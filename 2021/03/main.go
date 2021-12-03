package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

type Node struct {
	Bit int
	n   []*Node
}
type Numbers struct {
	Zeroes *[]Node
	Ones   *[]Node
}

func parseLine(s string) []int {
	result := make([]int, len(s))
	for k, c := range s {
		value, _ := strconv.Atoi(string(c))
		result[k] = value
	}
	return result
}
func parseInput(input []string) [][]int {

	result := make([][]int, len(input), len(input))
	for k, s := range input {
		if len(strings.TrimSpace(s)) > 0 {
			result[k] = parseLine(s)
		}
	}
	return result
}

func flipArray(input [][]int) [][]int {
	result := make([][]int, len(input[0]))
	for i, line := range input {
		for j, bit := range line {
			if len(result[j]) == 0 {
				result[j] = make([]int, len(input))
			}
			result[j][i] = bit
		}
	}
	return result
}
func arrayToInt(input []int)int {
	result:=0
	for k, v := range input {
		result = result + (1 << (len(input) - k-1))*v
	}
	return result
}
func part1(input [][]int) int {

	output := make([]int, len(input[0]), len(input[0]))
	for j := range input[0] {
		for i := range input {
			output[j] = output[j] + input[i][j]
		}
	}
	gamma := 0
	for k, v := range output {
		if (2 * v) > len(input) {
			gamma = gamma + (1 << (len(output) - k-1))
		}
	}
	epsilon := ((1 << len(output)) - 1) - gamma

	return gamma * epsilon
}

func findBitAndFilter(input [][]int, index int,valueToFilter func([][]int,int) int) []int {
	val := valueToFilter(input,index)
	filtered := make([][]int, 0, len(input))
	for _, number := range input {
		if number[index] == val {
			filtered = append(filtered, number)
		}
	}
	if len(filtered) > 1 {
		return findBitAndFilter(filtered, index+1, valueToFilter)
	}
	return filtered[0]
}
func filterForCO2(input [][]int,index int) int {
	sum := 0
	for _,v := range input{
		sum += v[index]
		if 2*sum >= len(input) {
			return 0
		}
	}
	return 1
}
func filterForO2(input [][]int,index int) int{
	sum := 0
	for _,v := range input{
		sum += v[index]
		if 2*sum >= len(input) {
			return 1
		}
	}
	return 0
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println("part1", part1(parseInput(trimmedInput)))
	fmt.Println("part2", part2(parseInput(trimmedInput)))

}

func part2(input [][]int) int {

	oxygen := arrayToInt(findBitAndFilter(input,0,filterForO2))
	co2 := arrayToInt(findBitAndFilter(input, 0, filterForCO2))
	fmt.Println("oxygen", oxygen)
	fmt.Println("co2", co2)
	return co2*oxygen
}
