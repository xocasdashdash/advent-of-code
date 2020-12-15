package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var t = flag.Int("t", 2020, "Number of turns to play")
var debug = flag.Bool("d", false, "In order to enable debugging")

type Number struct {
	LastTwoTurns []int
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	numbers := parse(trimmedInput)
	if *debug {
		fmt.Printf("Input numbers: %+v\n", numbers)
	}
	//This is a map of the number and an array
	//of length 2. The first position holds the last
	//time the number was said out loud, and the second one holds
	//the number of times it was said.
	lastTimeSeenNumber := make(map[int]Number, *t)

	lastNumberSpoken := 0
	saidNumbers := make([]int, 0, *t)
	var saidNumber int
	for i := range numbers {
		turn := i + 1
		saidNumber = numbers[i]
		saidNumbers = append(saidNumbers, saidNumber)
		lastNumberSpoken = numbers[i]
		lastTimeSeenNumber[lastNumberSpoken] = Number{LastTwoTurns: []int{turn}}
	}
	fmt.Printf("Stating numbers: %+v\n%+v\n", saidNumbers, lastTimeSeenNumber)
	numberToSay := 0
	for j := len(numbers); j <= *t; j++ {
		turn := j + 1
		saidNumbers = append(saidNumbers, numberToSay)
		var n Number
		var ok bool
		if n, ok = lastTimeSeenNumber[numberToSay]; !ok {
			lastTimeSeenNumber[numberToSay] = Number{LastTwoTurns: []int{turn}}
			numberToSay = 0
		} else {
			n = Number{LastTwoTurns: []int{turn, n.LastTwoTurns[0]}}
			lastTimeSeenNumber[numberToSay] = n
			numberToSay = n.LastTwoTurns[0] - n.LastTwoTurns[1]
		}
	}
	//printNumbers(saidNumbers)
	//fmt.Printf("Said numbers: %+v. %d number: %d\n", saidNumbers, *t, saidNumbers[*t-1])
	fmt.Printf("Number at %d turn: %d\n", *t, saidNumbers[*t-1])

}
func printNumbers(n []int) {
	for k, v := range n {
		fmt.Printf("Turn: %d => %d\n", k+1, v)
	}
}
func parse(s []string) []int {
	var result []int
	for i := range s {
		for _, candidateNumber := range strings.Split(s[i], ",") {
			n, _ := strconv.Atoi(candidateNumber)
			result = append(result, n)

		}
	}
	return result

}
