package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
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
				return nil, fmt.Errorf("unmatched ']' at position %d", i)
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
		return nil, fmt.Errorf("unmatched '[' at position %d", stack[0])
	}
	if len(stack) > 1 {
		return nil, fmt.Errorf("unmatched '[' at positions %v", stack)
	}
	return res, nil
}

func runBrainFuck(code string, requestInputChannel chan bool, inChannel chan byte, outChannel chan byte) error {
	defer close(requestInputChannel)
	defer close(outChannel)
	lookup, err := preprocess(code)
	if err != nil {
		return err
	}
	const TAPE_LENGTH = 30000
	tape := make([]byte, TAPE_LENGTH)
	ptr := 0
	for i := 0; i < len(code); i++ {
		// fmt.Printf("%v", code[i])
		switch code[i] {
		case '>':
			ptr++
			if ptr >= TAPE_LENGTH {
				return fmt.Errorf("pointer went too far right and exceeded tape length of %d", TAPE_LENGTH)
			}
		case '<':
			ptr--
			if ptr < 0 {
				return fmt.Errorf("pointer went too far left")
			}
		case '+':
			tape[ptr]++
		case '-':
			tape[ptr]--
		case '.':
			outChannel <- tape[ptr]
		case ',':
			requestInputChannel <- true
			val, more := <-inChannel
			if !more {
				return fmt.Errorf("expected more inputs, but input channel closed")
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

func takeInputs(requestInputChannel chan bool, inChannel chan byte) error {
	defer close(inChannel)
	var input string
	for {
		_, more := <-requestInputChannel
		if !more {
			break
		}
		fmt.Print("Enter an integer (to be converted into a byte)\n>> ")
		fmt.Scanln(&input)
		x, err := strconv.Atoi(input)
		// fmt.Printf("%v %v %v\n", input, x, err)
		if err != nil {
			fmt.Printf("Unable to convert %v into an int, stopped taking inputs.\n", input)
			return err
		}
		inChannel <- byte(x)
	}
	return nil
}

func printOutputs(outChannel chan byte, outputAsString bool) error {
	for {
		b, more := <-outChannel
		if !more {
			break
		}
		if outputAsString {
			fmt.Printf("%c", b)
		} else {
			fmt.Println(b)
		}
	}
	fmt.Println()
	return nil
}

func readInputs(requestInputChannel chan bool, inChannel chan byte, inputFileLocation string) error {
	defer close(inChannel)
	file, err := os.OpenFile(inputFileLocation, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file %v does not exist", inputFileLocation)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	
	reader := bufio.NewReader(file)
	for {
		_, more := <-requestInputChannel
		if !more {
			break
		}
		input, err := reader.ReadByte()
		if err != nil {
			return err
		}
		inChannel <- input
	}
	return nil
}

func main() {
	var program, inputFileLocation, outputFileLocation string
	var inputAsString, outputAsString bool
	flag.StringVar(&program, "program", "", "Program to run")
	flag.StringVar(&inputFileLocation, "input", "", "Input file location")
	flag.StringVar(&outputFileLocation, "output", "", "Output file location")
	flag.BoolVar(&inputAsString, "input-as-string", true, "Whether to input is read as string")
	flag.BoolVar(&outputAsString, "output-as-string", false, "Whether to output result as string")
	flag.Parse()

	fmt.Printf("Program: %v\n", program)
	fmt.Printf("Input file: %v\n", inputFileLocation)
	fmt.Printf("Output file: %v\n", outputFileLocation)

	if inputFileLocation == outputFileLocation && outputFileLocation != "" {
		fmt.Printf("Output file (%v) cannot be the same as input file (%v).\n", outputFileLocation, inputFileLocation)
		return
	}

	requestInputChannel := make(chan bool, 1)
	inChannel := make(chan byte, 1)
	outChannel := make(chan byte, 1)

	// example program that prints first n Fibonacci numbers
	// ">>+<<,[->>.<[->>+<<]>[-<+>>+<]>[-<+>]<<<]"
	if inputFileLocation == "" {
		go takeInputs(requestInputChannel, inChannel)
	} else {
		go readInputs(requestInputChannel, inChannel, inputFileLocation)
	}
	go runBrainFuck(program, requestInputChannel, inChannel, outChannel)
	printOutputs(outChannel, outputAsString)
}
