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
var debug = flag.Bool("d", false, "Set to enable debugging")
var maskRegex = regexp.MustCompile(`^mask = (.*)$`)
var memRegex = regexp.MustCompile(`^mem\[([0-9]*)\] = ([0-9]*)$`)

type Operation func(*VM)
type VM struct {
	Registries map[int]int
	BitMask    []byte
	Operation  []Operation
}

func (vm *VM) PrintRegistries() {
	if *debug {
		fmt.Printf("**********\n")
		for k, v := range vm.Registries {
			fmt.Printf("Address: %d, Value: %d\n", k, v)
		}
		fmt.Printf("**********\n")
	}

}
func (vm *VM) Run() {
	for _, o := range vm.Operation {
		o(vm)
		vm.PrintRegistries()
	}
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	vm := parse(trimmedInput, "v1")
	vm.Run()
	sum := 0
	for _, v := range vm.Registries {
		sum += v
	}
	fmt.Printf("Number of registries: %d\n", len(vm.Registries))
	fmt.Printf("Part1: %d\n", sum)
	vm = parse(trimmedInput, "v2")
	vm.Run()
	sum = 0
	for _, v := range vm.Registries {
		sum += v
	}
	fmt.Printf("Number of registries: %d\n", len(vm.Registries))
	fmt.Printf("Part2: %d\n", sum)

}
func v1(address []byte, value []byte) Operation {
	return func(vm *VM) {
		v, _ := strconv.Atoi(string(value))
		a, _ := strconv.Atoi(string(address))
		newValue := 0
		bit := 0
		for k, b := range vm.BitMask {
			if b == '1' {
				bit = 1
			} else if b == '0' {
				bit = 0
			} else {
				bit = v >> k & 0x01
			}
			add := bit * (2 << (k) / 2)
			newValue += add
		}
		vm.Registries[a] = newValue
	}
}
func v2(address []byte, value []byte) Operation {
	return func(vm *VM) {
		v, _ := strconv.Atoi(string(value))
		a, _ := strconv.Atoi(string(address))
		addresses := make([]int, 1, 100)
		var bit int
		for k, b := range vm.BitMask {
			if b == '1' {
				bit = 1
			} else if b == '0' {
				bit = a >> k & 0x01
				if bit == 0 {
					//If it's a 0 we don't need to add a new number
					continue
				}
			} else {
				currentAddresses := len(addresses)
				for i := 0; i < currentAddresses; i++ {
					add := addresses[i] + (2 << (k) / 2)
					addresses = append(addresses, add)
				}
				continue
			}
			for j := range addresses {
				add := bit * (2 << (k) / 2)
				addresses[j] += add
			}
		}
		for _, address := range addresses {
			vm.Registries[address] = v
		}
	}

}

func parse(s []string, version string) *VM {

	vm := VM{
		BitMask: make([]byte, 36),
	}
	vm.Registries = make(map[int]int)
	for _, instruction := range s {
		if maskRegex.MatchString(instruction) {
			//It's a mask
			matches := maskRegex.FindStringSubmatch(instruction)[1:]

			mask := []byte(matches[0])
			for i, j := 0, len(mask)-1; i < j; i, j = i+1, j-1 {
				mask[i], mask[j] = mask[j], mask[i]
			}
			vm.Operation = append(vm.Operation, func(mask []byte) Operation {
				return func(vm *VM) {
					vm.BitMask = mask
					//fmt.Printf("Bitmask: %s. Length: %d\n", string(vm.BitMask), len(vm.BitMask))
				}
			}(mask))
		} else if memRegex.MatchString(instruction) {
			//It's a memory instruction
			matches := memRegex.FindStringSubmatch(instruction)[1:]
			address := []byte(matches[0])
			value := []byte(matches[1])

			var o Operation
			if version == "v1" {
				o = v1(address, value)
			} else {
				o = v2(address, value)
			}
			vm.Operation = append(vm.Operation, o)
		}
	}

	return &vm
}
