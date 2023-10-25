package main

import "fmt"

func preprocess(code string) (map[int]int, err) {
	stack := make([]int, 0, len(code))
	res := make(map[int]int)
	var j int
	for i,c := range(code){
		switch c {
			case '[':
				stack = append(stack, i)
			case ']':
				if len(stack)==0 {
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
	if len(stack)==1 {
		return nil, fmt.Errorf("Unmatched '[' at position %d", stack[0])
	}
	if len(stack)>1 {
		return nil, fmt.Errorf("Unmatched '[' at positions %v", stack)
	}
	return res, nil
}

func runBrainFuck(code string, inChannel chan byte, outChannel chan byte) error {
	lookup, err := preprocess(code)
	if err != nil {
		return err
	}
	const TAPE_LENGTH := 30000
	tape := make([]byte, TAPE_LENGTG)
	ptr := 0
	for i:=0;i<len(code);i++ {
		switch code[i] {
			case '>':
				ptr++
				if ptr >= TAPE_LENGTH {
					return fmt.Error
				}
			case '<':
				ptr--
			case '+':
				tape[ptr]++
			case '-':
				tape[ptr]--
			case '.':
				outChannel <- tape[ptr]
			case ',':
				tape[ptr] <- inChannel
			case '[':
				if tape[ptr]==0 {
					i = lookup[i]
				}
			case ']':
				if tape[ptr]!=0 {
					i = lookup[i]
				}
			default:
				// do nothing
		}
	}
	return nil
}

func main() {


}
knsr