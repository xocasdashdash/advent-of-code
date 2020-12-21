package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var ruleToCheck = flag.Int("r", 0, "Use this rule to check")
var partA = flag.Bool("pA", false, "Part A")

//Rule A rule.
type Rule struct {
	ruleID             int
	rules              [][]*Rule
	OriginalExpression string
	value              string
	parsed             bool
}

//MatchText Matches text and returns the number of matching characters
func (r Rule) MatchText(input string) (result []int) {

	if len(r.rules) == 0 {
		if len(input) < len(r.value) {
			result = nil
			return
		}
		if input[:len(r.value)] == r.value {
			result = []int{len(r.value)}
			return
		}
	}
	var matches []int
	for _, groupRule := range r.rules {
		potentialMatches := []int{0}
		for _, localRule := range groupRule {
			var newPotentialMatches []int
			for _, m := range potentialMatches {
				matches := localRule.MatchText(input[m:])
				for _, v := range matches {
					newPotentialMatches = append(newPotentialMatches, v+m)
				}
			}
			potentialMatches = newPotentialMatches
		}
		matches = append(matches, potentialMatches...)
	}
	result = matches
	return
}

// StringWithLimit returns a regex string with a maximum depth to avoid infinite recursion
func (r Rule) StringWithLimit(depth int, maxDepth int) string {

	if len(r.rules) == 0 {
		if len(r.value) == 0 {
			panic(fmt.Sprintf("rule %d is not well generated", r.ruleID))
		}
		return r.value
	}
	if depth >= maxDepth {
		return ""
	}
	var ruleStringBuilder strings.Builder
	ruleStringBuilder.WriteString("(?:")
	depth++
	for k, ruleGroup := range r.rules {
		if k > 0 {
			ruleStringBuilder.WriteString("|")
		}
		for _, aRule := range ruleGroup {
			ruleStringBuilder.WriteString(aRule.StringWithLimit(depth, maxDepth))
		}
	}
	ruleStringBuilder.WriteString(")")
	return ruleStringBuilder.String()
}

//MythicalInformationSystem The information received from the elves
type MythicalInformationSystem struct {
	Rules    map[int]*Rule
	Messages []string
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	system := parse(trimmedInput)
	validRuleTest := 0
	validRuleRegex := 0
	rule := parseRule(system.Rules, *ruleToCheck)
	testMatcherP1 := regexp.MustCompile(fmt.Sprintf("^%s$", rule.StringWithLimit(0, 13)))
	for _, m := range system.Messages {
		match := rule.MatchText(m)
		for _, r := range match {
			if r == len(m) {
				validRuleTest++
				break
			}
		}
		if testMatcherP1.MatchString(m) {
			validRuleRegex++
		}
	}
	fmt.Printf("Part1(RuleTest): %d\n", validRuleTest)
	fmt.Printf("Part2(RuleRegex): %d\n", validRuleRegex)
	validRuleTest = 0
	validRuleRegex = 0
	system = parse(trimmedInput)
	system.Rules[8] = &Rule{
		OriginalExpression: "42 | 42 8",
	}
	system.Rules[11] = &Rule{
		OriginalExpression: "42 31 | 42 11 31",
	}
	ruleP2 := parseRule(system.Rules, *ruleToCheck)
	testMatcher := regexp.MustCompile(fmt.Sprintf("^%s$", ruleP2.StringWithLimit(0, 13)))
	for _, m := range system.Messages {
		match := ruleP2.MatchText(m)
		for _, r := range match {
			if r == len(m) {
				validRuleTest++
				break
			}
		}
		m2 := testMatcher.MatchString(m)
		if m2 {
			validRuleRegex++
			//return
		}
	}
	fmt.Printf("Part2(validRuleTest): %d\n", validRuleTest)
	fmt.Printf("Part2(validRuleRegex): %d\n", validRuleRegex)

}

func parseRule(parsedRules map[int]*Rule, ruleID int) *Rule {

	r, ok := parsedRules[ruleID]
	if ok && r.parsed {
		return r
	}

	rule := &Rule{
		ruleID:             ruleID,
		parsed:             true,
		OriginalExpression: r.OriginalExpression,
	}
	parsedRules[ruleID] = rule
	if strings.Index(r.OriginalExpression, "\"") != -1 {
		rule.value = strings.Replace(r.OriginalExpression, "\"", "", -1)
		return rule
	}

	for _, s := range strings.Split(strings.TrimSpace(r.OriginalExpression), "|") {
		rIds := strings.Split(strings.TrimSpace(s), " ")
		var groupRules []*Rule
		for _, r := range rIds {
			d, _ := strconv.Atoi(r)
			groupRules = append(groupRules, parseRule(parsedRules, d))
		}
		rule.rules = append(rule.rules, groupRules)
	}

	return rule
}
func parse(input []string) MythicalInformationSystem {
	var i int
	var ruleInput []string
	var messageInput []string
	for i = 0; input[i] != ""; i++ {
		ruleInput = append(ruleInput, input[i])
	}
	messageInput = input[i+1:]
	readRules := make(map[int]*Rule)
	for _, r := range ruleInput {
		splitRules := strings.Split(r, ":")
		ruleID, _ := strconv.Atoi(splitRules[0])
		readRules[ruleID] = &Rule{
			ruleID:             ruleID,
			OriginalExpression: strings.TrimSpace(splitRules[1]),
		}
	}

	//fmt.Printf("R2: %s", readRules["2"].String())
	return MythicalInformationSystem{
		Rules:    readRules,
		Messages: messageInput,
	}
}
