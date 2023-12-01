package main

import (
	_ "embed"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed input
var input string

//go:embed testInput
var testInput string

//go:embed testInputP2
var testInputP2 string
var testMode = flag.Bool("test", false, "Set to run using the testInput")

func calcSum(i []string) int {
	sum := 0
	for _, l := range i {
		firstInt := -1
		secondInt := -1
		for _, c := range strings.Split(l, "") {
			intVal, err := strconv.Atoi(c)
			if err != nil {
				continue
			}
			if firstInt == -1 {
				firstInt = intVal
			} else {
				secondInt = intVal
			}
		}
		if secondInt == -1 {
			secondInt = firstInt
		}
		sum = sum + firstInt*10 + secondInt
	}
	return sum
}
func main() {
	start := time.Now()
	flag.Parse()
	inputP1 := input
	inputP2 := input
	if *testMode {
		inputP1 = testInput
		inputP2 = testInputP2
	}
	trimmedInputP1 := strings.Split(strings.TrimSpace(string(inputP1)), "\n")
	trimmedInputP2 := strings.Split(strings.TrimSpace(string(inputP2)), "\n")

	replacedDigitsInput := replaceDigitsPart2(trimmedInputP2)
	part1 := calcSum(trimmedInputP1)
	part2 := calcSum(replacedDigitsInput)
	// fmt.Println("Replaced|\n\n", strings.Join(replacedDigitsInput, "\n"))
	fmt.Println("P1", part1)
	fmt.Println("P2", part2)
	fmt.Println("Took", time.Since(start))

}

func replaceDigitsPart2(input []string) []string {
	positionIndex := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
	}
	var re = regexp.MustCompile(`(?m).*[0-9]+.*`)
	for lineIndex, l := range input {
		minPosition := make([]int, len(positionIndex))
		maxPosition := make([]int, len(positionIndex))
		for p := range minPosition {
			minPosition[p] = 1000000000000
			maxPosition[p] = -1
		}
		found := false
		for k, v := range positionIndex {
			minPos := strings.Index(l, v)
			if minPos != -1 {
				found = true
				minPosition[k-1] = minPos
			}
		}
		if !found {
			// fmt.Println("skipping, not found")
			continue
		}
		minIndex := 0
		for index, pos := range minPosition {
			if pos < minPosition[minIndex] {
				minIndex = index
			}
		}
		valueToReplace := positionIndex[minIndex+1]
		if !re.MatchString(l[0:minPosition[minIndex]]) {
			// Theres no digit before the string, we replace
			// fmt.Println("Replacing", valueToReplace, "at", minPosition[minIndex], "with", strconv.Itoa(minIndex+1))
			input[lineIndex] = strings.Replace(l, valueToReplace, strconv.Itoa(minIndex+1), 1)
			l = input[lineIndex]
		}
		// At this point we found the minimal index of the string
		// Repeat the same for the max
		found = false
		for k, v := range positionIndex {
			maxPos := strings.LastIndex(l, v)
			if maxPos != -1 {
				found = true
				maxPosition[k-1] = maxPos
			}
		}
		if !found {
			// fmt.Println("skipping, not found")
			continue
		}
		maxIndex := len(positionIndex) - 1
		for index, pos := range maxPosition {
			if pos > maxPosition[maxIndex] {
				maxIndex = index
			}
		}
		valueToReplace = positionIndex[maxIndex+1]
		if re.MatchString(l[maxPosition[maxIndex]:]) {
			// Theres a digit before the string, we continue
			// fmt.Println("Theres a digit after the string, we continue")
			continue
		}
		//fmt.Println("Replacing", valueToReplace, "at", maxPosition[maxIndex], "with", strconv.Itoa(maxIndex+1))
		// fmt.Println("Replacing max", l[:maxPosition[maxIndex]], l[maxPosition[maxIndex]:])
		input[lineIndex] = l[:maxPosition[maxIndex]] + strings.Replace(l[maxPosition[maxIndex]:], valueToReplace, strconv.Itoa(maxIndex+1), 1)
		// fmt.Println("old line", l)
		// fmt.Println("New line", input[lineIndex])
	}
	return input
}
