package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const preambleLength = 25

func main() {

	input, _ := ioutil.ReadFile("testInput")
	numbers := make([]int, 0, 100)
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		converted, _ := strconv.Atoi(s)
		numbers = append(numbers, converted)
	}

	for i := preambleLength; i < len(numbers); i++ {
		fmt.Printf("Testing %d \n", numbers[i])
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
			fmt.Printf("Bad number: %d\n", candidate)
			break
		}

	}
}
