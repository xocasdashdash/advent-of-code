package main

import (
	"fmt"
)

const (
	puzzleInput = 325489
)

func main() {
	//part1()
	part2()
}

func part1() {
	a := 571
	mid := 285
	s := make([][]int, a)
	for i := range s {
		s[i] = make([]int, a)
	}
	s[mid][mid] = 1
	l := a * a
	k := 1
	for i := 0; i < l; i++ {
		for j := 0; j < 2*i+2; j++ {
			y := mid + i - j
			x := mid + i + 1
			k++
			if k >= puzzleInput {
				fmt.Println(abs(x-mid) + abs(y-mid))
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid - i - 1
			x := mid + i - j
			k++
			if k >= puzzleInput {
				fmt.Println(abs(x-mid) + abs(y-mid))
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid - i + j
			x := mid - i - 1
			k++
			if k >= puzzleInput {
				fmt.Println(abs(x-mid) + abs(y-mid))
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid + i + 1
			x := mid - i + j
			k++
			if k >= puzzleInput {
				fmt.Println(abs(x-mid) + abs(y-mid))
				return
			}
			s[y][x] = k
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part2() {
	a := 571
	mid := 285
	s := make([][]int, a)
	for i := range s {
		s[i] = make([]int, a)
	}
	s[mid][mid] = 1
	l := (a - 1) * (a - 1)
	for i := 0; i < l; i++ {
		for j := 0; j < 2*i+2; j++ {
			y := mid + i - j
			x := mid + i + 1
			k := sumSquare(s, y, x)
			if k > puzzleInput {
				fmt.Println(k)
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid - i - 1
			x := mid + i - j
			k := sumSquare(s, y, x)
			if k > puzzleInput {
				fmt.Println(k)
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid - i + j
			x := mid - i - 1
			k := sumSquare(s, y, x)
			if k > puzzleInput {
				fmt.Println(k)
				return
			}
			s[y][x] = k
		}
		for j := 0; j < 2*i+2; j++ {
			y := mid + i + 1
			x := mid - i + j
			k := sumSquare(s, y, x)
			if k > puzzleInput {
				fmt.Println(k)
				return
			}
			s[y][x] = k
		}
	}
}

func sumSquare(s [][]int, r, c int) int {
	return s[r][c+1] + s[r-1][c+1] + s[r-1][c] + s[r-1][c-1] + s[r][c-1] + s[r+1][c-1] + s[r+1][c] + s[r+1][c+1]
}
