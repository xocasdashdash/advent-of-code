package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

func calculatePriority(a string) int {
	lowercaseA := strings.ToLower(a)
	if a == lowercaseA {
		return int(a[0]) - int("a"[0]) + 1
	}
	return int(a[0]) - int("A"[0]) + 1 + 26

}

type Rucksack struct {
	First        []string
	Second       []string
	UniqueValues []string
}

func (r Rucksack) Repeated() int {
	firstUnique := make(map[string]struct{})
	for _, v := range r.First {
		firstUnique[v] = struct{}{}
	}
	for _, v := range r.Second {
		if _, ok := firstUnique[v]; ok {
			return calculatePriority(v)
		}
	}
	return -1
}

func (r *Rucksack) Unique() {
	unique := make(map[string]struct{})
	for _, v := range append(r.First, r.Second...) {
		unique[v] = struct{}{}
	}
	for k := range unique {
		r.UniqueValues = append(r.UniqueValues, k)
	}
}

func CommonElement(r []Rucksack) string {

	allUniques := make([]string, 0)
	for _, v := range r {
		allUniques = append(allUniques, v.UniqueValues...)
	}
	sort.Strings(allUniques)
	for k := range allUniques {
		found := true
		// Works for any number of rucksacks
		for i := 1; i < len(r); i++ {
			found = found && allUniques[k] == allUniques[k+i]
			if !found {
				break
			}
		}
		if found {
			return allUniques[k]
		}
	}
	return ""

}

func parseLine(l string) Rucksack {

	r := Rucksack{}
	r.First = strings.Split(l[:len(l)/2], "")
	r.Second = strings.Split(l[len(l)/2:], "")
	r.Unique()
	return r
}
func main() {
	t := time.Now()
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	part1 := 0
	fmt.Println("parsing took", time.Since(t))
	t = time.Now()
	rucksacks := make([]Rucksack, len(trimmedInput))
	for k, l := range trimmedInput {
		rucksacks[k] = parseLine(l)
		part1 += rucksacks[k].Repeated()
	}
	fmt.Println("part 1 took", time.Since(t))
	t = time.Now()
	part2 := 0
	for i := 0; i < len(rucksacks); i = i + 3 {
		sharedValue := CommonElement(rucksacks[i : i+3])
		part2 += calculatePriority(sharedValue)
	}
	fmt.Println("part 2 took", time.Since(t))
	fmt.Println("Part1", part1)
	fmt.Println("Part2", part2)
}
