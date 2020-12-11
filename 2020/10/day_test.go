package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestCalculate2(t *testing.T) {
	type input struct {
		file   string
		result int
	}
	inputs := []input{
		{
			file:   "testInput1",
			result: 8,
		},
		/*
			 {
				file:   "testInput2",
				result: 19208,
			},
		*/
	}
	for i := range inputs {

		input, _ := ioutil.ReadFile(inputs[i].file)

		output := calculate2(inputs[i].file, strings.Split(strings.TrimSpace(string(input)), "\n"))
		if output != inputs[i].result {
			t.Errorf("Output %d is not what is expected %d", output, inputs[i].result)
			t.Fail()
		}
	}
}
