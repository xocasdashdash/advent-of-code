package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Seat byte
type SeatRow []Seat
type SeatMap []SeatRow
type Coord struct {
	Row    int
	Column int
}

const Occupied = '#'
const Empty = 'L'
const Floor = '.'

var neighborCache = make(map[Coord][]Coord)
var checks int = 0

func main() {

	input, _ := ioutil.ReadFile("input")
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	m := parseMap(trimmedInput)
	fmt.Printf("Original map\n***************************\n")
	printMap(m)
	fmt.Printf("\n***************************\n")
	checksum := checkSumMap(m)
	for {
		m = evolve(m)
		newCheckSum := checkSumMap(m)
		if checksum == newCheckSum {
			fmt.Printf("\n*Stabilized*\n")
			break
		}
		checksum = newCheckSum
	}
	fmt.Printf("\nFinal map \n")
	printMap(m)
	fmt.Printf("Total checks: %d", checks)
	fmt.Printf("\n***************************\n")

	occupied := 0
	for r := range m {
		for c := range m[r] {
			if m[r][c] == Occupied {
				occupied++
			}
		}
	}
	fmt.Printf("Part2: %d", occupied)
}
func printMap(m SeatMap) {
	for r := range m {
		for c := range m[r] {
			fmt.Printf("%c", m[r][c])
		}
		fmt.Printf("\n")
	}
}
func checkSumMap(m SeatMap) string {
	result := make([]byte, len(m)*len(m[0]), len(m)*len(m[0]))
	for r := range m {
		for c := range m[r] {
			result[r*len(m[0])+c] = byte(m[r][c])
		}
	}
	return string(result)
}
func evolve(m SeatMap) SeatMap {
	evolvedSM := make(SeatMap, 0, len(m))
	for r := range m {
		row := make(SeatRow, 0, len(m[r]))
		for c := range m[r] {
			neighbors := candidateNeighbors(m, r, c)
			occupied := 0
			for _, n := range neighbors {
				if m[n.Row][n.Column] == Occupied {
					occupied++
				}
			}
			state := m[r][c]
			if state == Empty && occupied == 0 {
				state = Occupied
			} else if state == Occupied && occupied >= 5 {
				state = Empty
			}
			row = append(row, state)
		}
		evolvedSM = append(evolvedSM, row)
	}
	return evolvedSM
}

func candidateNeighbors(m SeatMap, row, column int) []Coord {

	var result []Coord
	type vector struct {
		i int
		j int
	}
	if neighborCache[Coord{row, column}] != nil {
		return neighborCache[Coord{row, column}]
	}
	vectors := make([]vector, 0, 8)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			v := vector{i, j}
			vectors = append(vectors, v)
		}
	}
	multiplier := 1
	for {
		if len(vectors) == 0 {
			break
		}
		newVectors := make([]vector, 0, len(vectors))
		for _, v := range vectors {
			candidateRow := row + v.i*multiplier
			candidateColumn := column + v.j*multiplier
			checks++
			if isValidNeighbor(m, candidateRow, candidateColumn) {
				s := m[candidateRow][candidateColumn]
				if s == Floor {
					newVectors = append(newVectors, v)
				} else {
					result = append(result, Coord{candidateRow, candidateColumn})
				}
			}
		}
		vectors = newVectors
		multiplier++
	}
	neighborCache[Coord{row, column}] = result
	return result

}
func isValidNeighbor(m SeatMap, row, column int) bool {
	if row < 0 {
		return false
	}
	if column < 0 {
		return false
	}
	if row >= len(m) {
		return false
	}
	if column >= len(m[row]) {
		return false
	}
	return true
}
func parseMap(m []string) SeatMap {
	sm := make(SeatMap, len(m))
	for r := range m {
		row := make(SeatRow, len(m[r]))
		for c := range m[r] {
			s := Seat(m[r][c])
			row[c] = s
		}
		sm[r] = row
	}
	return sm
}
