package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//Homework The input
type Homework struct {
	Input  string
	output int
}

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var debugFlag = flag.Bool("debug", false, "Run with debug")

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	homework := readHomework(trimmedInput)

	sum := 0
	for _, h := range homework {
		value, err := Parse(h.Input, "v1")
		if err != nil {
			panic(err)
		}
		if h.output != 0 && h.output != value && *debugFlag {
			fmt.Printf("Value does not match for line %s\n. Expected: %d, Got: %d\n", h.Input, h.output, value)
		} else {
			sum += value
		}
	}
	if sum != 0 {
		fmt.Printf("Part1: %d\n", sum)
	}
	sum = 0
	for _, h := range homework {
		if *debugFlag {
			fmt.Printf("Input %s\n", h.Input)
		}
		value, err := Parse(h.Input, "v2")
		if err != nil {
			panic(err)
		}
		sum += value
		if *debugFlag {
			if h.output != 0 && h.output != value && *debugFlag {
				fmt.Printf("KO! Value does not match for line %s Expected: %d, Got: %d\n", h.Input, h.output, value)
			} else {
				fmt.Printf("OK! Value does match for line %s. Expected: %d, Got: %d\n", h.Input, h.output, value)
			}
		}
	}
	if sum != 0 {
		fmt.Printf("Part2: %d\n", sum)
	}

}

func readHomework(lines []string) []Homework {
	var result []Homework
	for _, l := range lines {
		s := strings.Split(l, ",")
		var h Homework
		h.Input = strings.Replace(s[0], " ", "", -1)
		if len(s) > 1 {
			output, _ := strconv.Atoi(s[1])
			h.output = output
		}
		result = append(result, h)
	}
	return result
}
