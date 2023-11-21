package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var program, programFile, inputFileLocation, outputFileLocation string
	var inputAsString, outputAsString bool
	flag.StringVar(&program, "program", "", "Program to run")
	flag.StringVar(&programFile, "program-file", "", "File location of program to run")
	flag.StringVar(&inputFileLocation, "input", "", "File location of inputs")
	flag.StringVar(&outputFileLocation, "output", "", "File location of outputs")
	flag.BoolVar(&inputAsString, "input-as-string", true, "Whether to input is read as string")
	flag.BoolVar(&outputAsString, "output-as-string", false, "Whether to output result as string")
	flag.Parse()

	if (programFile == "" && program == "") || (programFile != "" && program != "") {
		fmt.Println("Must specify exactly one of -program or -program-file")
		return
	}

	if programFile != "" {
		text, err := os.ReadFile(programFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		program = string(text)
	}

	if programFile != "" {
		fmt.Printf("Program file: %v\n", programFile)
	} else {
		fmt.Printf("Program: %v\n", program)
	}
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
