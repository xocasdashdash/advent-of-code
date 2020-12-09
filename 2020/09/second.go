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
	targetNumber := make([]int, 0, 100)
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
			targetNumber = append(targetNumber, candidate)
			fmt.Printf("Bad number: %d\n", candidate)
		}
	}

	fmt.Printf("Found %d candidates: %+v\n", len(targetNumber), targetNumber)
	startRange := 0
	endRange := 1

	firstBadNumber := targetNumber[0]
	totalSums := 0
	for {
		currentSum := 0
		minNumber := firstBadNumber + 1
		maxNumber := 0
		totalSums++
		for i := startRange; i < endRange; i++ {
			currentSum += numbers[i]
			if numbers[i] > maxNumber {
				maxNumber = numbers[i]
			}
			if numbers[i] < minNumber {
				minNumber = numbers[i]
			}
		}
		if currentSum == firstBadNumber {
			fmt.Printf("Found range! %d, %d after %d sums\n", startRange, endRange, totalSums)
			fmt.Printf("Min in range :%d, max in range:%d, sum :%d\n", minNumber, maxNumber, minNumber+maxNumber)
			break
		}
		if currentSum < firstBadNumber {
			endRange++
		} else if currentSum > firstBadNumber {
			startRange++
		}
	}
}
