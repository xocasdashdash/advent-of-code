package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

type BusTimes struct {
	MinTime int
	BusIds  []int
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)

	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	busTimes := parseBustimes(trimmedInput)
	fmt.Printf("%+v", busTimes)
	closestBusID := 999999999
	closestBusTime := 999999999
	for _, busID := range busTimes.BusIds {

		rest := busTimes.MinTime % busID
		if rest == 0 {
			fmt.Printf("Found exact match \n")
			closestBusID = busID
			break
		}
		m := busTimes.MinTime / busID
		nextBus := (m + 1) * (busID)
		//fmt.Printf("M: %d, Bid: %d, MinTime: %d, nexTime: %d\n", m, busID, busTimes.MinTime, nextBus)
		if nextBus < closestBusTime {
			closestBusID = busID
			closestBusTime = nextBus
		} else {
			fmt.Printf("Next Bus(%d) is later than current one (%d)\n", busID, closestBusID)
		}
	}
	fmt.Printf("Part 1: %d | %d | %d", closestBusID*(closestBusTime-busTimes.MinTime), closestBusTime, busTimes.MinTime)

}

func parseBustimes(info []string) BusTimes {

	r := BusTimes{}
	r.MinTime, _ = strconv.Atoi(info[0])

	var busIds []int
	for _, busID := range strings.Split(info[1], ",") {
		if busID == "x" {
			continue
		}
		parsedBusID, _ := strconv.Atoi(busID)
		busIds = append(busIds, parsedBusID)
	}
	r.BusIds = busIds
	return r

}
