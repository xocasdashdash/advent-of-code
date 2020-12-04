package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
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

const minRequiredFields = 7

func lineParser(line []byte) (*Passport, error) {
	tokens := bytes.Split(line, []byte(" "))
	p := Passport{}
	requiredFields := 0
	for _, token := range tokens {
		fieldAndValue := bytes.Split(token, []byte(":"))
		if len(fieldAndValue) != 2 {
			continue
		}
		field := fieldAndValue[0]
		//v := fieldAndValue[1]
		switch string(field) {
		case "byr":
			requiredFields++
		case "iyr":
			requiredFields++
		case "eyr":
			requiredFields++
		case "hgt":
			requiredFields++
		case "hcl":
			requiredFields++
		case "ecl":
			requiredFields++
		case "pid":
			requiredFields++
		case "cid":

		}
	}
	if requiredFields != minRequiredFields {
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
			passportLine = append(passportLine, ' ')
			passportLine = append(passportLine, b...)
		}
	}
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
	fmt.Printf("Valid passports %d", validPassports)
}
