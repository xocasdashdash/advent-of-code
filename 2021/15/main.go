package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

type Coord struct {
	X int
	Y int
}
type Node struct {
	Value int
}

var visited map[Coord]bool

func parseInput(input []string) CostMap {
	result := make(CostMap)
	for x, l := range input {
		for y, c := range l {
			v, _ := strconv.Atoi(string(c))
			result[Coord{X: x, Y: y}] = v
		}
	}
	return result
}

type CostMap map[Coord]int
func (c CostMap)String()string{
	sb := strings.Builder{}
	for i := 0; i<10;i++ {
		for j:=0;j<10;j++ {
			sb.WriteString(fmt.Sprintf("|\t%d\t|", c[Coord{i,j}]))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
func (c CostMap) Neighbors(p Coord) []Coord {
	neighbors := make([]Coord, 0, 10)
	possibleNeighbors := []Coord{
		{p.X, p.Y + 1}, 
		{p.X, p.Y - 1},
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
	}
	for _, p := range possibleNeighbors {
		if _, ok := c[p]; ok {
			neighbors = append(neighbors, p)
		}
	}
	return neighbors
}
func initializeCostMap(rows, cols int) CostMap {
	m := make(map[Coord]int)
	v := 100000000000000000
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			m[Coord{i, j}] = v
		}
	}
	return m
}
func filterNeighbors(coords []Coord, visitedNodes map[Coord]bool) []Coord {
	res := make([]Coord, 0, len(coords))
	for _, c := range coords {
		if _, ok := visitedNodes[c]; !ok {
			res = append(res, c)
		}
	}
	return res
}
func minimalNeighbor(neighborCoords []Coord, nodes CostMap, visited map[Coord]bool) *Coord {
	var minCoord *Coord
	if len(neighborCoords) == 0 {
		return nil
	}
	value := 1000000000000000000
	for _, c := range neighborCoords {
		if value >= nodes[c] && !visited[c]{
			value = nodes[c]
			t := c
			minCoord = &t
		}
	}
	return minCoord
}
func visitNodes(nodes CostMap, costs CostMap) CostMap {
	visitedNodes := make(map[Coord]bool)
	currentNode := Coord{0, 0}
	costs[currentNode] = 0
	pendingVisit := make([]Coord, 0, 10)
	prevVisited := len(visitedNodes)
	for len(visitedNodes) != len(nodes) {
		if prevVisited == len(visitedNodes) {
			fmt.Println("Bad")
		}
		neighbors := costs.Neighbors(currentNode)
		pendingVisit = append(pendingVisit, neighbors...)
		for _, n := range neighbors {
			currentCost := costs[n]
			newCost := costs[currentNode] + nodes[n]
			if newCost <= currentCost {
				costs[n] = newCost
			}
		}
		prevVisited = len(visitedNodes)
		visitedNodes[currentNode] = true
		nextNode := minimalNeighbor(neighbors, nodes, visitedNodes)
		if nextNode == nil {
			for k,possibleNextNode  := range pendingVisit {
				if ok := visitedNodes[possibleNextNode];!ok {
					currentNode = possibleNextNode
					pendingVisit = pendingVisit[k:]
					break
				}
			}
		}else {
			currentNode = *nextNode
		}

	}
	return costs
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	nodes := parseInput(trimmedInput)
	costMap := initializeCostMap(len(trimmedInput), len(trimmedInput[0]))
	costMap = visitNodes(nodes, costMap)
	fmt.Println("nodeMap")
	fmt.Println(nodes)
	fmt.Println("costMap")
	fmt.Println(costMap)
	fmt.Println("9,9", costMap[Coord{9,9}])

}
