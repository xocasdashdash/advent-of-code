package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var initialValue = flag.Int("n", 0, "Value to start from ")

type Buses struct {
	BusIds       []int
	BusPositions map[int]int
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)

	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	busTimes := parseBustimes(trimmedInput)
	candidateTime := busTimes.BusIds[0] * (*initialValue / busTimes.BusIds[0])
	currentTries := 0
	foundNumbers := -1
	currentIncrement := 1

	for {
		currentTries++
		candidateTime += currentIncrement
		found := true
		for i := range busTimes.BusIds {
			currentBusID := busTimes.BusIds[i]
			if (candidateTime+busTimes.BusPositions[currentBusID])%currentBusID != 0 {
				found = false
				break
			} else if i > foundNumbers {
				fmt.Printf("Increments: %d\n", currentIncrement)
				fmt.Println("Tries", currentTries)
				currentIncrement *= currentBusID
				foundNumbers++
			}
		}
		if found {
			break
		}
	}
	for i := range busTimes.BusIds {
		currentBusID := busTimes.BusIds[i]
		position := busTimes.BusPositions[currentBusID]
		if (candidateTime+position)%(currentBusID) != 0 {
			fmt.Printf("Discarding time: %d\n", candidateTime)
			os.Exit(1)
		}
	}
	fmt.Printf("Increment: %d\n", currentIncrement)
	fmt.Println("Tries", currentTries)
	fmt.Printf("Part 2: %d", candidateTime)

}

func parseBustimes(info []string) Buses {

	r := Buses{}

	var busIds []int
	busPositions := make(map[int]int)
	maxBusID := 0
	for k, busID := range strings.Split(info[1], ",") {
		if busID == "x" {
			continue
		}
		parsedBusID, _ := strconv.Atoi(busID)
		busPositions[parsedBusID] = k
		busIds = append(busIds, parsedBusID)
		if parsedBusID > maxBusID {
			maxBusID = parsedBusID
		}
	}

	r.BusIds = busIds
	r.BusPositions = busPositions
	return r

}
