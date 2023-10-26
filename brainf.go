package main

import (
	"fmt"
	"strconv"
	"sync"
)

func preprocess(code string) (map[int]int, error) {
	stack := make([]int, 0, len(code))
	res := make(map[int]int)
	var j int
	for i, c := range code {
		switch c {
		case '[':
			stack = append(stack, i)
		case ']':
			if len(stack) == 0 {
				return nil, fmt.Errorf("Unmatched ']' at position %d", i)
			}
			j = stack[len(stack)-1]
			res[i] = j
			res[j] = i
			stack = stack[:len(stack)-1]
		default:
			// do nothing
		}
	}
	if len(stack) == 1 {
		return nil, fmt.Errorf("Unmatched '[' at position %d", stack[0])
	}
	if len(stack) > 1 {
		return nil, fmt.Errorf("Unmatched '[' at positions %v", stack)
	}
	return res, nil
}

func runBrainFuck(code string, inChannel chan byte, outChannel chan byte) error {
	defer close(outChannel)
	lookup, err := preprocess(code)
	if err != nil {
		return err
	}
	const TAPE_LENGTH = 30000
	tape := make([]byte, TAPE_LENGTH)
	ptr := 0
	for i := 0; i < len(code); i++ {
		fmt.Println(code[i])
		switch code[i] {
		case '>':
			ptr++
			if ptr >= TAPE_LENGTH {
				return fmt.Errorf("Pointer went too far right and exceeded tape length of %d", TAPE_LENGTH)
			}
		case '<':
			ptr--
			if ptr < 0 {
				return fmt.Errorf("Pointer went too far left")
			}
		case '+':
			tape[ptr]++
		case '-':
			tape[ptr]--
		case '.':
			outChannel <- tape[ptr]
		case ',':
			val, more := <-inChannel
			if !more {
				return fmt.Errorf("Expected more inputs, but input channel closed")
			}
			tape[ptr] = val
		case '[':
			if tape[ptr] == 0 {
				i = lookup[i]
			}
		case ']':
			if tape[ptr] != 0 {
				i = lookup[i]
			}
		default:
			// do nothing
		}
	}
	return nil
}

func takeInputs(inChannel chan byte) {
	defer close(inChannel)
	var input string
	for {
		fmt.Println("Enter an integer (to be converted into a byte)")
		fmt.Scanln(&input)
		x, err := strconv.Atoi(input)
		fmt.Println(input, x, err)
		if err != nil {
			fmt.Printf("Unable to convert %v into an int, stopped taking inputs.\n", input)
			break
		}
		inChannel <- byte(x)
	}
}

func printOutputs(outChannel chan byte) {
	for {
		b, more := <-outChannel
		if !more {
			break
		}
		fmt.Print(b)
	}
	fmt.Println()
}

func main() {
	inChannel := make(chan byte, 1)
	outChannel := make(chan byte, 1)
	program := ",+."
	fmt.Println(program)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		takeInputs(inChannel)
		go printOutputs(outChannel)
		go runBrainFuck(program, inChannel, outChannel)
	}()
	wg.Wait()
}
