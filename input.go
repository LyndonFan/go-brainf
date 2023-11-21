package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type InputConfig struct {
	RequestInputChannel *chan bool
	InChannel           *chan byte
	InputFileLocation   string
	ReadInputAsString   bool
}

func takeInputs(inputConfig InputConfig) error {
	defer close(*inputConfig.InChannel)
	var input string
	for {
		_, more := <-*inputConfig.RequestInputChannel
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
		*inputConfig.InChannel <- byte(x)
	}
	return nil
}

func readInputs(config InputConfig) error {
	defer close(*config.InChannel)
	file, err := os.OpenFile(config.InputFileLocation, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file %v does not exist", config.InputFileLocation)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		_, more := <-*config.RequestInputChannel
		if !more {
			break
		}
		input, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*config.InChannel <- input
	}
	return nil
}
