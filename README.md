# Go-BrainF

A parser for the [brainf\*ck](https://en.wikipedia.org/wiki/Brainfuck) programming language.

## Quickstart

`./go-brainf -program "<PROGRAM>"`

It will then print the `<PROGRAM>`. If it needs input, it will pause and ask you for it. All inputs should be integers, all outputs are bytes / integers from 0 to 255 (inclusive) and printed immediately. The quotation marks are to ensure the terminal reads the program as an argument, not for the terminal.

Depending on your platform, replace `go-brainf` below with:

-   Linux (tested on Ubuntu, not sure about others): `go-brainf.sh`
-   MacOS: `go-brainf`
-   Windows: TBD

## Example programs

-   cat: asks for a number, prints it: `",."`
-   addition: `",>,<[->>+<<]>[->+<]>."`
-   multiplication: `",>,<[->[->+>+<<]>[-<+>]<<]>>>."` (Also in [program.txt](program.txt))
-   first n Fibonacci numbers (mod 256): `">>+<<,[->>.<[->>+<<]>[-<+>>+<]>[-<+>]<<<]"`

## Usage

```bash
Usage of ./go-brainf:
  -input string
        File location of inputs
  -input-as-string
        Whether to input is read as string (default true)
  -output string
        File location of outputs
  -output-as-string
        Whether to output result as string
  -program string
        Program to run
  -program-file string
        File location of program to run
```
