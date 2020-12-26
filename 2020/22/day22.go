package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Player struct {
	deck []int
}

func (p Player) calcScore() int {

	result := 0
	for i := 0; i < len(p.deck); i++ {
		result += (i + 1) * (p.deck[len(p.deck)-i-1])
	}
	return result
}

var inputFile = flag.String("f", "testInput", "Relative file path to use as input.")

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	p1, p2 := parse(trimmedInput)
	for len(p1.deck) != 0 && len(p2.deck) != 0 {
		fmt.Printf("P1: %+v\n", p1)
		fmt.Printf("P2: %+v\n", p2)
		p1Card := p1.deck[0]
		p1.deck = p1.deck[1:]
		p2Card := p2.deck[0]
		p2.deck = p2.deck[1:]
		if p1Card > p2Card {
			p1.deck = append(p1.deck, p1Card, p2Card)
		} else if p2Card > p1Card {
			p2.deck = append(p2.deck, p2Card, p1Card)
		}
	}
	fmt.Printf("P1: %+v\n", p1)
	fmt.Printf("P2: %+v\n", p2)
	fmt.Println("P1 Score", p1.calcScore())
	fmt.Println("P2 Score", p2.calcScore())

}

func parse(s []string) (p1 Player, p2 Player) {

	var currentPlayer *Player
	currentPlayer = &p1
	for _, l := range s[1:] {
		if strings.Index(l, "Player ") == 0 {
			currentPlayer = &p2
		} else if l != "" {
			card, _ := strconv.Atoi(strings.TrimSpace(l))
			currentPlayer.deck = append(currentPlayer.deck, card)
		}
	}
	return p1, p2

}
