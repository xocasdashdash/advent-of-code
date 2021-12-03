package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func parseLine(s string)[]int{
	result := make([]int, len(s))
	for k, c := range s {
		value, _ := strconv.Atoi(string(c))
		result[k] = value
	}
	return result
}
func parseInput(input []string)[][]int{
	
	result := make([][]int, len(input), len(input))
	for k, s := range input {
		result[k] = parseLine(s)
	}
	return result
}


func flipArray (input [][]int)[][]int {
	result := make([][]int,len(input[0]))
	for i,line := range input{
		for j, bit := range line {
			if (len(result[j]) == 0) {
				result[j] = make([]int, len(input))
			}
			result[j][i] = bit
		}

	}
	return result
}
func findMostCommonBit(input []int) int{
	sum := 0
	for _, v := range input {
		sum += v
		if sum >= len(input)/2 {
			return 1
		}
	}
	return 0
}
func calcPowerFactors(input [][]int ) (int,int){
	gamma := 0
	epsilon := 0
	for k := range input {
		sum := 0
		totalBits := len(input[k])
		bitNumber := (len(input) - k)-1
		for _, b := range input[k]{
			sum += b
			if sum > totalBits/2{
				gamma = gamma + (1 << bitNumber)
				break
			}
		}
		if sum <= totalBits/2{
			epsilon = epsilon + (1 << bitNumber)
		}
	}
	return gamma, epsilon
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	gamma, epsilon := calcPowerFactors(flipArray(parseInput(trimmedInput)))
	fmt.Println("part1", gamma* epsilon)
}
