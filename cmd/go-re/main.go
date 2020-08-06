package main

import (
	"flag"
	"runtime"
	"strings"

	gore "github.com/lebauce/go-re"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func main() {
	var (
		input, output string
		includes      arrayFlags
		defines       arrayFlags
		warnings      arrayFlags
		verbose       bool
	)

	flag.StringVar(&input, "c", "-", "Path of the C file to be compiled. Defaults to standard input")
	flag.StringVar(&output, "o", "-", "Path of the file to be written. Defaults to standard output")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.Var(&includes, "I", "Include directory")
	flag.Var(&defines, "D", "Preprocessor define")
	flag.Var(&warnings, "W", "Warnings")
	flag.Parse()

	compiler := gore.NewEBPFCompiler(verbose)

	var cflags []string

	for _, includeDir := range includes {
		cflags = append(cflags, "-I"+includeDir)
	}

	for _, define := range defines {
		cflags = append(cflags, "-D"+define)
	}

	for _, warning := range warnings {
		cflags = append(cflags, "-W"+warning)
	}

	if err := compiler.CompileToObjectFile(input, output, cflags); err != nil {
		panic(err)
	}

	runtime.GC()
}
