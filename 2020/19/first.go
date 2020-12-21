package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var ruleToCheck = flag.String("r", "0", "Use this rule to check")

//Rule A rule.
type Rule struct {
	ruleID             string
	LRule              []*Rule
	RRule              []*Rule
	OriginalExpression string
	value              string
	parsed             bool
}

//MatchBranch recursively iterates through a string and validates a rule
func (r Rule) MatchBranch(rList []*Rule, input string) (bool, string) {
	matches := false
	pendingMatchText := input
	for _, localRule := range rList {
		//fmt.Printf("Applying rule(%d) - '%s' for input '%s'\n", localRule.ruleID, localRule.String(), inputToMatch)
		matches, pendingMatchText = localRule.MatchText(pendingMatchText)
		if !matches {
			return matches, ""
		}
	}
	return matches, pendingMatchText
}

var maxDepthForRule = make(map[int]map[string]int, 10)

//MatchText Matches text
func (r Rule) MatchText(input string) (matches bool, leftOverString string) {

	if len(r.LRule) > 0 {
		matches, leftOverString = r.MatchBranch(r.LRule, input)
		if matches {
			return
		}
	}
	if len(r.RRule) > 0 {
		matches, leftOverString = r.MatchBranch(r.RRule, input)
		if matches {
			return
		}
	}

	if len(r.RRule) == 0 && len(r.LRule) == 0 && input != "" {
		matches = input[:1] == r.value
		leftOverString = input[1:]
	}
	return

}
func (r Rule) regexWithLimit(depth int, maxDepth int) string {

	if len(r.LRule) == 0 {
		if len(r.value) == 0 {
			panic(fmt.Sprintf("rule %d is not well generated", r.ruleID))
		}
		return r.value
	}
	if depth >= maxDepth {
		return ""
	}
	var chars strings.Builder
	chars.WriteString("(?:")
	depth++
	for _, aRule := range r.LRule {
		chars.WriteString(aRule.regexWithLimit(depth, maxDepth))
	}
	if len(r.RRule) > 0 {
		chars.WriteString("|")
		for _, aRule := range r.RRule {
			chars.WriteString(aRule.regexWithLimit(depth, maxDepth))
		}
	}
	chars.WriteString(")")
	return chars.String()

}

//MythicalInformationSystem The information received from the elves
type MythicalInformationSystem struct {
	Rules    map[string]*Rule
	Messages []string
}

func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	system := parse(trimmedInput)
	validMessages := 0
	system = parse(trimmedInput)
	P1 := parseRule(system.Rules, *ruleToCheck)
	for _, m := range system.Messages {
		match, unmatchedChars := P1.MatchText(m)
		if match && len(unmatchedChars) == 0 {
			validMessages++
		}
	}
	fmt.Printf("Part1: %d\n", validMessages)
	validMessages = 0
	system = parse(trimmedInput)
	system.Rules["8"] = &Rule{
		ruleID:             "8",
		OriginalExpression: "42 | 42 8",
	}
	system.Rules["11"] = &Rule{
		ruleID:             "11",
		OriginalExpression: "42 31 | 42 11 31",
	}
	ruleP2 := parseRule(system.Rules, *ruleToCheck)
	//13 is the minimun value from experimentation
	testMatcher := regexp.MustCompile(fmt.Sprintf("^%s$", ruleP2.regexWithLimit(0, 13)))
	part2Regex := 0
	for _, m := range system.Messages {
		match, unmatchedChars := ruleP2.MatchText(m)
		if match && len(unmatchedChars) == 0 {
			validMessages++
		}
		if testMatcher.MatchString(m) {
			part2Regex++
		}
	}
	fmt.Printf("Part2(recursive): %d\n", validMessages)
	fmt.Printf("Part2(part2Regex): %d\n", part2Regex)
}

func parseRuleList(parsedRules map[string]*Rule, rulesID string) []*Rule {
	rIds := strings.Split(rulesID, " ")
	var result []*Rule
	for _, r := range rIds {
		if r == "" {
			continue
		}
		result = append(result, parseRule(parsedRules, r))
	}
	return result
}
func parseRule(parsedRules map[string]*Rule, ruleID string) *Rule {

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
		rule.ruleID = ruleID
		rule.value = strings.Replace(r.OriginalExpression, "\"", "", -1)
		return rule
	}
	sides := strings.Split(r.OriginalExpression, "|")
	rule.LRule = parseRuleList(parsedRules, sides[0])
	if len(sides) > 1 {
		rule.RRule = parseRuleList(parsedRules, sides[1])
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
	readRules := make(map[string]*Rule)
	for _, r := range ruleInput {
		splitRules := strings.Split(r, ":")
		readRules[splitRules[0]] = &Rule{
			ruleID:             splitRules[0],
			OriginalExpression: strings.TrimSpace(splitRules[1]),
		}
	}

	return MythicalInformationSystem{
		Rules:    readRules,
		Messages: messageInput,
	}
}
