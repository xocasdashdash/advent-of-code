package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Bag struct {
	Name     string
	Contains []*Bag
}

func (b *Bag) String() string {
	s := "|Name: %s|"
	c := make([]string, len(b.Contains), len(b.Contains))
	for k, v := range b.Contains {
		c[k] = v.Name
	}
	return fmt.Sprintf(s, b.Name) //, strings.Join(c, ","))
}

type ContainedBy map[string][]string

func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)
	containedBy := make(ContainedBy, 1000)
	for s.Scan() {
		l := s.Text()
		bags := strings.Split(l, "bags contain")
		bagName := strings.TrimSpace(bags[0])
		contains := strings.Split(bags[1], ",")
		for _, b := range contains {
			b = strings.TrimRight(b, ".")
			if b == "no other bags" {
				continue
			}
			bTokens := strings.Split(strings.TrimSpace(b), " ")
			bName := strings.TrimSpace(bTokens[1] + " " + bTokens[2])
			fmt.Printf("Bag Name: %s\n", bName)
			if _, ok := containedBy[bName]; !ok {
				containedBy[bName] = make([]string, 0, 100)
			}
			containedBy[bName] = append(containedBy[bName], bagName)
		}
	}
	fmt.Printf("Contained By %+v", containedBy)
	seen := make(map[string]bool, len(containedBy))
	var visitor func(n string) int
	visitor = func(n string) int {
		if seen[n] {
			return 0
		}
		seen[n] = true
		total := 1
		if contained, ok := containedBy[n]; ok {
			for _, v := range contained {
				total = total + visitor(v)
				fmt.Printf("Visited %s. Total %d\n", v, total)
			}
		}
		return total
	}
	fmt.Printf("\n%s\n", containedBy["bright white"])
	t := visitor("shiny gold") - 1
	fmt.Printf("\nTotal %d", t)
}
