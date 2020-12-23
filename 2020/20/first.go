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
	"bottom", "right", "top", "left",
}

type Bit struct {
	X, Y int
}
type TileConfiguration struct {
	rotations int
	flip      bool
}
type NeighborMatch struct {
	ID         int
	selfTC     TileConfiguration
	neighborTC TileConfiguration
}

func (nc NeighborMatch) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("{ ID: %d, TC: { rotations: %d, flipped %+v }", nc.ID, nc.selfTC.rotations, nc.selfTC.flip))
	return sb.String()
}

type Tile struct {
	ID          int
	image       map[Bit]bool
	neighborIDs []int
	//Map with up to four keys top,bottom,left,right
	//This will allow me to start from wherever and
	//go wherever "forcing" my neighbors to be configured in a specific
	//way. This is a list to allow for the case when you can have neighbors match in multiple configurations
	neighborConfiguration map[string]map[int][]NeighborMatch
	TileConfig            TileConfiguration
	//When set, flip and rotate operations will panic
	locked bool
}

func breakFN() {
	if true {
		//runtime.Breakpoint()
	}
}
func (t *Tile) resetTile() {
	reverseRotations := 4 - t.TileConfig.rotations
	for i := 0; i < reverseRotations; i++ {
		t.rotateTile()
	}
	t.TileConfig.rotations = 0
	if t.TileConfig.flip {
		t.flipTile()
	}
	t.TileConfig.flip = false
}
func (t *Tile) flipTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}
	t.TileConfig.flip = !t.TileConfig.flip
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			key := Bit{X: r, Y: c}
			currentBit := t.image[key]
			if currentBit {
				nextKey := Bit{X: 9 - r, Y: c}
				newBitMap[nextKey] = true
			}
		}
	}
	t.image = newBitMap
	return nil
}

type neighborFunc func(tt *Tile) int

func (t *Tile) rotateTile() error {
	if t.locked {
		return fmt.Errorf(fmt.Sprintf("Tile with id %d, is locked. Cannot modify configuration", t.ID))
	}

	t.TileConfig = TileConfiguration{
		rotations: (t.TileConfig.rotations + 1) % 4,
		flip:      t.TileConfig.flip,
	}
	newBitMap := make(map[Bit]bool, len(t.image))
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			key := Bit{X: r, Y: c}
			currentBit := t.image[key]
			if currentBit {
				nextKey := Bit{X: 9 - c, Y: r}
				newBitMap[nextKey] = true
			}
		}
	}
	t.image = newBitMap
	return nil
}
func (t *Tile) topNeighborFunc() neighborFunc {
	return func(tt *Tile) int {
		localResult := 0
		for x := 0; x < 10; x++ {
			keyToCheck := Bit{x, 0}
			if t.image[keyToCheck] {
				localResult += (2 << x) / 2
			}
		}
		return localResult
	}
}
func (t *Tile) bottomNeighborFunc() neighborFunc {
	return func(tt *Tile) int {

		localResult := 0
		for x := 0; x < 10; x++ {
			keyToCheck := Bit{x, 9}
			if t.image[keyToCheck] {
				localResult += (2 << x) / 2
			}
		}
		return localResult
	}
}
func (t *Tile) rightNeighborFunc() neighborFunc {

	return func(tt *Tile) int {
		localResult := 0
		for x := 0; x < 10; x++ {
			keyToCheck := Bit{9, x}
			if tt.image[keyToCheck] {
				localResult += (2 << x) / 2
			}
		}
		return localResult
	}
}
func (t *Tile) leftNeighborFunc() neighborFunc {
	return func(tt *Tile) int {
		localResult := 0
		for x := 0; x < 10; x++ {
			keyToCheck := Bit{0, x}
			if tt.image[keyToCheck] {
				localResult += (2 << x) / 2
			}
		}
		return localResult
	}
}

//reconfigures the tile until the border specified matches the input
func (t *Tile) matchBorder(input int, border string) error {

	var borderFunc neighborFunc
	switch border {
	case "top":
		borderFunc = t.topNeighborFunc()
	case "bottom":
		borderFunc = t.bottomNeighborFunc()
	case "left":
		borderFunc = t.leftNeighborFunc()
	case "right":
		borderFunc = t.rightNeighborFunc()
	default:
		panic("impossible")
	}
	var borderResult int
	if t.locked {
		borderResult = borderFunc(t)
		if borderResult == input {
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("No border matches. Border result %d", borderResult))
	}
	//if t.ID == 1489 || t.ID == 2311 || t.ID == 2729 || t.ID == 2473 {
	//	runtime.Breakpoint()
	//}
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			borderResult = borderFunc(t)
			if borderResult == input {
				return nil
			}
			t.flipTile()
		}
		t.rotateTile()
	}
	return fmt.Errorf(fmt.Sprintf("No border matches. Border result %d", borderResult))

}
func (t *Tile) localNeighbor(side string) []int {

	rotations := 0
	switch strings.ToLower(side) {
	case "right":
		rotations = 1
	case "bottom":
		rotations = 2
	case "left":
		rotations = 3
	}
	reverseRotations := 4 - rotations
	result := make([]int, 2)
	var index int
	for rotations > 0 {
		t.rotateTile()
		rotations--
	}
	if t.TileConfig.flip {
		index = 1
	}
	result[index] = t.topNeighborFunc()(t)
	t.flipTile()
	if index == 1 {
		index = 0
	} else {
		index = 1
	}
	result[index] = t.topNeighborFunc()(t)
	t.flipTile()
	for reverseRotations > 0 && reverseRotations != 4 {
		t.rotateTile()
		reverseRotations--
	}
	return result
}

func (t Tile) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tile %d:\n", t.ID))
	if !t.locked {
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if t.image[Bit{i, j}] {
					sb.WriteString("#")
				} else {
					sb.WriteString(".")
				}
			}
			sb.WriteString("\n")
		}
	}
	sb.WriteString("Neighbors\n")
	for _, b := range borders {
		sb.WriteString(fmt.Sprintf("%s possible neighbors:(%d) \n", b, len(t.neighborConfiguration[b])))
		//for _, matchConfig := range t.neighborConfiguration[b] {
		//	sb.WriteString(fmt.Sprintf("\tSelf config: %+v|", matchConfig.selfTC))
		//	sb.WriteString(fmt.Sprintf("Neighbor config: %+v\n", matchConfig.neighborTC))
		//}
		sb.WriteString("\n")
	}
	return sb.String()

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
	return list

}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	tiles := parse(trimmedInput)
	//for _, t := range tiles {
	//	fmt.Printf("%s\n", t)
	//}
	var tileNeighborCandidates map[int][]int
	tileNeighborCandidates = make(map[int][]int)
	for _, t := range tiles {
		candidateNeighbors := t.candidateNeighbors()
		sort.Ints(candidateNeighbors)
		tileNeighborCandidates[t.ID] = candidateNeighbors
	}
	var tileNeighbors map[int]map[int][]int
	tileNeighbors = make(map[int]map[int][]int, 10)
	//At this point we have all the candidate neighbors for a
	//tile. Now we need to filter them

	for tileID, neighborsWallSignatures := range tileNeighborCandidates {
		tileNeighborsID, ok := tileNeighbors[tileID]
		if !ok {
			tileNeighborsID = make(map[int][]int)
		}
		for _, neighborsWallSignature := range neighborsWallSignatures {
			for possibleNeighbor, availableNeighbors := range tileNeighborCandidates {
				if possibleNeighbor == tileID {
					continue
				}
				for _, availableNeighbor := range availableNeighbors {
					if availableNeighbor == neighborsWallSignature {
						tileNeighborsID[possibleNeighbor] = uniqueInts(append(tileNeighborsID[possibleNeighbor], neighborsWallSignature))
						tileNeighbors[tileID] = tileNeighborsID
						//This also implies the reverse. And it also implies that a neighbor
						//of my neighbor cannot be my neighbor.
						foundNeighbors, ok := tileNeighbors[possibleNeighbor]
						if !ok {
							foundNeighbors = make(map[int][]int)
						}
						foundNeighbors[tileID] = uniqueInts(append(foundNeighbors[tileID], availableNeighbor))
						tileNeighbors[possibleNeighbor] = foundNeighbors
					}
				}
				//if len(tileNeighborsID) == 4 {
				//	break
				//}
			}
		}
	}
	result := 1
	corners := make([]int, 0, 4)
	for tileID, n := range tileNeighbors {
		//fmt.Printf("neighbors for %d can be (%d) %+v\n", tileID, len(n), n)
		if len(n) == 2 {
			result *= tileID
			corners = append(corners, tileID)
		}
	}
	fmt.Printf("Part1: %d\n", result)

	startTile := 1951
	//toVisit := []int{startTile}
	t := tiles[startTile]

	tiles[startTile] = t

	visitedTiles := make(map[int]bool)

	visitList := []int{startTile}
	//The working tilemap configuration. I should only write here once.
	tilemap := make(map[int]TileConfiguration)
	for len(visitedTiles) != len(tiles) {
		var nextVisitList []int

		for _, tileID := range visitList {
			currentTile := tiles[tileID]
			currentTile.resetTile()
			visitedTiles[tileID] = true
			if currentTile.ID == 1171 {
				breakFN()
			}
			if len(currentTile.neighborConfiguration) == len(tileNeighbors[currentTile.ID]) {
				//We've already gotten all the information
				continue
			}
			currentNeighbors := tileNeighbors[currentTile.ID]
			for neighborID, matchingBorders := range currentNeighbors {
				//if _, ok := visitedTiles[neighborID]; ok {
				//	//If i have already visited the node we can skip it
				//	continue
				//}
				nextVisitList = append(nextVisitList, neighborID)
				neighborTile := tiles[neighborID]
				var matchBorderError error

				for _, border := range borders {
					//if _, ok := currentTile.neighborConfiguration[border]; ok {
					//	fmt.Printf("Tile %d already has a neighbor at %s\n", currentTile.ID, border)
					//	continue
					//}
					for _, matchingBorder := range matchingBorders {
						//fmt.Printf("Matching tile(%d) border %s with neighbor (%d). Trying side %s with signature %d\n", currentTile.ID, border, neighborTile.ID, complementaryBorders[border], matchingBorder)
						matchBorderError = neighborTile.matchBorder(matchingBorder, complementaryBorders[border])
						if matchBorderError == nil {
							//fmt.Printf("Got match! tile %d(%s) neighbour of %d(%s)\n", currentTile.ID, border, neighborTile.ID, complementaryBorders[border])
							if currentTile.neighborConfiguration == nil {
								currentTile.neighborConfiguration = make(map[string]map[int][]NeighborMatch)
							}
							if currentTile.neighborConfiguration[border] == nil {
								currentTile.neighborConfiguration[border] = make(map[int][]NeighborMatch)
							}
							//Set neighbor configuration for current tile
							currentTile.neighborConfiguration[border][neighborTile.ID] = append(currentTile.neighborConfiguration[border][neighborTile.ID], NeighborMatch{
								ID:         neighborTile.ID,
								selfTC:     currentTile.TileConfig,
								neighborTC: neighborTile.TileConfig,
							})
							currentTile.neighborIDs = uniqueInts(append(currentTile.neighborIDs, neighborTile.ID))
							tiles[currentTile.ID] = currentTile
							//We also set the neighborconfiguration for the other tile.
							if neighborTile.neighborConfiguration == nil {
								neighborTile.neighborConfiguration = make(map[string]map[int][]NeighborMatch)
							}
							if neighborTile.neighborConfiguration[complementaryBorders[border]] == nil {
								neighborTile.neighborConfiguration[complementaryBorders[border]] = make(map[int][]NeighborMatch)
							}
							neighborTile.neighborConfiguration[complementaryBorders[border]][currentTile.ID] = append(neighborTile.neighborConfiguration[complementaryBorders[border]][currentTile.ID], NeighborMatch{
								ID:         currentTile.ID,
								neighborTC: currentTile.TileConfig,
								selfTC:     neighborTile.TileConfig,
							})
							neighborTile.neighborIDs = uniqueInts(append(neighborTile.neighborIDs, currentTile.ID))
							tiles[neighborTile.ID] = neighborTile
							break
						}
						neighborTile.resetTile()
					}
				}
			}
			currentTile.resetTile()
			tiles[tileID] = currentTile
		}

		visitList = uniqueInts(nextVisitList)
		nextVisitList = make([]int, 0)

	}

	fmt.Printf("debug test")
	//i := 0
	for _, v := range tiles {
		fmt.Printf("%+v\n", v)
	}
	//			nextVisit = append(nextVisit, tiles[id].neighborConfiguration[b].ID)
	//		}
	//	}
	//	toVisit = nextVisit
	//	nextVisit = make([]int, 0)
	//	if i > 200 {
	//		panic("bad")
	//	}
	//}

	//Idea 1: In order to solve this, treat this like a sudoku problem
	//We have some constraints (the four corners must be in the corners (0,0),(0,n),(n,0),(n,n) )
	//And because we know what corners match or not we can detect the specific corner that we have.
	//At the beginning, we need to pick one corner and give it the position (0,0) as a reference.
	//This makes it that the neighbors of this corner will be in known positions with a specific rotation and flipped or not
	//An example, if the neighbors are on the L,S it means we're on the top right corner. Then we match the S side with the north side of the known neighbor
	//and now we know that neighbors position, orientation and rotation to match this situation.

	//Idea 2: You can also find the center ingredient, with this you can find the correct configuration for it as it has 4 known constraints
	//Once you have that, you can find the configuration of all your neighbors (maximum 8 configurations), and move from there. Once that is done
	//You can find the position of any other member and their configuration.
	//=> Bad idea, this works on the minimal input (3x3) as there's one center, we have multiple cases of tiles with 4 borders.

	//Idea 3: Pick one corner, for each neighbor find the right configuration for it, save it and then go to one of the neighbors.
	//Repeat this until you have visited every node.

	//Part 2: Once you have the board configuration, move it to a new tile, print the map and ran a regex scan for every 3 lines

}

func parse(s []string) map[int]Tile {
	result := make(map[int]Tile, 10)

	currentTile := Tile{}
	tileIndex := 0
	rowIndex := 0
	currentTile.image = make(map[Bit]bool)
	for _, l := range s {
		if l == "" {
			tileIndex++
			result[currentTile.ID] = currentTile
			tileIndex++
			rowIndex = 0
			currentTile = Tile{}
			currentTile.image = make(map[Bit]bool)
			continue
		}
		if strings.Index(l, "Tile") == 0 {
			tileID, _ := strconv.Atoi(
				strings.TrimRight(strings.Split(l, " ")[1], ":"))
			currentTile.ID = tileID
		} else {
			for colIndex, c := range l {
				if c == '#' {
					currentTile.image[Bit{rowIndex, colIndex}] = true
				} else {
					currentTile.image[Bit{rowIndex, colIndex}] = false
				}
			}
			rowIndex++
		}
	}
	result[currentTile.ID] = currentTile

	return result

}
