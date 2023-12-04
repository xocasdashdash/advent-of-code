package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string
var testMode = flag.Bool("test", false, "Set to run using the testInput")

func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		input = testInput
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println("Took", time.Since(start))

	cards, winningNumbers := parseInput(trimmedInput)
	totalPoints := calcTotalPoints(cards, winningNumbers)
	fmt.Println("Part 1", totalPoints)
}

func intersectNumbers(a []int, b []int) []int {
	result := make([]int, 0)

	bMap := make(map[int]bool, 0)
	for _, n := range b {
		bMap[n] = true
	}
	for _, n := range a {
		if _, ok := bMap[n]; ok {
			result = append(result, n)
		}
	}
	fmt.Println("intersection of ", a, b, "is", result)
	return result
}
func calcTotalPoints(cards map[int]Card, winningNumbers [][]int) int {
	result := 0
	for _, c := range cards {
		interSection := intersectNumbers(c.numbers, winningNumbers[c.index-1])
		if len(interSection) == 0 {
			continue
		}
		partialPoints := int(math.Pow(float64(2), float64(len(interSection)-1)))
		result += partialPoints
		fmt.Println("Card", c.index-1, "has this matches", interSection, partialPoints)
	}
	return result
}

type Card struct {
	index   int
	numbers []int
}

func parseNumbers(input string) []int {
	result := make([]int, 0)
	for _, cardNumbers := range strings.Split(strings.TrimSpace(input), " ") {
		candidateNumber := strings.TrimSpace(cardNumbers)
		if candidateNumber == "" {
			continue
		}
		cardNumber, err := strconv.Atoi(candidateNumber)
		if err != nil {
			panic(err)
		}
		result = append(result, cardNumber)
	}
	return result
}

func parseCard(input string) Card {
	result := Card{}
	parts := strings.Split(strings.TrimSpace(input), ":")

	index, err := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(parts[0], "Card", "")))
	if err != nil {
		panic("error parsing " + err.Error())
	}
	result.index = index
	result.numbers = parseNumbers(parts[1])

	return result

}
func parseInput(input []string) (map[int]Card, [][]int) {

	cards := make(map[int]Card, 0)
	winningNumbers := make([][]int, 0)
	for _, l := range input {
		parts := strings.Split(l, "|")
		card := parseCard(parts[0])
		numbers := parseNumbers(parts[1])
		winningNumbers = append(winningNumbers, numbers)
		cards[card.index-1] = card
	}
	return cards, winningNumbers
}
