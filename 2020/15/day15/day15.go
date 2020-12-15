package day15

type Game struct {
	StartNumbers []int
}

func (g *Game) PlayNRounds(rounds int) int {
	seen := make(map[int]int)
	var lastNumber int
	for i, n := range g.StartNumbers {
		turn := i + 1
		lastNumber = n
		seen[lastNumber] = turn
	}
	numberToSay := 0
	var j, n int
	var ok bool
	for j = len(g.StartNumbers); j < rounds; j++ {
		lastNumber = numberToSay
		if n, ok = seen[lastNumber]; ok {
			numberToSay = j + 1 - n
		} else {
			numberToSay = 0
		}
		seen[lastNumber] = j + 1
	}

	return lastNumber
}
