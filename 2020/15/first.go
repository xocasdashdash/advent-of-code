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
var t = flag.Int("t", 2020, "Number of turns to play")
var debug = flag.Bool("d", false, "In order to enable debugging")

type Number struct {
	A int
	B int
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	fmt.Printf("INput: %s\n", trimmedInput)
	start := time.Now()
	numbers := parse(trimmedInput)
	fmt.Printf("Time to parse: %s\n", time.Since(start))
	if *debug {
		fmt.Printf("Input numbers: %+v\n", numbers)
	}
	//This is a map of the number and an array
	//of length 2. The first position holds the last
	//time the number was said out loud, and the second one holds
	//the number of times it was said.
	lastTimeSeenNumber := make(map[int]Number, *t)

	lastNumberSpoken := 0
	saidNumbers := make([]int, *t+1, *t+1)
	var saidNumber int
	start = time.Now()
	for i := range numbers {
		turn := i + 1
		saidNumber = numbers[i]
		saidNumbers[i] = saidNumber
		lastNumberSpoken = numbers[i]
		lastTimeSeenNumber[lastNumberSpoken] = Number{A: turn}
	}
	if *debug {
		fmt.Printf("Starting numbers: %+v\n%+v\n", saidNumbers, lastTimeSeenNumber)
	}
	numberToSay := 0
	currentNumber := 0
	lastNumber := 0
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Printf("Speed: %d numbers/second\n", currentNumber-lastNumber)
			lastNumber = currentNumber
		}
	}()
	t2 := time.Now()
	for i := 0; i <= *t; i++ {
		currentNumber++
	}
	fmt.Printf("Took %s to loop through %d numbers\n", time.Since(t2), *t)
	lastNumber = 0
	currentNumber = 0

	for j := len(numbers); j <= *t; j++ {
		currentNumber++
		turn := j + 1
		saidNumbers[j] = numberToSay
		var n Number
		var ok bool
		if n, ok = lastTimeSeenNumber[numberToSay]; !ok {
			lastTimeSeenNumber[numberToSay] = Number{A: turn}
			numberToSay = 0
		} else {
			oldTurn := n.A
			n.A = turn
			lastTimeSeenNumber[numberToSay] = n
			numberToSay = turn - oldTurn
		}
	}
	fmt.Printf("Time to Finish: %s\n", time.Since(start))
	//printNumbers(saidNumbers)
	//fmt.Printf("Said numbers: %+v. %d number: %d\n", saidNumbers, *t, saidNumbers[*t-1])
	fmt.Printf("Number at %d: %d\n", *t, saidNumbers[*t-1])

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
