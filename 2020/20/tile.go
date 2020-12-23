package main

import "sort"

func (t *Tile) candidateNeighbors() []int {
	var result []int
	result = append(result, t.localNeighbor("top")...)
	result = append(result, t.localNeighbor("right")...)
	result = append(result, t.localNeighbor("bottom")...)
	result = append(result, t.localNeighbor("left")...)
	sort.Ints(result)
	return result
}
