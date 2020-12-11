package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Seat byte
type SeatRow []Seat
type SeatMap []SeatRow
type Placing struct {
	Row    int
	Column int
	State  Seat
}

const Occupied = '#'
const Empty = 'L'
const Floor = '.'

func main() {

	input, _ := ioutil.ReadFile("input")
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	m := parseMap(trimmedInput)
	printMap(m)
	fmt.Printf("\n***************************\n")
	checksum := checkSumMap(m)
	for {
		m = evolve(m)
		newCheckSum := checkSumMap(m)
		if checksum == newCheckSum {
			break
		}
		checksum = newCheckSum
		//printMap(m)

	}
	printMap(m)
	fmt.Printf("\n***************************\n")
	//e, o, f := candidateNeighbors(m, 0, 3)
	//fmt.Printf("R:0,c:3\n")
	//fmt.Printf("Empty: %d, Occupied: %d, floor: %d\n", len(e), len(o), len(f))
	//fmt.Printf("E: %+v\n", e)
	//fmt.Printf("O: %+v\n", o)
	//fmt.Printf("F: %+v\n", f)

	occupied := 0
	for r := range m {
		for c := range m[r] {
			if m[r][c] == Occupied {
				occupied++
			}
		}
	}
	fmt.Printf("Part1: %d", occupied)
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
			_, occupiedPlacings, _ := candidateNeighbors(m, r, c)
			state := m[r][c]
			if state == Empty && len(occupiedPlacings) == 0 {
				state = Occupied
			} else if state == Occupied && len(occupiedPlacings) >= 4 {
				state = Empty
			}
			row = append(row, state)
		}
		evolvedSM = append(evolvedSM, row)
	}
	return evolvedSM
}

func candidateNeighbors(m SeatMap, row, column int) ([]Placing, []Placing, []Placing) {

	empty := make([]Placing, 0, 8)
	occupied := make([]Placing, 0, 8)
	floor := make([]Placing, 0, 8)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			candidateRow := row + i
			candidateColumn := column + j
			if candidateColumn == column && candidateRow == row {
				continue
			}
			if isValidNeighbor(m, candidateRow, candidateColumn) {
				s := m[candidateRow][candidateColumn]
				if s == Empty {
					empty = append(empty, Placing{
						Row:    candidateRow,
						Column: candidateColumn,
						State:  s,
					})
				} else if s == Occupied {
					occupied = append(occupied, Placing{
						Row:    candidateRow,
						Column: candidateColumn,
						State:  s,
					})
				} else {
					floor = append(floor, Placing{
						Row:    candidateRow,
						Column: candidateColumn,
						State:  s,
					})
				}
			}
		}
	}
	return empty, occupied, floor

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
