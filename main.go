package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

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

	if inputFileLocation == "" {
		go takeInputs(requestInputChannel, inChannel)
	} else {
		go readInputs(requestInputChannel, inChannel, inputFileLocation)
	}
	go runBrainFuck(program, requestInputChannel, inChannel, outChannel)
	printOutputs(outChannel, outputAsString)
}
