package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bag struct {
	Name     string
	Contains map[string]int
}

func (b *Bag) String() string {
	s := "|Name: %s|"
	return fmt.Sprintf(s, b.Name)
}

type BagList map[string]*Bag

func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)
	bagList := make(BagList, 1000)
	for s.Scan() {
		l := s.Text()
		bags := strings.Split(l, "bags contain")
		bagName := strings.TrimSpace(bags[0])
		contains := strings.Split(bags[1], ",")
		var bag *Bag
		var ok bool
		if bag, ok = bagList[bagName]; !ok {
			bag = &Bag{
				Name:     bagName,
				Contains: make(map[string]int, 1000),
			}
			bagList[bag.Name] = bag
		}
		for _, b := range contains {
			b = strings.TrimSpace(strings.TrimRight(b, "."))
			if b == "no other bags" {
				continue
			}
			bTokens := strings.Split(strings.TrimSpace(b), " ")
			bName := strings.TrimSpace(bTokens[1] + " " + bTokens[2])
			bTotal, _ := strconv.Atoi(bTokens[0])
			bagList[bag.Name].Contains[bName] = bTotal
		}
	}
	var visitor func(string) int
	visitor = func(n string) int {
		total := 1
		defer fmt.Printf("Visited %s. Total %d\n", n, total)
		if bag, ok := bagList[n]; ok {
			for name, number := range bag.Contains {
				total = total + (number * visitor(name))
			}
		}
		return total
	}
	t := visitor("shiny gold") - 1
	fmt.Printf("\nTotal %d", t)
}
