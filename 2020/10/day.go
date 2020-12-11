package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func main() {

	inputs := []string{"input"}
	for k := range inputs {
		input, _ := ioutil.ReadFile(inputs[k])
		trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
		fmt.Printf("Part 1: %d\n", calculate1(trimmedInput))
		fmt.Printf("Part 2: %d\n", calculate2(inputs[k], trimmedInput))
	}
}

func calculate1(i []string) int {
	adapters := make([]int, 0, 100)

	for _, s := range i {
		converted, _ := strconv.Atoi(s)
		adapters = append(adapters, converted)
	}
	sort.Ints(adapters)
	differences := make(map[int]int)

	currentJolt := 0

	for k := range adapters {
		currentAdapter := adapters[k]
		difference := currentAdapter - currentJolt
		differences[difference] = differences[difference] + 1
		currentJolt = currentAdapter
	}
	differences[3]++
	return differences[3] * differences[1]
}

func calculate2(name string, i []string) int {
	adapters := make([]int, 0, 100)
	exists := make(map[int]int)

	for _, s := range i {
		converted, _ := strconv.Atoi(s)
		adapters = append(adapters, converted)
	}
	sort.Ints(adapters)
	for k := range adapters {
		exists[adapters[k]] = k
	}
	diffs := []int{1, 2, 3}
	ways := make([]int, len(adapters))
	ways[len(adapters)-1] = 1
	//fmt.Printf("Adapters: %d -  %+v\n", len(adapters), adapters)
	//fmt.Printf("Ways: %d - %+v\n", len(ways), ways)
	for k := len(adapters) - 2; k >= 0; k-- {
		sum := 0
		for _, d := range diffs {
			doesTheNumberExists := adapters[k] + d
			if pos, ok := exists[doesTheNumberExists]; ok {
				sum += ways[pos]
			}
		}
		ways[k] = sum
	}
	ret := 0
	for v := 1; v <= 3; v++ {
		if pos, ok := exists[v]; ok {
			ret += ways[pos]
		}
	}
	//fmt.Printf("Ways: %#v\n", ways)
	//fmt.Printf("Exists: %#v\n", exists)
	//fmt.Printf("Return: %+v\n", ret)
	return ret
}
