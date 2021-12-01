package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func sumWindow(list []int, windowSize int, windowIndex int) int {
	sum :=0
	for i:=0;i<windowSize;i++{
		sum = sum + list[i+windowIndex]
	}

	return sum
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	integerInput := make([]int, len(trimmedInput))
	for k := range trimmedInput {
		tmp, _ := strconv.Atoi(trimmedInput[k])
		integerInput[k] = tmp
	}
	windowSizes := []int{1,3}
	for _, windowSize := range windowSizes {
		increases := 0
		for i:=windowSize;i<len(integerInput);i++{
			first := sumWindow(integerInput, windowSize, i-windowSize)
			second := sumWindow(integerInput, windowSize, i-windowSize+1)
			if second > first {
				increases = increases + 1
			}
		}
		fmt.Printf("%d\n", increases)
	}
}
