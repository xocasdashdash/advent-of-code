package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	deck        []int
	seenDecks   map[string]bool
	currentGame int
	ID          int
}

func (p Player) String() string {

	return fmt.Sprintf("{ Player %d, CurrentGame: %d, Deck: %+v }", p.ID, p.currentGame, p.deck)
}
func (p Player) calcScore() int {

	result := 0
	for i := 0; i < len(p.deck); i++ {
		result += (i + 1) * (p.deck[len(p.deck)-i-1])
	}
	return result
}
func (p *Player) generateDeckSignature() string {
	var sb strings.Builder
	for _, card := range p.deck {
		sb.WriteString(fmt.Sprintf("%d", card))
	}
	return sb.String()
}

var inputFile = flag.String("f", "testInput2", "Relative file path to use as input.")

var globalGameCounter = 1

//PlayRecursiveCombat Plays recursive combat
func PlayRecursiveCombat(p1, p2 Player) (Player, Player) {
	round := 1
	currentGame := globalGameCounter
	// fmt.Printf("=== Game %d ===\n\n", currentGame)
	globalGameCounter++
	p1.currentGame = currentGame
	p2.currentGame = currentGame

	seenDecks := make(map[int]map[string]bool)
	seenDecks[p1.ID] = make(map[string]bool)
	seenDecks[p2.ID] = make(map[string]bool)
	for len(p1.deck) != 0 && len(p2.deck) != 0 {
		// fmt.Printf("-- Round %d (Game %d) --\n", round, currentGame)
		// fmt.Printf("P1: %+v\n", p1)
		// fmt.Printf("P2: %+v\n", p2)
		if seenDecks[p1.ID][p1.generateDeckSignature()] || seenDecks[p2.ID][p2.generateDeckSignature()] {
			// fmt.Printf("Repeated deck! on the same game. P1 wins the game")
			p2.deck = make([]int, 0)
			return p1, p2
		}
		seenDecks[p1.ID][p1.generateDeckSignature()] = true
		seenDecks[p2.ID][p2.generateDeckSignature()] = true
		p1Card := p1.deck[0]
		p1.deck = p1.deck[1:]
		p2Card := p2.deck[0]
		p2.deck = p2.deck[1:]
		if p1Card <= len(p1.deck) && p2Card <= len(p2.deck) {
			// fmt.Printf("Playing a sub-game to determine the winner...\n")
			p1Prime := p1
			p2Prime := p2
			p1Prime.deck = make([]int, p1Card, p1Card)
			copy(p1Prime.deck, p1.deck[:p1Card])
			// fmt.Printf("Copied %d to p1Prime deck\n", copied)
			p2Prime.deck = make([]int, p2Card, p2Card)
			copy(p2Prime.deck, p2.deck[:p2Card])
			// fmt.Printf("Copied %d to p2Prime deck\n", copied)
			p1Prime, _ = PlayRecursiveCombat(p1Prime, p2Prime)
			// fmt.Printf("Back to game %d\n", currentGame)
			if p1Prime.calcScore() > 0 {
				// fmt.Printf("P1 won the subgame appending %d,%d to %+v\n", p1Card, p2Card, p1.deck)
				p1.deck = append(p1.deck, p1Card, p2Card)
			} else {
				// fmt.Printf("P2 won the subgame appending %d,%d to %+v\n", p2Card, p1Card, p2.deck)
				p2.deck = append(p2.deck, p2Card, p1Card)
			}
		} else {
			if p1Card > p2Card {
				// fmt.Printf("P1 won the round appending %d,%d to %+v\n", p1Card, p2Card, p1.deck)
				p1.deck = append(p1.deck, p1Card, p2Card)
			} else if p2Card > p1Card {
				// fmt.Printf("P2 won the round appending %d,%d to %+v\n", p2Card, p1Card, p2.deck)
				p2.deck = append(p2.deck, p2Card, p1Card)
			}
		}
		// fmt.Printf("-- Round %d (Game %d) --\n", round, currentGame)
		round++
	}
	if p1.calcScore() == 0 {
		// fmt.Printf("P2 won the game!\n\n")
	} else {
		// fmt.Printf("P1 won the game!\n\n")
	}
	// fmt.Printf("P1: %+v\n", p1)
	// fmt.Printf("P2: %+v\n", p2)
	return p1, p2
}

func playNormalCombat(p1, p2 Player) {
	seenDecks := make(map[int]map[string]bool)
	seenDecks[p1.ID] = make(map[string]bool)
	seenDecks[p2.ID] = make(map[string]bool)
	for len(p1.deck) != 0 && len(p2.deck) != 0 {
		if seenDecks[p1.ID][p1.generateDeckSignature()] || seenDecks[p2.ID][p2.generateDeckSignature()] {
			return
		}
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
	part1 := p1.calcScore()
	if p1.calcScore() == 0 {
		part1 = p2.calcScore()
	}
	fmt.Printf("Part1: %d\n", part1)
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	var p1, p2 Player
	p1, p2 = parse(trimmedInput)
	playNormalCombat(p1, p2)
	p1, p2 = parse(trimmedInput)
	t := time.Now()
	p1, p2 = PlayRecursiveCombat(p1, p2)
	fmt.Printf("Recursive combat took %s\n", time.Since(t))
	part2 := p1.calcScore()
	if p1.calcScore() == 0 {
		fmt.Printf("P2 won!")
		part2 = p2.calcScore()
	} else {
		fmt.Printf("P1 won!")
	}
	fmt.Printf("Part2: %d\n", part2)

}

func parse(s []string) (p1 Player, p2 Player) {

	var currentPlayer *Player
	currentPlayer = &p1
	currentPlayer.ID = 1
	for _, l := range s[1:] {
		currentPlayer.seenDecks = make(map[string]bool)
		if strings.Index(l, "Player ") == 0 {
			currentPlayer = &p2
			currentPlayer.ID = 2
		} else if l != "" {
			card, _ := strconv.Atoi(strings.TrimSpace(l))
			currentPlayer.deck = append(currentPlayer.deck, card)
		}
	}
	return p1, p2

}
