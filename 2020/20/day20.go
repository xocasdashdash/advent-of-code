package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "testInput", "Relative file path to use as input.")

var complementaryBorders = map[string]string{
	"top":    "bottom",
	"left":   "right",
	"right":  "left",
	"bottom": "top",
}
var borders = []string{
	"top", "right", "bottom", "left",
}
var flipped = []bool{
	false, true,
}

type Bit struct {
	Row int
	Col int
}

type TileConfiguration struct {
	rotations int
	flip      bool
}
type NeighborMatch struct {
	ID         int
	neighborTC TileConfiguration
}

func (nc NeighborMatch) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("{ ID: %d, TC: { rotations: %d, flipped %+v }", nc.ID, nc.neighborTC.rotations, nc.neighborTC.flip))
	return sb.String()
}

type Tile struct {
	ID                    int
	image                 map[Bit]bool
	neighborIDs           []int
	neighborConfiguration map[string]int
	locked                bool
	rotations             int
	flipped               bool
	sideSignatures        map[string]map[bool]int
}

func (t Tile) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tile %d:\n", t.ID))
	for row := 0; row <= 9; row++ {
		for col := 0; col <= 9; col++ {
			if t.image[Bit{Row: row, Col: col}] {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *Tile) rotateTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}
	t.rotations = (t.rotations + 1) % 4
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := 0; r <= 9; r++ {
		for c := 0; c <= 9; c++ {
			key := Bit{Row: r, Col: c}
			currentBit := t.image[key]
			nextKey := Bit{
				Row: 9 - c,
				Col: r,
			}
			if currentBit {
				newBitMap[nextKey] = true
			} else {
				newBitMap[nextKey] = false
			}
		}
	}
	t.image = newBitMap
	return nil
}

var checks = 0

func (t *Tile) compareSideToTile(tt Tile, border string) bool {
	var leftPoint, rightPoint Bit
	var rowIncreaseL, colIncreaseL int
	var rowIncreaseR, colIncreaseR int

	switch complementaryBorders[border] {
	case "top":
		leftPoint = Bit{Row: 0, Col: 0}
		rightPoint = Bit{Row: 9, Col: 0}
		colIncreaseL = 1
		colIncreaseR = 1
		rowIncreaseL = 0
		rowIncreaseR = 0
	case "bottom":
		leftPoint = Bit{Row: 9, Col: 0}
		rightPoint = Bit{Row: 0, Col: 0}
		colIncreaseL = 1
		colIncreaseR = 1
		rowIncreaseL = 0
		rowIncreaseR = 0
	case "right":
		leftPoint = Bit{Row: 0, Col: 9}
		rightPoint = Bit{Row: 0, Col: 0}
		colIncreaseL = 0
		colIncreaseR = 0
		rowIncreaseL = 1
		rowIncreaseR = 1
	case "left":
		leftPoint = Bit{Row: 0, Col: 0}
		rightPoint = Bit{Row: 0, Col: 9}
		colIncreaseL = 0
		colIncreaseR = 0
		rowIncreaseL = 1
		rowIncreaseR = 1
	}
	//fmt.Printf("Comparing side: %s\n", complementaryBorders[border])
	for p := 0; p <= 9; p++ {
		if t.image[leftPoint] != tt.image[rightPoint] {
			checks++
			//fmt.Printf("Tile %d is different at (%d,%d) from tile %d at (%d,%d)\n", t.ID, leftPoint.Row, leftPoint.Col, tt.ID, rightPoint.Row, rightPoint.Col)
			return false
		}
		leftPoint.Col += colIncreaseL
		rightPoint.Col += colIncreaseR
		leftPoint.Row += rowIncreaseL
		rightPoint.Row += rowIncreaseR
	}
	//fmt.Printf("Side %s matches between tile %d and %d\n", complementaryBorders[border], t.ID, tt.ID)
	//fmt.Printf("SideA: |%s|\nSideB: |%s|\n", t.printSide(complementaryBorders[border]), tt.printSide(border))
	return true
}
func (t *Tile) printSides() {

	for range borders {
		for range flipped {
			//fmt.Printf("Border %s,flipped: %+v,  result: %+d\n", b, f, t.getSideSignature(b))
			t.flipTile()
		}
		//t.rotateTile()
	}
}
func (t *Tile) flipTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}
	t.flipped = !t.flipped
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := 0; r <= 9; r++ {
		for c := 0; c <= 9; c++ {
			key := Bit{Row: r, Col: c}
			currentBit := t.image[key]
			nextKey := Bit{
				Row: 9 - r,
				Col: c,
			}
			if currentBit {
				newBitMap[nextKey] = true
			} else {
				newBitMap[nextKey] = false
			}
		}
	}
	t.image = newBitMap
	return nil
}

type borderFunc func() int

func (t *Tile) topNeighborFunc() int {

	result := 0
	for col := 0; col <= 9; col++ {
		keyToCheck := Bit{Row: 0, Col: col}
		if t.image[keyToCheck] {
			result += 1 << col
		}
	}
	return result
}

func (t *Tile) bottomNeighborFunc() int {

	result := 0
	for col := 0; col <= 9; col++ {
		keyToCheck := Bit{Row: 9, Col: col}
		if t.image[keyToCheck] {
			result += 1 << col
		}
	}
	return result
}

func (t *Tile) rightNeighborFunc() int {
	result := 0
	for row := 0; row <= 9; row++ {
		keyToCheck := Bit{Row: row, Col: 9}
		if t.image[keyToCheck] {
			result += 1 << row
		}
	}
	return result
}
func (t *Tile) leftNeighborFunc() int {
	result := 0
	for row := 0; row <= 9; row++ {
		keyToCheck := Bit{Row: row, Col: 0}
		if t.image[keyToCheck] {
			result += 1 << row
		}
	}
	return result
}
func (t *Tile) getAllPossibleSides() []int {

	var result []int
	for j := 0; j < 2; j++ {
		for range borders {
			result = append(result, t.topNeighborFunc())
			t.rotateTile()
		}
		t.flipTile()
	}
	sort.Ints(result)
	return result

}

func (t *Tile) matchBorder(tt Tile, border string) error {

	var borderResult int
	if t.locked {
		if t.compareSideToTile(tt, border) {
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("No border matches at current configuration rot: %d, f: %+v\n", t.rotations, t.flipped))
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if t.compareSideToTile(tt, border) {
				return nil
			}
			t.flipTile()
		}
		t.rotateTile()
	}
	return fmt.Errorf(fmt.Sprintf("No border matches. Border result %d", borderResult))
}

func (t *Tile) printSide(side string) string {
	var result int
	switch side {
	case "top":
		result = t.topNeighborFunc()
	case "bottom":
		result = t.bottomNeighborFunc()
	case "left":
		result = t.leftNeighborFunc()
	case "right":
		result = t.rightNeighborFunc()
	}

	r := strconv.FormatInt(int64(result), 2)

	r = strings.ReplaceAll(r, "1", "#")
	r = strings.ReplaceAll(r, "0", ".")
	currentLength := len(r)
	padding := 10 - currentLength
	return strings.Repeat(".", padding) + r
}

// func (t *Tile) updateSides() {

// 	t.sideSignatures["top"][t.flipped] = t.topNeighborFunc()
// 	t.sideSignatures["right"][t.flipped] = t.rightNeighborFunc()
// 	t.sideSignatures["bottom"][t.flipped] = t.bottomNeighborFunc()
// 	t.sideSignatures["left"][t.flipped] = t.leftNeighborFunc()
// 	t.flipTile()
// 	t.sideSignatures["top"][t.flipped] = t.topNeighborFunc()
// 	t.sideSignatures["right"][t.flipped] = t.rightNeighborFunc()
// 	t.sideSignatures["bottom"][t.flipped] = t.bottomNeighborFunc()
// 	t.sideSignatures["left"][t.flipped] = t.leftNeighborFunc()
// 	t.flipTile()
// }
func (t *Tile) getSideSignature(side string) int {

	var result int
	switch side {
	case "top":
		result = t.topNeighborFunc()
		return result
	case "bottom":
		result = t.bottomNeighborFunc()
		return result
	case "left":
		result = t.leftNeighborFunc()
		return result
	case "right":
		result = t.rightNeighborFunc()
		return result
	}
	panic("impossible")

}
func uniqueInts(input []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	sort.Ints(list)
	return list
}

func buildTilemap(startTile int, tilemap map[int]Tile) map[int]Tile {

	visitedTiles := make(map[int]bool)
	visitList := []int{startTile}

	startingTile := tilemap[startTile]
	tilemap[startTile] = startingTile
	for len(visitedTiles) != len(tilemap) {
		//fmt.Printf("Next visit: %+v\n", visitList)
		var nextVisitList []int
		if len(visitList) == 0 {
			//fmt.Printf("%+v", visitedTiles)
			for tID := range tilemap {
				if _, ok := visitedTiles[tID]; !ok {
					fmt.Printf("Did not visit %d\n", tID)
				}
			}
			panic("Bad")
		}
		for _, tileID := range visitList {
			//fmt.Printf("Visiting %d\n", tileID)
			if visitedTiles[tileID] {
				//fmt.Printf("Already seen %d\n", tileID)
				continue
			}
			currentTile := tilemap[tileID]
			visitedTiles[tileID] = true
			if len(currentTile.neighborConfiguration) == len(currentTile.neighborIDs) {
				continue
			}

			//currentTile.printSides()
			for _, neighborID := range currentTile.neighborIDs {
				if visitedTiles[neighborID] {
					//fmt.Printf("Skipping %d as we already visited it\n", neighborID)
					continue
				}
				currentNeighbor := tilemap[neighborID]
				if len(currentNeighbor.neighborConfiguration) == len(currentNeighbor.neighborIDs) {
					continue
				}
				nextVisitList = append(nextVisitList, neighborID)
				var matchError error
				for _, border := range borders {
					if currentNeighbor.neighborConfiguration != nil {
						_, ok := currentNeighbor.neighborConfiguration[complementaryBorders[border]]
						if ok {
							continue
						}
					}

					//This tile already has a neighbor here
					if currentTile.neighborConfiguration != nil {
						_, ok := currentTile.neighborConfiguration[border]
						if ok {
							continue
						}
					}
					borderToMatch := currentTile.getSideSignature(border)
					matchError = currentNeighbor.matchBorder(currentTile, border)
					if matchError == nil {
						candidateMatch := currentNeighbor.getSideSignature(complementaryBorders[border])
						if candidateMatch != borderToMatch {
							fmt.Printf("************\n")
							fmt.Printf("Error matching %s(%d) to %s(%d)\n", border, currentTile.ID, complementaryBorders[border], currentNeighbor.ID)
							fmt.Printf("SideA:\t%s\n", currentTile.printSide(border))
							fmt.Printf("SideB:\t%s\n", currentNeighbor.printSide(complementaryBorders[border]))
							fmt.Printf("SideAComp:\t%s\n", currentTile.printSide(complementaryBorders[border]))
							fmt.Printf("SideBComp:\t%s\n", currentNeighbor.printSide(border))
							fmt.Printf("%+v\n", currentTile)
							fmt.Printf("%+v\n", currentNeighbor)
							fmt.Printf("************\n")
						}
						// fmt.Printf("Current neighbor config: %+v\n", currentNeighbor.TileConfig)
						//We have a match
						if currentTile.neighborConfiguration == nil {
							currentTile.neighborConfiguration = make(map[string]int)
						}
						currentTile.neighborConfiguration[border] = currentNeighbor.ID
						currentTile.locked = true
						tilemap[currentTile.ID] = currentTile

						if currentNeighbor.neighborConfiguration == nil {
							currentNeighbor.neighborConfiguration = make(map[string]int)
						}
						currentNeighbor.neighborConfiguration[complementaryBorders[border]] = currentTile.ID
						currentNeighbor.locked = true
						tilemap[currentNeighbor.ID] = currentNeighbor
						break
					}
				}
				if matchError != nil {
					panic("impossible")
				}
			}
		}
		visitList = uniqueInts(nextVisitList)
	}
	return tilemap
}
func findPossibleNeighbors(tilemap map[int]Tile) map[int]Tile {
	//result := make(map[int]Tile)
	tileIDs := make([]int, 0, len(tilemap))
	for _, t := range tilemap {
		tileIDs = append(tileIDs, t.ID)
	}
	sort.Ints(tileIDs)
	foundNeighbors := make(map[int][]int)
	for _, tID := range tileIDs {
		if len(foundNeighbors[tID]) == 4 {
			continue
		}
		t := tilemap[tID]
		tileSides := t.getAllPossibleSides()
		tileSidesHashMap := make(map[int]bool)
		for _, ts := range tileSides {
			tileSidesHashMap[ts] = true
		}
		for _, possibleNeighborID := range tileIDs {
			if possibleNeighborID == t.ID {
				continue
			}
			if len(foundNeighbors[possibleNeighborID]) == 4 {
				continue
			}
			possibleNeighbor := tilemap[possibleNeighborID]
			neighborSides := possibleNeighbor.getAllPossibleSides()
			for _, neighboneighborSide := range neighborSides {
				checks++
				if tileSidesHashMap[neighboneighborSide] {
					//Found a neighbor
					foundNeighbors[t.ID] = append(foundNeighbors[t.ID], possibleNeighbor.ID)
					foundNeighbors[possibleNeighbor.ID] = append(foundNeighbors[possibleNeighbor.ID], t.ID)
					break
				}
			}
		}
	}
	for tID, n := range foundNeighbors {
		t := tilemap[tID]
		t.neighborIDs = uniqueInts(n)
		tilemap[tID] = t
	}
	return tilemap

}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	tiles := parse(trimmedInput)
	tiles = findPossibleNeighbors(tiles)
	corners := make([]int, 0)
	result := 1
	for _, t := range tiles {
		//fmt.Printf("Neighbors (%d): %+v\n", t.ID, t.neighborIDs)
		if len(t.neighborIDs) == 2 {
			corners = append(corners, t.ID)
			result *= t.ID
		}
	}
	//fmt.Printf("Corners: %+v", corners)
	fmt.Printf("Part1: %d\n", result)
	sort.Ints(corners)
	//fmt.Printf("Corners: %+v\n", uniqueInts(corners))

	tiles = buildTilemap(tiles[1427].ID, tiles)
	fmt.Printf("Checks :%d\n", checks)
	// for _, t := range tiles {
	// 	fmt.Printf("%d - \nNeighbors:\n", t.ID)
	// 	for side, neighbor := range t.neighborConfiguration {
	// 		fmt.Printf("\t%s - %d\n", side, neighbor)
	// 	}
	// }
}

var allSides = make(map[int][]int)

func parse(s []string) map[int]Tile {
	result := make(map[int]Tile, 10)

	currentTile := Tile{}
	rowIndex := 0
	currentTile.image = make(map[Bit]bool)

	for _, l := range s {
		if l == "" {
			result[currentTile.ID] = currentTile
			rowIndex = 0
			currentTile = Tile{}
			currentTile.image = make(map[Bit]bool)
		} else if strings.Index(l, "Tile") == 0 {
			tileID, _ := strconv.Atoi(
				strings.TrimRight(strings.Split(l, " ")[1], ":"))
			currentTile.ID = tileID
		} else {
			for colIndex, c := range l {
				if c == '#' {
					currentTile.image[Bit{Row: rowIndex, Col: colIndex}] = true
				} else {
					currentTile.image[Bit{Row: rowIndex, Col: colIndex}] = false
				}
			}
			rowIndex++
		}
	}
	result[currentTile.ID] = currentTile
	for _, t := range result {
		t.sideSignatures = make(map[string]map[bool]int)
		for _, f := range flipped {
			for _, b := range borders {
				if t.sideSignatures[b] == nil {
					t.sideSignatures[b] = make(map[bool]int)
				}
				s := t.getSideSignature(b)
				allSides[s] = uniqueInts(append(allSides[s], t.ID))
				t.sideSignatures[b][f] = s
			}
			t.flipTile()
		}
		result[t.ID] = t

	}
	return result

}
