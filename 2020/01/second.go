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
	for i, c1 := range numbers {
		for j, c2 := range numbers[i+1:] {
			if c1+c2 < 2020 {
				for k, c3 := range numbers[j+1:] {
					checks = checks + 1
					d := c1 + c2 + c3
					if d == 2020 {
						fmt.Printf("c1:%d, c2:%d, c3:%d\n", c1, c2, c3)
						fmt.Printf("r: %d\n", c1*c2*c3)
						fmt.Printf("Checks %d\n", checks)
						fmt.Printf("Indexes %d, %d, %d\n", i, j, k)
						os.Exit(0)
					}
				}
			}
		}
	}

}
