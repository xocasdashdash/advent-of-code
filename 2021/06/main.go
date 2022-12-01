package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func parseInput(input []string)[]int {
	line := input[0]
	result := make([]int, 0,10)
	for _,v := range strings.Split(line,",") {
		temp, _ := strconv.Atoi(v)
		result = append(result, temp)
	}
	return result
}
func reproduce(fishes []int, days int) int {
	daysArray := make([]int,9,9)
	for _, v := range fishes {
		daysArray[v] = daysArray[v] + 1
	}
	for i:=0;i<days;i++ {
		nextGen := daysArray[0]
		copy(daysArray, daysArray[1:])
		daysArray[6] = daysArray[6] + nextGen
		daysArray[len(daysArray)-1] = nextGen
	}
	sum := 0
	for _, v := range daysArray {
		sum += v
	}
	return sum
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	fishes := parseInput(trimmedInput)
	t1 := time.Now()
	p1 := reproduce(fishes,80)
	took1 := time.Since(t1)
	p2 := reproduce(fishes, 256)
	took2 := time.Since(t1) - took1
	fmt.Println("Part1", p1)
	fmt.Println("Took", took1)
	fmt.Println("Part2", p2)
	fmt.Println("Took", took2)

}
