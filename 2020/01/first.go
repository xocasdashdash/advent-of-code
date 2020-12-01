package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	f, _ := os.Open("input")

	numbers := make([]int, 0, 20)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		numbers = append(numbers, n)
	}
	checks := 0
	for k, c1 := range numbers {
		if c1 < 2020 {
			for _, c2 := range numbers[k+1:] {
				checks = checks + 1

				d := c1 + c2
				if d == 2020 {
					fmt.Printf("c1:%d, c2:%d\n", c1, c2)
					fmt.Printf("r: %d\n", c1*c2)
					fmt.Printf("Checks %d", checks)
					os.Exit(0)
				}
			}
		}
	}

}
