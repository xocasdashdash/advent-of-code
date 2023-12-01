package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string
var testMode = flag.Bool("test", false, "Set to run using the testInput")

func calcSum([][]string) 
func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		input = testInput
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println("Took", time.Since(start))

	sum := 0
	for _, l := range trimmedInput {
		firstInt := -1
		secondInt := -1
		for _, c := range strings.Split(l, "") {
			fmt.Println("c", c)
			intVal, err := strconv.Atoi(c)
			if err != nil {
				continue
			}
			if firstInt == -1 {
				firstInt = intVal
			} else {
				secondInt = intVal
			}
		}
		if secondInt == -1 {
			secondInt = firstInt
		}
		sum = sum + firstInt*10 + secondInt
	}
	fmt.Println("Sum", sum)

}
