package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	input := read()

	var id1, id2 string
	h := 0
	sort.Strings(input)
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}

loop:
	for i := range input {
		for j := i + 1; j < len(input); j++ {
			h = h + 1
			if differByOne(input[i], input[j]) {
				id1, id2 = input[i], input[j]
				break loop
			}
		}
	}
	fmt.Printf("%d", h)
	if id1 == "" {
		panic("not found")
	}

	s := strings.Builder{}
	for k := 0; k < len(id1); k++ {
		if id1[k] == id2[k] {
			s.WriteByte(id1[k])
		}
	}

	fmt.Printf("Answer: %s\n", id1)
}

func read() (input []string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		input = append(input, s.Text())
	}
	return input
}

func differByOne(a, b string) bool {
	diff := 0
	for k := 0; k < len(a); k++ {
		if a[k] != b[k] {
			diff++
		}
		if diff > 1 {
			return false
		}
	}
	return diff == 1
}
