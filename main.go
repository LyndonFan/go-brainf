package main

import (
	"flag"
	"fmt"
)

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

	inputConfig := InputConfig{
		RequestInputChannel: &requestInputChannel,
		InChannel:           &inChannel,
		InputFileLocation:   inputFileLocation,
		ReadInputAsString:   inputAsString,
	}
	outputConfig := OutputConfig{
		OutputChannel:      &outChannel,
		OutputFileLocation: outputFileLocation,
		OutputAsString:     outputAsString,
	}

	if inputFileLocation == "" {
		go takeInputs(inputConfig)
	} else {
		go readInputs(inputConfig)
	}
	go runBrainFuck(program, requestInputChannel, inChannel, outChannel)
	printOutputs(outputConfig)
}
