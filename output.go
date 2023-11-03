package main

import "fmt"

type OutputConfig struct {
	OutputChannel      *chan byte
	OutputFileLocation string
	OutputAsString     bool
}

func printOutputs(outputConfig OutputConfig) error {
	for {
		b, more := <-*outputConfig.OutputChannel
		if !more {
			break
		}
		if outputConfig.OutputAsString {
			fmt.Printf("%c", b)
		} else {
			fmt.Println(b)
		}
	}
	fmt.Println()
	return nil
}
