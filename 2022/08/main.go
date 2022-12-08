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

type Coord struct {
	Row int
	Col int
}
type Direction string
type Tree struct {
	Height      int
	ScenicScore int
	Coord       Coord
}

var directions = []Coord{{Row: 1, Col: 0}, {Row: -1, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: -1}}

type TreeMap map[Coord]*Tree

func (t Tree) String() string {
	return fmt.Sprintf("{Height : %d , Coord %+v}", t.Height, t.Coord)
}
func (t *Tree) UpdateScenicScore(tmap TreeMap) {
	scenicScore := 1
	for _, d := range directions {
		currentCoord := AddCoords(t.Coord, d)
		neighbor, ok := tmap[currentCoord]
		visibleNeighbors := 0
		for ok {
			// If a neighbor exists, we count it
			visibleNeighbors = visibleNeighbors + 1
			// If the neighbors height is taller, we break out
			if neighbor.Height >= t.Height {
				break
			}
			currentCoord = AddCoords(neighbor.Coord, d)
			neighbor, ok = tmap[currentCoord]
		}
		scenicScore = scenicScore * visibleNeighbors
	}
	t.ScenicScore = scenicScore
}
func AddCoords(c Coord, c1 Coord) Coord {
	return Coord{Row: c.Row + c1.Row, Col: c.Col + c1.Col}
}
func visibleTrees(t TreeMap, startCoord Coord, direction Coord) []Coord {
	currentMaxHeight := -1
	result := make([]Coord, 0)
	currentCoord := startCoord
	currentTree, ok := t[currentCoord]
	// Do this to simplify going out of the map.
	nextCoord := AddCoords(currentCoord, direction)
	_, ok = t[nextCoord]
	for ok {
		if currentTree.Height > currentMaxHeight {
			currentMaxHeight = currentTree.Height
			result = append(result, currentCoord)
			// This is the max available height
			if currentMaxHeight == 9 {
				return result
			}
		}
		currentCoord = AddCoords(currentCoord, direction)
		currentTree, ok = t[currentCoord]
	}
	return result
}
func parseLines(input []string) TreeMap {
	result := make(TreeMap)
	for i, l := range input {
		for j, c := range strings.Split(l, "") {
			p := Coord{Row: i, Col: j}
			h, _ := strconv.Atoi(c)
			treeCandidate := new(Tree)
			treeCandidate.Height = h
			treeCandidate.Coord = p
			result[p] = treeCandidate
		}
	}
	return result
}
func main() {
	start := time.Now()
	flag.Parse()
	if *testMode {
		input = testInput
	}
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	t := time.Now()
	tmap := parseLines(trimmedInput)
	fmt.Println("Took(Parsing)", time.Since(t))
	t = time.Now()
	totalVisibleTrees := make([]Coord, 0)
	// for i := 0; i < len(trimmedInput); i++ {
	// 	for j := 0; j < len(trimmedInput); j++ {
	// 		v := tmap[Coord{i, j}]
	// 		fmt.Printf("\t%d", v.Height)
	// 	}
	// 	fmt.Printf("\n")
	// }
	for _, r := range []int{0, len(trimmedInput) - 1} {
		for c := 0; c < len(trimmedInput); c++ {
			for _, d := range directions {
				found := visibleTrees(tmap, Coord{Row: r, Col: c}, d)
				totalVisibleTrees = append(totalVisibleTrees, found...)
			}
		}
	}

	for _, c := range []int{0, len(trimmedInput) - 1} {
		for r := 0; r < len(trimmedInput); r++ {
			for _, d := range directions {
				found := visibleTrees(tmap, Coord{Row: r, Col: c}, d)
				totalVisibleTrees = append(totalVisibleTrees, found...)
			}
		}
	}
	uniqueCoords := make(map[Coord]bool)
	for _, v := range totalVisibleTrees {
		uniqueCoords[v] = true
	}
	fmt.Println("Took(Part1)", time.Since(t))
	fmt.Println("Part1", len(uniqueCoords))
	t = time.Now()
	maxScenicScore := 0
	for i := 0; i < len(trimmedInput); i++ {
		for j := 0; j < len(trimmedInput); j++ {
			t := tmap[Coord{j, i}]
			t.UpdateScenicScore(tmap)
			if t.ScenicScore > maxScenicScore {
				maxScenicScore = t.ScenicScore
			}
		}
	}
	fmt.Println("Took(Part2)", time.Since(t))
	fmt.Println("Part2", maxScenicScore)
	fmt.Println("Total time", time.Since(start))

	// for i := 0; i < len(trimmedInput); i++ {
	// 	for j := 0; j < len(trimmedInput); j++ {
	// 		v := tmap[Coord{i, j}]
	// 		if _, ok := uniqueCoords[Coord{i, j}]; ok {
	// 			fmt.Printf("\t*")

	// 		} else {
	// 			fmt.Printf("\t%d", v.Height)

	// 		}
	// 	}
	// 	fmt.Printf("\n")
	// }

}
