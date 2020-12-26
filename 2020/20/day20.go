package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

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

type Tile struct {
	ID          int
	image       map[Bit]bool
	neighborIDs []int
	neighbors   map[string]int
	locked      bool
	rotations   int
	flipped     bool
	minRowPos   int
	maxRowPos   int
	minColPos   int
	maxColPos   int
}

func (t Tile) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tile %d:\n", t.ID))
	sb.WriteString(fmt.Sprintf("Min position: %d,%d | Max position: %d,%d\n", t.minRowPos, t.minColPos, t.maxRowPos, t.maxColPos))

	sb.WriteString(" " + strings.Repeat("-", t.maxRowPos-t.minRowPos+2))
	sb.WriteString("\n")
	for row := t.minRowPos; row <= t.maxRowPos; row++ {
		sb.WriteString("|")
		for col := t.minColPos; col <= t.maxColPos; col++ {
			if t.image[Bit{Row: row, Col: col}] {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteString("|")
		sb.WriteString("\n")
	}
	sb.WriteString(" " + strings.Repeat("-", t.maxRowPos-t.minRowPos+2))
	sb.WriteString("\n")
	return sb.String()
}

func (t *Tile) rotateTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}
	t.rotations = (t.rotations + 1) % 4
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := t.minRowPos; r <= t.maxRowPos; r++ {
		for c := t.minColPos; c <= t.maxColPos; c++ {
			key := Bit{Row: r, Col: c}
			currentBit := t.image[key]
			nextKey := Bit{
				Row: t.maxRowPos - c,
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
		leftPoint = Bit{Row: t.minRowPos, Col: t.minColPos}
		rightPoint = Bit{Row: tt.maxRowPos, Col: tt.minColPos}
		colIncreaseL = 1
		colIncreaseR = 1
		rowIncreaseL = 0
		rowIncreaseR = 0
	case "bottom":
		leftPoint = Bit{Row: t.maxRowPos, Col: t.minColPos}
		rightPoint = Bit{Row: t.minColPos, Col: tt.minColPos}
		colIncreaseL = 1
		colIncreaseR = 1
		rowIncreaseL = 0
		rowIncreaseR = 0
	case "right":
		leftPoint = Bit{Row: t.minRowPos, Col: t.maxColPos}
		rightPoint = Bit{Row: tt.minRowPos, Col: tt.minColPos}
		colIncreaseL = 0
		colIncreaseR = 0
		rowIncreaseL = 1
		rowIncreaseR = 1
	case "left":
		leftPoint = Bit{Row: t.minRowPos, Col: t.minColPos}
		rightPoint = Bit{Row: t.minRowPos, Col: tt.maxColPos}
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
func (t *Tile) stripSides() {
	newBitMap := make(map[Bit]bool, len(t.image))
	var row, col int
	for row = t.minRowPos + 1; row <= (t.maxRowPos - 1); row++ {
		for col = t.minColPos + 1; col <= (t.maxColPos - 1); col++ {
			key := Bit{Col: col, Row: row}
			remappedKey := Bit{Col: col - 1, Row: row - 1}
			newBitMap[remappedKey] = t.image[key]
		}
	}
	t.maxRowPos = t.maxRowPos - 2
	t.maxColPos = t.maxColPos - 2
	t.image = newBitMap

}
func (t *Tile) flipTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}
	t.flipped = !t.flipped
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := t.minRowPos; r <= t.maxRowPos; r++ {
		for c := t.minColPos; c <= t.maxColPos; c++ {
			key := Bit{Row: r, Col: c}
			currentBit := t.image[key]
			nextKey := Bit{
				Row: t.maxRowPos - r,
				Col: c,
			}
			newBitMap[nextKey] = currentBit
		}
	}
	t.image = newBitMap
	return nil
}

func (t *Tile) topNeighborFunc() int {

	result := 0
	for col := t.minColPos; col <= t.maxColPos; col++ {
		keyToCheck := Bit{Row: 0, Col: col}
		if t.image[keyToCheck] {
			result += 1 << col
		}
	}
	return result
}

func (t *Tile) bottomNeighborFunc() int {

	result := 0
	for col := t.minColPos; col <= t.maxColPos; col++ {
		keyToCheck := Bit{Row: 9, Col: col}
		if t.image[keyToCheck] {
			result += 1 << col
		}
	}
	return result
}

func (t *Tile) rightNeighborFunc() int {
	result := 0
	for row := t.minRowPos; row <= t.maxRowPos; row++ {
		keyToCheck := Bit{Row: row, Col: 9}
		if t.image[keyToCheck] {
			result += 1 << row
		}
	}
	return result
}
func (t *Tile) leftNeighborFunc() int {
	result := 0
	for row := t.minRowPos; row <= t.maxRowPos; row++ {
		keyToCheck := Bit{Row: row, Col: 0}
		if t.image[keyToCheck] {
			result += 1 << row
		}
	}
	return result
}

//In order to avoid conflicts of orientation, it's
//easier to just check a single place and just rotate the
//tile.
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
	//startingTile.flipTile()
	tilemap[startTile] = startingTile
	for len(visitedTiles) != len(tilemap) {
		//fmt.Printf("Next visit: %+v\n", visitList)
		var nextVisitList []int
		if len(visitList) == 0 {
			panic("impossible ")
		}
		for _, tileID := range visitList {
			//fmt.Printf("Visiting %d\n", tileID)
			if visitedTiles[tileID] {
				//fmt.Printf("Already seen %d\n", tileID)
				continue
			}
			currentTile := tilemap[tileID]
			visitedTiles[tileID] = true
			if len(currentTile.neighbors) == len(currentTile.neighborIDs) {
				continue
			}

			for _, neighborID := range currentTile.neighborIDs {
				if visitedTiles[neighborID] {
					//fmt.Printf("Skipping %d as we already visited it\n", neighborID)
					continue
				}
				currentNeighbor := tilemap[neighborID]
				if len(currentNeighbor.neighbors) == len(currentNeighbor.neighborIDs) {
					continue
				}
				nextVisitList = append(nextVisitList, neighborID)
				var matchError error
				for _, border := range borders {
					if currentNeighbor.neighbors != nil {
						_, ok := currentNeighbor.neighbors[complementaryBorders[border]]
						if ok {
							continue
						}
					}

					//This tile already has a neighbor here
					if currentTile.neighbors != nil {
						_, ok := currentTile.neighbors[border]
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
						if currentTile.neighbors == nil {
							currentTile.neighbors = make(map[string]int)
						}
						currentTile.neighbors[border] = currentNeighbor.ID
						currentTile.locked = true
						tilemap[currentTile.ID] = currentTile

						if currentNeighbor.neighbors == nil {
							currentNeighbor.neighbors = make(map[string]int)
						}
						currentNeighbor.neighbors[complementaryBorders[border]] = currentTile.ID
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
	tp1 := time.Now()
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
	sort.Ints(corners)
	fmt.Printf("Corners: %+v\n", corners)
	fmt.Printf("Part1: %d\n", result)
	//fmt.Printf("Corners: %+v\n", uniqueInts(corners))
	//as it is built right now we can start from any place
	fmt.Printf("P1 Took %s\n", time.Since(tp1))
	t := time.Now()
	defer func() {
		fmt.Printf("Took %s\n", time.Since(t))
	}()
	tiles = buildTilemap(corners[0], tiles)
	var startingTile int
	for _, c := range corners {
		t := tiles[c]
		_, hasTop := t.neighbors["top"]
		_, hasLeft := t.neighbors["left"]
		if !hasTop && !hasLeft {
			startingTile = c
			break
		}
	}
	fmt.Printf("\n%d -> %+v\n", startingTile, tiles[startingTile].neighbors)
	//This fusion only works starting at the top corners.
	bigPicture := fuseTiles(startingTile, tiles)
	for range borders {
		for range flipped {
			monstersFound := findMonsters(bigPicture)
			if monstersFound != 0 {
				fmt.Printf("Found %d monsters with %d rotations, %+v fliped \n", findMonsters(bigPicture), bigPicture.rotations, bigPicture.flipped)
				totalHashtags := 0
				for _, p := range bigPicture.image {
					if p {
						totalHashtags++
					}
				}
				fmt.Printf("Rough index:%d\n", totalHashtags-monstersFound*15)
				//We can only find monsters in one specific configuration.
				return
			}
			bigPicture.flipTile()
		}
		bigPicture.rotateTile()
	}

}

func getMonsterPoints(b Bit, maxRow, maxCol int) []Bit {
	var result []Bit

	if b.Col+19 > maxCol {
		return result
	}
	if b.Row+2 > maxRow {
		return result
	}
	result = append(result, Bit{Col: b.Col + 18, Row: b.Row})

	result = append(result, Bit{Col: b.Col, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 5, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 6, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 11, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 12, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 17, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 18, Row: b.Row + 1})
	result = append(result, Bit{Col: b.Col + 19, Row: b.Row + 1})

	result = append(result, Bit{Col: b.Col + 1, Row: b.Row + 2})
	result = append(result, Bit{Col: b.Col + 4, Row: b.Row + 2})
	result = append(result, Bit{Col: b.Col + 7, Row: b.Row + 2})
	result = append(result, Bit{Col: b.Col + 10, Row: b.Row + 2})
	result = append(result, Bit{Col: b.Col + 13, Row: b.Row + 2})
	result = append(result, Bit{Col: b.Col + 16, Row: b.Row + 2})

	return result

}
func findMonsters(tile Tile) int {

	var monsters int
	for row := tile.minRowPos; row <= tile.maxRowPos; row++ {
		for col := tile.minColPos; col <= tile.maxColPos; col++ {
			pointsToCheck := getMonsterPoints(Bit{Row: row, Col: col}, tile.maxRowPos, tile.maxColPos)
			if len(pointsToCheck) == 0 {
				continue
			}
			matches := true
			for _, p := range pointsToCheck {
				if !tile.image[p] {
					matches = false
					break
				}
			}
			if matches {
				monsters++
			}
		}
	}
	return monsters

}

//Right now this works by creating a spiral, so it works in a square as long as you start from one of the outer
//borders. In order for this to work we would need to be able to walk back
// and go to a point where going in a different direction is possible (A* perhaps?)
func fuseTiles(startingTile int, tilemap map[int]Tile) Tile {

	visitedTiles := make(map[int]bool)
	currentTile := tilemap[startingTile]
	rowIndex := 0
	colIndex := 0
	bigPicture := Tile{}
	bigPicture.image = make(map[Bit]bool)
	directions := []string{"right", "top", "left", "bottom"}
	directionIndex := 0

	type rowColIndexModifier struct {
		row int
		col int
	}
	rowColIndexModifierMap := make(map[string]rowColIndexModifier)
	rowColIndexModifierMap["right"] = rowColIndexModifier{row: 0, col: 1}
	rowColIndexModifierMap["top"] = rowColIndexModifier{row: -1, col: 0}
	rowColIndexModifierMap["bottom"] = rowColIndexModifier{row: 1, col: 0}

	if _, ok := currentTile.neighbors["top"]; ok {
		rowColIndexModifierMap["top"] = rowColIndexModifier{row: 1, col: 0}
		rowColIndexModifierMap["bottom"] = rowColIndexModifier{row: -1, col: 0}
	}
	rowColIndexModifierMap["left"] = rowColIndexModifier{row: 0, col: -1}

	currentDirectionModifier := rowColIndexModifierMap[directions[directionIndex]]
	numberOfDirectionsChanges := 0
	visitationOrder := make([]int, 0, 144)
	for len(visitedTiles) != len(tilemap) {

		visitedTiles[currentTile.ID] = true
		visitationOrder = append(visitationOrder, currentTile.ID)
		//Because we start at the top left corner, we mark that point as 0,0. This means that we need to
		//invert the Y axis to drive it down.
		//currentTile.invertRowCoords()
		//currentTile.flipTile()
		currentTile.stripSides()

		//remap all the bits using the row and col index
		var row, col int
		for row = currentTile.minRowPos; row <= currentTile.maxRowPos; row++ {
			transplantedRow := row + rowIndex*8
			for col = currentTile.minColPos; col <= currentTile.maxRowPos; col++ {
				k := Bit{Col: col, Row: row}
				remappedKey := Bit{Col: k.Col + colIndex*8, Row: transplantedRow}
				bigPicture.image[remappedKey] = currentTile.image[k]
				//Add all the points to the "bigPicture"
				if remappedKey.Col < bigPicture.minColPos {
					bigPicture.minColPos = remappedKey.Col
				} else if remappedKey.Col > bigPicture.maxColPos {
					bigPicture.maxColPos = remappedKey.Col
				}
				if remappedKey.Row < bigPicture.minRowPos {
					bigPicture.minRowPos = remappedKey.Row

				} else if remappedKey.Row > bigPicture.maxRowPos {
					bigPicture.maxRowPos = remappedKey.Row
				}
			}
		}
		//To pick the next tile to visit:
		for {
			// If you have been everywhere, congratulations!
			if len(visitedTiles) == len(tilemap) {
				break
			}
			// Check if a neighbor is available at the current directions[directionIndex] and hasn't been visited.
			if neighbor, ok := currentTile.neighbors[directions[directionIndex]]; ok && !visitedTiles[neighbor] {
				currentTile = tilemap[neighbor]
				numberOfDirectionsChanges = 0
				break
			} else {
				// If there's none, or it has been visited, increment the directionIndex and try again
				directionIndex = (directionIndex + 1) % len(directions)
				// Whenever the direction index changes you need to update the rowIndexModifier and colIndexModifiers to the correct values
				currentDirectionModifier = rowColIndexModifierMap[directions[directionIndex]]
				numberOfDirectionsChanges++
			}
			if numberOfDirectionsChanges > 4 {
				//If this happens it means you're lost and there's no way to continue
				panic("bad")
			}
		}
		colIndex += currentDirectionModifier.col
		rowIndex += currentDirectionModifier.row

	}
	return bigPicture

}

var allSides = make(map[int][]int)

func parse(s []string) map[int]Tile {
	result := make(map[int]Tile, 10)

	currentTile := Tile{}
	rowIndex := 0
	currentTile.image = make(map[Bit]bool)
	for _, l := range s {
		currentTile.maxColPos = 9
		currentTile.maxRowPos = 9
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
		t.minColPos = 0
		t.minRowPos = 0
		t.maxColPos = 9
		t.maxRowPos = 9
		result[t.ID] = t

	}
	return result

}
