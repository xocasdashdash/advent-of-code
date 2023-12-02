package main

import (
	_ "embed"
	"flag"
	"fmt"
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
	games := parseGames(trimmedInput)
	maxSet := Set{
		blue:  14,
		green: 13,
		red:   12,
	}
	validGames := validateGames(maxSet, games)
	fmt.Println("P1", addGames(validGames))
	fmt.Println("P2", part2(games))
}

func part2(games []Game) int {
	sum := 0
	for _, game := range games {
		minSet := Set{}
		for _, s := range game.sets {
			if s.blue > minSet.blue {
				minSet.blue = s.blue
			}
			if s.green > minSet.green {
				minSet.green = s.green
			}
			if s.red > minSet.red {
				minSet.red = s.red
			}
		}
		power := (minSet.blue * minSet.green * minSet.red)
		sum += power
	}
	return sum
}

type Set struct {
	blue  int
	red   int
	green int
}
type Game struct {
	sets []Set
	id   int
}

func addGames(games []Game) int {
	sum := 0
	for _, g := range games {
		sum += g.id
	}
	return sum
}
func validateGames(maxSet Set, games []Game) []Game {
	result := make([]Game, 0, len(games))
	for _, g := range games {
		valid := true
		for _, s := range g.sets {
			if s.blue > maxSet.blue {
				valid = false
				break
			}
			if s.green > maxSet.green {
				valid = false
				break
			}
			if s.red > maxSet.red {
				valid = false
				break
			}
		}
		if valid {
			result = append(result, g)
		}
	}
	return result
}

func parseGames(input []string) []Game {
	result := make([]Game, 0, len(input))
	for _, l := range input {
		if strings.TrimSpace(l) == "" {
			continue
		}
		g := Game{}
		parts := strings.Split(l, ":")
		// Get game id
		gameIdParts := strings.Split(strings.TrimSpace(parts[0]), " ")
		g.id, _ = strconv.Atoi(gameIdParts[1])
		g.sets = make([]Set, 0, 10)
		for _, splitSets := range strings.Split(strings.TrimSpace(parts[1]), ";") {
			s := Set{}
			setSplit := strings.Split(strings.TrimSpace(splitSets), ",")
			for _, setInfo := range setSplit {
				coloringParts := strings.Split(strings.TrimSpace(setInfo), " ")
				quantity, _ := strconv.Atoi(coloringParts[0])
				switch coloringParts[1] {
				case "red":
					s.red = quantity
				case "blue":
					s.blue = quantity
				case "green":
					s.green = quantity
				}
			}
			g.sets = append(g.sets, s)
		}
		result = append(result, g)
	}
	return result
}
