package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const preambleLength = 25

func main() {

	input, _ := ioutil.ReadFile("input")
	numbers := make([]int, 0, 100)
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		converted, _ := strconv.Atoi(s)
		numbers = append(numbers, converted)
	}
	targetNumber := 0
	for i := preambleLength; i < len(numbers); i++ {
		preamble := make(map[int]bool, preambleLength)
		found := false
		for k := i - preambleLength; k < i; k++ {
			preamble[numbers[k]] = true
		}
		candidate := numbers[i]
		for j := range preamble {
			search := candidate - j
			if preamble[search] {
				found = true
				break
			}
		}

		if !found {
			targetNumber = candidate
			break
		}
	}

	fmt.Printf("Found candidates: %+v\n", targetNumber)
	startRange := 0
	endRange := 0
	totalSums := 0
	var currentSum int
	for {
		totalSums++
		if currentSum == targetNumber {
			fmt.Printf("Found range! %d, %d after %d sums\n", startRange, endRange, totalSums)
			break
		} else if currentSum < targetNumber {
			currentSum += numbers[endRange]
			endRange++
		} else if currentSum > targetNumber {
			currentSum -= numbers[startRange]
			startRange++
		}
	}
	maxNumber := 0
	minNumber := currentSum + 1
	for i := startRange; i < endRange; i++ {
		if numbers[i] > maxNumber {
			maxNumber = numbers[i]
		}
		if numbers[i] < minNumber {
			minNumber = numbers[i]
		}
	}
	fmt.Printf("Min in range(%d-%d) :%d, max in range:%d, sum :%d\n", startRange, endRange, minNumber, maxNumber, minNumber+maxNumber)

}
