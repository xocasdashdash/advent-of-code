package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string
var testMode = flag.Bool("test", false, "Set to run using the testInput")

type Item struct {
	Inspected bool
	Worry     uint64
}
type Operation struct {
	Expression string
}

func (o *Operation) Apply(oldValue uint64) uint64 {

	tokens := strings.Split(o.Expression, " ")
	cursor := 3
	leftSide := oldValue
	op := tokens[cursor]
	cursor++
	rightSide := oldValue
	if tokens[cursor] != "old" {
		tmpRightSide, _ := strconv.Atoi(tokens[cursor])
		rightSide = uint64(tmpRightSide)
	}

	switch op {
	case "*":
		return leftSide * rightSide
	case "+":
		return leftSide + rightSide
	default:
		panic("bad value")
	}
}

type Monkey struct {
	_             struct{}
	Items         []*Item
	Index         int
	Divisible     uint64
	TrueIndex     int
	FalseIndex    int
	ActionTrue    *Monkey
	ActionFalse   *Monkey
	op            *Operation
	ActivityIndex int
}

func parseMonkey(input string) *Monkey {
	lines := strings.Split(input, "\n")
	cursor := 0

	var monkeyIndex int
	fmt.Sscanf(strings.TrimSpace(lines[cursor]), "Monkey %d:", &monkeyIndex)
	cursor++
	var itemString string
	itemString = strings.TrimSpace(strings.Split(lines[cursor], ":")[1])
	items := make([]*Item, 0, 1)
	for _, item := range strings.Split(itemString, ",") {
		worry, _ := strconv.Atoi(strings.TrimSpace(item))
		items = append(items, &Item{Worry: uint64(worry)})
	}
	cursor++
	var operationString string
	operationString = strings.TrimSpace(strings.Split(lines[cursor], ":")[1])
	operation := Operation{
		Expression: operationString,
	}
	cursor++
	var divisible int
	fmt.Sscanf(strings.TrimSpace(lines[cursor]), "Test: divisible by %d", &divisible)

	cursor++
	var trueIndex int
	fmt.Sscanf(strings.TrimSpace(lines[cursor]), "If true: throw to monkey %d", &trueIndex)

	cursor++
	var falseIndex int
	fmt.Sscanf(strings.TrimSpace(lines[cursor]), "If false: throw to monkey %d", &falseIndex)

	return &Monkey{
		Index:      monkeyIndex,
		op:         &operation,
		Divisible:  uint64(divisible),
		TrueIndex:  trueIndex,
		FalseIndex: falseIndex,
		Items:      items,
	}
}

func doTheMonkeyBusiness(monkeys []*Monkey, rounds int, relief bool) int {

	var divisor uint64
	divisor = 3

	if relief == false {
		divisor = 1
		// https://en.wikipedia.org/wiki/Chinese_remainder_theorem
		for _, m := range monkeys {
			divisor = divisor * m.Divisible
		}
	}
	for r := 1; r <= rounds; r++ {
		for _, monkey := range monkeys {
			for _, item := range monkey.Items {
				item.Worry = monkey.op.Apply(item.Worry)
				if !relief {
					item.Worry = item.Worry % divisor
				} else {
					item.Worry = item.Worry / divisor
				}
				targetMonkey := monkey.ActionTrue
				if item.Worry%monkey.Divisible != 0 {
					targetMonkey = monkey.ActionFalse
				}
				targetMonkey.Items = append(targetMonkey.Items, item)
				monkey.ActivityIndex++
			}
			monkey.Items = make([]*Item, 0)
		}
	}
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].ActivityIndex > monkeys[j].ActivityIndex
	})
	return monkeys[0].ActivityIndex * monkeys[1].ActivityIndex
}
func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		input = testInput
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n\n")
	monkeysP1 := make([]*Monkey, len(trimmedInput))
	monkeysP2 := make([]*Monkey, len(trimmedInput))
	for monkeyIndex, m := range trimmedInput {
		monkeysP1[monkeyIndex] = parseMonkey(m)
		monkeysP2[monkeyIndex] = parseMonkey(m)
	}
	for _, m := range monkeysP1 {
		m.ActionTrue = monkeysP1[m.TrueIndex]
		m.ActionFalse = monkeysP1[m.FalseIndex]
	}

	for _, m := range monkeysP2 {
		m.ActionTrue = monkeysP2[m.TrueIndex]
		m.ActionFalse = monkeysP2[m.FalseIndex]
	}

	fmt.Println("Part1", doTheMonkeyBusiness(monkeysP1, 20, true))
	fmt.Println("Part2", doTheMonkeyBusiness(monkeysP2, 10000, false))
	fmt.Println("Took", time.Since(start))
}
