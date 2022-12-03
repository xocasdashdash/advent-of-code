package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

const (
	Lose int = iota
	Draw
	Win
)

type Result struct {
	Points int
}
type Move interface {
	Compare(i Move) int
	Name() string
	Points() int
}

type Paper struct{}

func (p *Paper) Compare(i Move) int {
	switch i.(type) {
	case *Rock:
		return Win
	case *Scissors:
		return Lose
	default:
		return Draw
	}
}
func (p *Paper) Points() int {
	return 2
}
func (p *Paper) Name() string { return "paper" }

type Scissors struct{}

func (p *Scissors) Compare(i Move) int {
	switch i.(type) {
	case *Rock:
		return Lose
	case *Paper:
		return Win
	default:
		return Draw
	}
}
func (s *Scissors) Points() int {
	return 3
}
func (p *Scissors) Name() string { return "scissors" }

type Rock struct{}

func (r *Rock) Compare(i Move) int {

	switch i.(type) {
	case *Paper:
		return Lose
	case *Scissors:
		return Win
	default:
		return Draw
	}
}
func (p *Rock) Name() string { return "rock" }

var MoveMap map[string]Move = map[string]Move{
	"A": &Rock{},
	"B": &Paper{},
	"C": &Scissors{},
	"X": &Rock{},
	"Y": &Paper{},
	"Z": &Scissors{},
}

func (r *Rock) Points() int {
	return 1
}

type Round struct {
	Move        Move
	CounterMove Move
	Result      int
}

var CounterMoveMap map[string]Move = map[string]Move{
	"rock":     &Paper{},
	"paper":    &Scissors{},
	"scissors": &Rock{},
}

func (r *Round) CalculateCounterMove() Move {

	switch r.Result {
	case Win:
		return CounterMoveMap[r.Move.Name()]
	case Lose:
		return CounterMoveMap[CounterMoveMap[r.Move.Name()].Name()]
	default:
		return r.Move
	}

}

var ResultMap map[string]int = map[string]int{
	"X": Lose,
	"Y": Draw,
	"Z": Win,
}

func parseLine(l string, part2 bool) Round {
	parts := strings.Split(l, " ")
	if part2 == true {
		r := Round{
			Move:   MoveMap[parts[0]],
			Result: ResultMap[parts[1]],
		}
		r.CounterMove = r.CalculateCounterMove()
		return r
	}
	return Round{
		Move:        MoveMap[parts[1]],
		CounterMove: MoveMap[parts[0]],
		Result:      MoveMap[parts[1]].Compare(MoveMap[parts[0]]),
	}
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")

	rounds := make([]Round, len(trimmedInput), len(trimmedInput))
	roundsP2 := make([]Round, len(trimmedInput), len(trimmedInput))
	pointsP1 := 0
	pointsP2 := 0
	for k, l := range trimmedInput {
		rounds[k] = parseLine(l, false)
		roundsP2[k] = parseLine(l, true)
		pointsP1 += 3*rounds[k].Result + rounds[k].Move.Points()
		pointsP2 += 3*roundsP2[k].Result + roundsP2[k].CounterMove.Points()
		fmt.Println("PointsP2", pointsP2)
		fmt.Printf("Round %#v\n", roundsP2[k])

	}
	fmt.Println("Part1", pointsP1)
	fmt.Println("Part2", pointsP2)

}
