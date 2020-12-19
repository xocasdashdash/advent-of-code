package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Set to run for partB")

//Rule A rule.
type Rule struct {
	ruleID             int
	LRule              []*Rule
	RRule              []*Rule
	Expression         string
	OriginalExpression string
	chars              string
	parsed             bool
}

//MatchText recursively iterates through a string and validates a rule
var numberOfChecks int
var visitedRules = make(map[string]bool)

func (r Rule) MatchBranch(rList []*Rule, input string, lazyRules map[int]bool) (bool, string) {
	matches := true
	localBranchInput := input
	for _, localRule := range rList {
		//fmt.Printf("Applying rule(%d) - '%s' for input '%s'\n", localRule.ruleID, localRule.String(), inputToMatch)
		matches, localBranchInput = localRule.MatchText(localBranchInput, lazyRules)
		if !matches {
			break
		}
	}
	return matches, localBranchInput
}

//MatchText Matches text
func (r Rule) MatchText(input string, lazyRules map[int]bool) (bool, string) {
	matches := false
	if len(input) == 0 {
		return false, ""
	}
	numberOfChecks++
	var localBranchInput string
	if len(r.LRule) > 0 {
		matches, localBranchInput = r.MatchBranch(r.LRule, input, nil)
		if matches {
			return matches, localBranchInput
		}
	}

	if len(r.RRule) > 0 {
		matches = false
		matches, localBranchInput = r.MatchBranch(r.RRule, input, nil)
		return matches, localBranchInput
	}

	//fmt.Printf("Comparing '%s' to '%s' at position %d. LRules: %d, RRules: %d\n", input, r.chars, position, len(r.LRule), len(r.RRule))
	matches = input[:1] == r.chars
	return matches, input[1:]
}

func (r Rule) Children() []Rule {
	var result []Rule

	if len(r.LRule) > 0 {
		for _, l := range r.LRule {
			result = append(result, l.Children()...)
		}
	}
	if len(r.RRule) > 0 {
		for _, l := range r.RRule {
			result = append(result, l.Children()...)
		}
	}
	if len(r.LRule) == 0 {
		result = append(result, r)
	}
	return result
}
func (r Rule) String() string {

	if len(r.LRule) == 0 {
		if len(r.chars) == 0 {
			//panic(fmt.Sprintf("rule %d is not well generated", r.ruleID))
		}
		return r.chars
	}
	var chars strings.Builder
	chars.WriteString("(?:")
	for _, aRule := range r.LRule {
		chars.WriteString(fmt.Sprintf("%s", aRule))
	}
	if len(r.RRule) > 0 {
		chars.WriteString("|")
		for _, aRule := range r.RRule {
			chars.WriteString(fmt.Sprintf("%s", aRule))
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
	rule := parseRule(system.Rules, "0")
	fmt.Printf("Using rule: '%s'\n", rule)
	validRuleTest := 0
	for _, m := range system.Messages {
		//fmt.Printf("Testing %s\n", m)
		match, matchingChars := rule.MatchText(m, nil)
		if match {
			if len(matchingChars) == 0 {
				validRuleTest++
			} else {
				fmt.Printf("Rest of string '%s' from '%s' does not match\n", matchingChars, m)
			}
		}
	}
	specialRules := map[int]bool{
		8:  true,
		11: true,
	}
	validRuleTest = 0
	for _, m := range system.Messages {
		//fmt.Printf("Testing %s\n", m)
		match, matchingChars := rule.MatchText(m, specialRules)
		if match {
			if len(matchingChars) == 0 {
				validRuleTest++
			} else {
				fmt.Printf("Rest of string '%s' from '%s' does not match\n", matchingChars, m)
			}
		}
	}
	fmt.Printf("Part1(RuleTest): %d\n", validRuleTest)
	fmt.Printf("Checks: %d\n", numberOfChecks)
}

func parseRule(parsedRules map[string]*Rule, ruleID string) *Rule {

	r, ok := parsedRules[ruleID]
	if ok && r.parsed {
		return r
	}
	rid, _ := strconv.Atoi(ruleID)

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
		r := Rule{
			ruleID: rid,
			chars:  strings.Replace(r.OriginalExpression, "\"", "", -1),
			parsed: true,
		}
		r.Expression = r.String()
		return &r
	}
	rule := Rule{
		ruleID: rid,
		parsed: true,
	}
	sides := strings.Split(r.OriginalExpression, "|")
	for k, side := range sides {
		rIds := strings.Split(strings.TrimSpace(side), " ")
		if k == 0 {
			//Left side
			for _, r := range rIds {
				var parsedRule *Rule
				parsedRule, ok = parsedRules[r]
				if !ok || parsedRule.parsed == false {
					parsedRule = parseRule(parsedRules, r)
				}
				if parsedRule != nil {
					rule.LRule = append(rule.LRule, parsedRule)
				}
			}
		} else {
			//Right side
			for _, r := range rIds {
				var parsedRule *Rule
				parsedRule, ok = parsedRules[r]
				if !ok || parsedRule.parsed == false {
					parsedRule = parseRule(parsedRules, r)
				}
				if parsedRule != nil {
					rule.RRule = append(rule.RRule, parsedRule)
				}
			}
		}
	}
	rule.Expression = rule.String()

	return &rule
}
func parseRules(parsedRules map[string]*Rule) map[string]*Rule {

	if parsedRules == nil {
		parsedRules = make(map[string]*Rule)
	}
	for ruleID := range parsedRules {
		fmt.Printf("Parsing rule: %s\n", ruleID)
		r := parseRule(parsedRules, ruleID)
		r.parsed = true
		if r != nil {
			parsedRules[ruleID] = r
		}
	}
	return parsedRules
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
