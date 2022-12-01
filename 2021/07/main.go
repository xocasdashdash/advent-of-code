package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func calcMedian(input []int) int{
	sort.Ints(input)
	return input[len(input)/2]
}
func calcMean(input []int) float64 {
	s := 0.0
	for _, v := range input {
		s+=float64(v)
	}	
	return s/float64(len(input))
}
func findNeareast (mean float64) ([]int) {

	return []int{
		int(math.Floor(mean)),int(math.Ceil(mean)),
	}
}
func parseInput(input []string)[]int {
	line := input[0]
	result := make([]int,0,strings.Count(line,","))
	for _, v := range strings.Split(line,",") {
		tmp, _ := strconv.Atoi(v)
		result = append(result, tmp)
	}
	return result
}
func calculateFuelP1(input []int, median int) int {
	result := 0
	for _, v := range input {
		if v > median {
			result = result + v - median
		}else {
			result = result + median -v
		}
	}
	return result
}

func calculateFuelP2(input []int, candidate int) int {
	total := 0
	for _, v := range input {
		if candidate == v {
			continue
		}
		diff := candidate - v
		if diff < 0{
			diff = -diff
		}
		total = total + ((diff * (diff + 1 ))/2)
	}
	return total
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	parsedInput := parseInput(trimmedInput)
	median := calcMedian(parsedInput)
	fmt.Println("median", median)
	part1 := calculateFuelP1(parsedInput, median)
	fmt.Println("part1", part1)
	mean := calcMean(parsedInput)
	candidates := findNeareast(
		mean,
	)
	fmt.Println("candidates", candidates)
	fmt.Println("part2 - c1", candidates[0], calculateFuelP2(parseInput(trimmedInput), candidates[0]))
	fmt.Println("part2 - c2", candidates[1], calculateFuelP2(parseInput(trimmedInput), candidates[1]))

}
