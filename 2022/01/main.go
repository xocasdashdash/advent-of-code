package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

type Elf struct {
	Calories int
}
type Elves []Elf

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	army := make(Elves, 0, 0)

	var currentElf Elf
	for _, line := range trimmedInput {
		if line == "" {
			army = append(army, currentElf)
			currentElf = Elf{}
			continue
		}
		calories, e := strconv.Atoi(strings.TrimSpace(line))
		if e != nil {
			panic(e)
		}
		currentElf.Calories += calories
	}
	army = append(army, currentElf)
	sort.Slice(army, func(i, j int) bool {
		return army[i].Calories > army[j].Calories
	})

	fmt.Println("Max elf", army[0].Calories)
	// Part2
	topThree := army[0].Calories + army[1].Calories + army[2].Calories
	fmt.Println("Top three", topThree)

}
