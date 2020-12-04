package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Passport struct {
	byr string
	iyr string
	eyr string
	hgt string
	hcl string
	ecl string
	pid string
	cid string
}

var ErrInvalidPassport error = errors.New("this passport is invalid")

var r *regexp.Regexp = regexp.MustCompile(`^(?P<height>[0-9]{2,3})(cm|in)$`)
var rHCL *regexp.Regexp = regexp.MustCompile(`#[0-9abcdef]{6}$`)
var rEcl *regexp.Regexp = regexp.MustCompile(`^(amb|blu|brn|gry|grn|hzl|oth)$`)
var rPid *regexp.Regexp = regexp.MustCompile(`^[0-9]{9}$`)

const minRequiredFields = 7

func validateIntegerInRange(s string, min, max int) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if i < min || i > max {
		return false
	}
	return true
}
func lineParser(line []byte) (*Passport, error) {
	tokens := bytes.Split(line, []byte(" "))
	p := Passport{}
	checkSumFields := 0
	for _, token := range tokens {
		fieldAndValue := bytes.Split(token, []byte(":"))
		if len(fieldAndValue) != 2 {
			continue
		}
		field := fieldAndValue[0]
		v := string(fieldAndValue[1])
		switch string(field) {
		case "byr":
			if len(v) > 4 {
				return nil, ErrInvalidPassport
			}
			if !validateIntegerInRange(v, 1920, 2002) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "iyr":
			if len(v) > 4 {
				return nil, ErrInvalidPassport
			}
			if !validateIntegerInRange(v, 2010, 2020) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "eyr":
			if len(v) > 4 {
				return nil, ErrInvalidPassport
			}

			if !validateIntegerInRange(v, 2020, 2030) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "hgt":
			if !r.MatchString(v) {
				return nil, ErrInvalidPassport
			}
			m := r.FindStringSubmatch(v)

			if len(m) != 3 {
				return nil, ErrInvalidPassport
			}

			if m[2] == "in" && !validateIntegerInRange(m[1], 59, 76) {
				return nil, ErrInvalidPassport

			} else if m[2] == "cm" && !validateIntegerInRange(m[1], 150, 193) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "hcl":
			if !rHCL.MatchString(v) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "ecl":
			if !rEcl.MatchString(v) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "pid":
			if !rPid.MatchString(v) {
				return nil, ErrInvalidPassport
			}
			checkSumFields++
		case "cid":

		}
	}
	if checkSumFields != minRequiredFields {
		return nil, ErrInvalidPassport
	}
	return &p, nil
}
func main() {
	f, _ := os.Open("input")
	s := bufio.NewScanner(f)

	passports := make([][]byte, 0, 100)
	passportLine := make([]byte, 0, 100)
	passportsRead := 0
	maxPassportLength := 0
	for s.Scan() {
		b := s.Bytes()
		if len(b) == 0 {
			passportsRead++
			//Emptyline, new passport
			newLine := make([]byte, len(passportLine), len(passportLine))
			copy(newLine, passportLine)
			passports = append(passports, newLine)
			passportLine = make([]byte, 0, 100)
		} else {
			passportLine = append(append(passportLine, ' '), b...)
		}
		if len(passportLine) > maxPassportLength {
			maxPassportLength = len(passportLine)
		}
	}
	//If there was any data left it's from a valid passport
	if len(passportLine) != 0 {
		passports = append(passports, passportLine)
		passportsRead++
	}
	validPassports := 0
	for _, line := range passports {
		_, err := lineParser(line)
		if err == nil {
			validPassports++
		}
	}
	fmt.Printf("Passports read %+v\n", passportsRead)
	fmt.Printf("Valid passports %d", validPassports)
}
