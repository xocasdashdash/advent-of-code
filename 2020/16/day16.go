package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Name string
	Attr string
	Min  int
	Max  int
}
type Ticket struct {
	Attrs []int
}
type TicketSystem struct {
	Rules         []Rule
	Attributes    []string
	MyTicket      *Ticket
	NearbyTickets []Ticket
}

func (ts TicketSystem) ValidateTickets() []int {
	var result []int
	invalidNumbers := make(map[int]bool, 10)
	hits := 0
	for _, t := range ts.NearbyTickets {
		for _, a := range t.Attrs {
			if _, ok := invalidNumbers[a]; ok {
				continue
			}
			valid := false
			for _, r := range ts.Rules {
				hits++
				if a >= r.Min && a <= r.Max {
					invalidNumbers[a] = true
					valid = true
					break
				}
			}
			if !valid {
				//fmt.Printf("Ticket at %d : %+v is not valid\n", ticketIndex, t)
				result = append(result, a)
			}
		}
	}
	fmt.Printf("Hit: %d\n", hits)
	return result
}

func (ts TicketSystem) FilterTicket(invalidNumbers []int) []Ticket {
	invalidNumbersMap := make(map[int]bool, len(invalidNumbers))
	for _, d := range invalidNumbers {
		invalidNumbersMap[d] = true
	}
	var result []Ticket
	for _, t := range ts.NearbyTickets {
		appendTicket := true
		for _, attr := range t.Attrs {
			if _, ok := invalidNumbersMap[attr]; ok {
				appendTicket = false
				break
			}
		}
		if appendTicket {
			result = append(result, t)
		}
	}
	return result
}

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

//
var attributeRegex = regexp.MustCompile(`^([A-Za-z ]+)+: ([0-9\-]+)* or ([0-9\-]+)*$`)
var ticketRegex = regexp.MustCompile(`[0-9]+,?`)
var columnRegex = regexp.MustCompile(`departure.*`)
var checks = 0

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	ts := parse(trimmedInput)
	invalidTickets := ts.ValidateTickets()
	fmt.Printf("Invalid tickets: %+v\n", invalidTickets)
	sum := 0
	for _, d := range invalidTickets {
		sum += d
	}
	ts.NearbyTickets = ts.FilterTicket(invalidTickets)
	//fmt.Printf("Leftover tickets: %d/%d\n", len(ts.NearbyTickets), len(ts.NearbyTickets)+len(invalidTickets))
	//fmt.Printf("Failure ratio: %d\n", sum)
	fmt.Printf("Part 1: %d\n", sum)
	columnMapping := ts.columnMapping()
	fmt.Printf("Column mapping: %#v\n", columnMapping)
	result := 1
	for attr, c := range columnMapping {
		if columnRegex.MatchString(attr) {
			result *= ts.MyTicket.Attrs[c]
		}
	}
	fmt.Printf("Part 2: %d\n", result)
}

func reducePossibleValues(
	result map[string]int,
	possibleValues map[string]map[int]bool,
	foundValues map[string]int) (map[string]int, map[string]map[int]bool, map[string]int) {

	loops := 0

	change := true
	for change {
		change = false
		loops++
		for attr := range foundValues {
			if foundValues[attr] == 1 {
				//fmt.Printf("Found a candidate column for attr %s. Possible : %+v\n", attr, possibleValues[attr])
				var foundColumn int = -1
				for candidateColumn, found := range possibleValues[attr] {
					if found {
						if foundColumn != -1 {
							//This is a safety check
							panic("impossible")
						}
						foundColumn = candidateColumn
					}
				}
				//Marking column foundColumn as not available anymore
				for a := range possibleValues {
					v := possibleValues[a][foundColumn]
					if v {
						change = true
						foundValues[a] = foundValues[a] - 1
					}
					possibleValues[a][foundColumn] = false
					checks++
				}
				result[attr] = foundColumn
				delete(foundValues, attr)
				delete(possibleValues, attr)
				//fmt.Printf("Candidate found at %d for attr %s\n", foundColumn, attr)
			}
		}
		//fmt.Printf("Setting %s to false on col %d\n", k, foundColumn)
	}
	fmt.Printf("Returning after %d loops \n", loops)
	return result, possibleValues, foundValues
}
func (ts TicketSystem) columnMapping() map[string]int {

	//This is basically a SUDOKU!
	possibleValues := make(map[string]map[int]bool, len(ts.Attributes))
	foundValues := make(map[string]int, len(ts.Attributes))
	result := make(map[string]int)
	for _, v := range ts.Attributes {
		possibleValues[v] = make(map[int]bool)
		foundValues[v] = len(ts.Attributes)
		for k := range ts.MyTicket.Attrs {
			possibleValues[v][k] = true
		}
	}
	// NOTE: THIS APPROACH DOES NOT WORK. AT SOME POINT YOU CAN HAVE MORE THAN ONE COLUMN MATCH MORE THAN
	// ONE VALUE
	// Idea 1:
	// Loop through attributes
	// For each attribute
	// Loop through every column and check if the rules for this attribute match every value in the column
	// If yes, add the column to the valid list
	// Idea 2:
	// Add all attributes that we haven't mapped to a list
	// Go through every value of a column applying every rule for every attribute, if a rule doesn't match remove the attribute as a candidate
	// If there's only one left, it must the correct attribute for that column.
	// unmappedAttrs := make([]string, len(ts.Attributes))
	attrRules := make(map[string][]Rule)
	for _, v := range ts.Rules {
		if _, ok := attrRules[v.Attr]; ok {
			attrRules[v.Attr] = append(attrRules[v.Attr], v)
		} else {
			attrRules[v.Attr] = []Rule{v}
		}
	}
	//fmt.Printf("All rules %+v\n%+v\n", attrRules, ts.Rules)

	columns := len(ts.Attributes)
	colValue := make(map[int][]int, len(ts.NearbyTickets))
	for _, t := range ts.NearbyTickets {
		for k, a := range t.Attrs {
			if _, ok := colValue[k]; ok {
				colValue[k] = append(colValue[k], a)
			} else {
				colValue[k] = []int{a}
			}
		}
	}
	for columnIndex := 0; columnIndex < columns; columnIndex++ {
		var d int

		for _, d = range colValue[columnIndex] {
			//fmt.Printf("Column : %d, row: %d, value: %d\n", columnIndex, row, d)
			for a, rules := range attrRules {
				if _, ok := result[a]; ok {
					continue
				}
				valid := false
				for _, r := range rules {
					checks++
					if d <= r.Max && d >= r.Min {
						valid = true
						//fmt.Printf("Column : %d, row: %d, value: %d is valid for attr %s: %+v\n", columnIndex, row, d, a, rules)
						break
					}
				}
				if !valid {
					possibleValues[a][columnIndex] = valid
					foundValues[a] = foundValues[a] - 1
					//fmt.Printf("Attr %s, col %d, \nFound %#v\nPossible %#v", a, columnIndex, foundValues, possibleValues)
				}
			}
			//fmt.Printf("Possible Values after col %d,row %d, %+v. Found values: %+v\n", columnIndex, row, possibleValues, foundValues)
		}
		//fmt.Printf("Found values: %+v\n", foundValues)
	}
	result, possibleValues, foundValues = reducePossibleValues(result, possibleValues, foundValues)

	fmt.Printf("Total number of checks: %d\n", checks)
	return result

}
func parse(s []string) TicketSystem {
	ts := TicketSystem{}
	ticketIndex := 0
	for _, l := range s {
		if attributeRegex.MatchString(l) {
			matches := attributeRegex.FindStringSubmatch(l)[1:]
			name := matches[0]
			ts.Attributes = append(ts.Attributes, name)
			for k, m := range matches[1:] {
				r := Rule{
					Attr: name,
					Name: fmt.Sprintf("%s-%d", name, k),
				}
				n := strings.Split(m, "-")
				minInt, _ := strconv.Atoi(n[0])
				maxInt, _ := strconv.Atoi(n[1])
				r.Min = minInt
				r.Max = maxInt
				ts.Rules = append(ts.Rules, r)
			}
		} else if ticketRegex.MatchString(l) {
			ticket := Ticket{}
			//fmt.Printf("Nearby ticket %d -> %s\n", ticketIndex, l)
			for _, attr := range strings.Split(l, ",") {
				d, _ := strconv.Atoi(attr)
				ticket.Attrs = append(ticket.Attrs, d)
			}
			if ts.MyTicket == nil {
				ts.MyTicket = &ticket
			} else {
				ts.NearbyTickets = append(ts.NearbyTickets, ticket)
				//fmt.Printf("Adding ticket %+v \n", ticket)
				ticketIndex++
			}

		}
	}
	return ts

}
