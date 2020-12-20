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

//Rule A rule.
type Rule struct {
	ruleID             int
	LRule              []*Rule
	RRule              []*Rule
	Expression         string
	OriginalExpression string
	chars              string
	parsed             bool
	parsing            bool
}

//MatchText recursively iterates through a string and validates a rule
var numberOfChecks int
var visitedRules = make(map[string]bool)

//MatchBranch recursively iterates through a string and validates a rule
func (r Rule) MatchBranch(currentDepth int, rList []*Rule, input string, lazyRules map[int]bool) (bool, string) {
	currentDepth++
	matches := false
	localBranchInput := input
	for _, localRule := range rList {
		//fmt.Printf("Applying rule(%d) - '%s' for input '%s'\n", localRule.ruleID, localRule.String(), inputToMatch)
		matches, localBranchInput = localRule.MatchText(currentDepth, localBranchInput, lazyRules)
		if !matches {
			break
		}
	}
	return matches, localBranchInput
}

var maxDepthForRule = make(map[int]map[string]int, 10)

//MatchText Matches text
func (r Rule) MatchText(currentDepth int, input string, lazyRules map[int]bool) (bool, string) {
	matches := false
	//defer func() {
	//	currentMaxDepth, ok := maxDepthForRule[r.ruleID]
	//	if !ok {
	//		maxDepthForRule[r.ruleID] = make(map[string]int, 1)
	//		currentMaxDepth = maxDepthForRule[r.ruleID]
	//	}
	//	if d, _ := currentMaxDepth["depth"]; d < currentDepth {
	//		currentMaxDepth["depth"] = d
	//	}
	//	currentMaxDepth["matches"] = 1
	//	if !matches {
	//		currentMaxDepth["matches"] = 0
	//	}
	//	//fmt.Printf("result %v at depth %d for rule %d\n", matches, currentDepth, r.ruleID)
	//}()

	numberOfChecks++
	var localBranchInput string
	var matchesL, matchesR bool
	if len(r.LRule) > 0 {
		matchesL, localBranchInput = r.MatchBranch(currentDepth, r.LRule, input, nil)
	}
	rLeftInput := localBranchInput
	if len(r.RRule) > 0 {
		matchesR, localBranchInput = r.MatchBranch(currentDepth, r.RRule, input, nil)
	}
	if matchesL {
		localBranchInput = rLeftInput
	}
	matches = matchesL || matchesR
	if len(input) == 0 {
		//fmt.Printf("not matching because no more input\n")
		return false, ""
	}
	if len(r.RRule) == 0 && len(r.LRule) == 0 {
		matches = input[:1] == r.chars
		localBranchInput = input[1:]
	}
	return matches, localBranchInput

}
func (r Rule) StringWithLimit(depth int, maxDepth int) string {

	if len(r.LRule) == 0 {
		if len(r.chars) == 0 {
			panic(fmt.Sprintf("rule %d is not well generated", r.ruleID))
		}
		return r.chars
	}
	if depth >= maxDepth {
		return ""
	}
	var chars strings.Builder
	chars.WriteString("(?:")
	depth++
	for _, aRule := range r.LRule {
		chars.WriteString(aRule.StringWithLimit(depth, maxDepth))
	}
	if len(r.RRule) > 0 {
		chars.WriteString("|")
		for _, aRule := range r.RRule {
			chars.WriteString(aRule.StringWithLimit(depth, maxDepth))
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
	//fmt.Printf("Read rules %#v", system.Rules)
	validRuleTest := 0
	validRuleRegex := 0
	rule := parseRule(system.Rules, "0")
	testMatcherP1 := regexp.MustCompile(fmt.Sprintf("^%s$", rule.StringWithLimit(0, 13)))
	fmt.Printf("Using rule: '%s'\n", rule)
	for _, m := range system.Messages {
		//fmt.Printf("Testing %s\n", m)
		match, matchingChars := rule.MatchText(0, m, nil)
		if match {
			if len(matchingChars) == 0 {
				validRuleTest++
			}
		}
		if testMatcherP1.MatchString(m) {
			validRuleRegex++
		}
	}
	fmt.Printf("Part1(RuleTest): %d\n", validRuleTest)
	fmt.Printf("Part2(RuleRegex): %d\n", validRuleRegex)
	validRuleTest = 0
	system = parse(trimmedInput)
	system.Rules["8"] = &Rule{
		ruleID:             8,
		OriginalExpression: "42 | 42 8",
	}
	system.Rules["11"] = &Rule{
		ruleID:             11,
		OriginalExpression: "42 31 | 42 11 31",
	}
	ruleP2 := parseRule(system.Rules, "0")
	fmt.Printf("Using rule %s\n", ruleP2.StringWithLimit(0, 9))
	testMatcher := regexp.MustCompile(fmt.Sprintf("^%s$", ruleP2.StringWithLimit(0, 13)))
	part2Valid := 0
	part2Regex := 0
	for _, m := range system.Messages {
		match, matchingChars := ruleP2.MatchText(0, m, nil)
		//fmt.Printf("Leftover chars %s\n", matchingChars)
		if match {
			if len(matchingChars) == 0 {
				//fmt.Printf("%s\n", m)
				part2Valid++
			}
		}
		m2 := testMatcher.MatchString(m)
		if m2 {
			part2Regex++
		}
		if match != m2 {
			//fmt.Printf("%s\n", m)
		}
	}
	fmt.Printf("Part2(RuleTest): %d\n", part2Valid)
	fmt.Printf("Part2(part2Regex): %d\n", part2Regex)
	fmt.Printf("Checks: %d\n", numberOfChecks)
}

func parseRuleList(parsedRules map[string]*Rule, rulesID string) []*Rule {
	rIds := strings.Split(rulesID, " ")
	var result []*Rule
	for _, r := range rIds {
		var parsedRule *Rule
		//var ok bool
		parsedRule, _ = parsedRules[r]
		if parsedRule != nil && parsedRule.parsed == false {
			parsedRule = parseRule(parsedRules, r)
		}
		if parsedRule != nil {
			result = append(result, parsedRule)
		}
	}
	return result
}
func parseRule(parsedRules map[string]*Rule, ruleID string) *Rule {
	//fmt.Printf("Going in for rule %s\n", ruleID)

	r, ok := parsedRules[ruleID]
	if ok && r.parsed {
		fmt.Printf("Rule %s, already parsed\n", ruleID)
		return r
	}
	rid, _ := strconv.Atoi(ruleID)
	rule := &Rule{
		ruleID:             rid,
		parsed:             true,
		OriginalExpression: r.OriginalExpression,
	}
	parsedRules[ruleID] = rule
	if strings.Index(r.OriginalExpression, "|") == -1 && strings.Index(r.OriginalExpression, " ") == -1 {
		_, err := strconv.Atoi(r.OriginalExpression)
		if err == nil {
			//If it's a single char we can ignore the rule and just map it to the reference
			var parsedRule *Rule
			var ok bool
			parsedRule, ok = parsedRules[r.OriginalExpression]
			if !ok || parsedRule.parsed == false {
				parsedRule = parseRule(parsedRules, r.OriginalExpression)
			}
			return &Rule{
				ruleID:             rid,
				LRule:              []*Rule{parsedRule},
				OriginalExpression: r.OriginalExpression,
				parsed:             true,
			}
		}
		rule.ruleID = rid
		rule.chars = strings.Replace(r.OriginalExpression, "\"", "", -1)
		//rule.Expression = r.String()
		return rule
	}
	sides := strings.Split(r.OriginalExpression, "|")
	rule.LRule = parseRuleList(parsedRules, sides[0])
	if len(sides) > 1 {
		rule.RRule = parseRuleList(parsedRules, sides[1])
	}
	//rule.Expression = rule.String()

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
		ruleID, _ := strconv.Atoi(splitRules[0])
		readRules[splitRules[0]] = &Rule{
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
