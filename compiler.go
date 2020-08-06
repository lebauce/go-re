package gore

/*
#cgo LDFLAGS: -lz -lm -ldl -lncurses -lLLVM-9 -lclangTooling -lclangFrontendTool -lclangFrontend -lclangDriver -lclangSerialization -lclangCodeGen -lclangParse -lclangSema -lclangStaticAnalyzerFrontend -lclangStaticAnalyzerCheckers -lclangStaticAnalyzerCore -lclangAnalysis -lclangARCMigrate -lclangRewrite -lclangRewriteFrontend -lclangEdit -lclangAST -lclangLex -lclangBasic -lclang
#cgo CPPFLAGS: -I/usr/include -D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS

#include <stdio.h>
#include <stdlib.h>
#include "shim.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type EBPFCompiler struct {
	compiler *C.struct_bpf_compiler

	verbose       bool
	defaultCflags []string
}

func (e *EBPFCompiler) CompileToObjectFile(inputFile, outputFile string, cflags []string) error {
	inputC := C.CString(inputFile)
	defer C.free(unsafe.Pointer(inputC))

	outputC := C.CString(outputFile)
	defer C.free(unsafe.Pointer(outputC))

	cflagsC := make([]*C.char, len(e.defaultCflags)+len(cflags)+1)
	for i, cflag := range e.defaultCflags {
		cflagsC[i] = C.CString(cflag)
	}
	for i, cflag := range cflags {
		cflagsC[i] = C.CString(cflag)
	}
	cflagsC[len(cflagsC)-1] = nil

	defer func() {
		for _, cflag := range cflagsC {
			if cflag != nil {
				C.free(unsafe.Pointer(cflag))
			}
		}
	}()

	verboseC := C.char(0)
	if e.verbose {
		verboseC = 1
	}

	if err := C.bpf_compile_to_object_file(e.compiler, inputC, outputC, (**C.char)(&cflagsC[0]), verboseC); err != 0 {
		errs := C.GoString(C.bpf_compiler_get_errors(e.compiler))
		return errors.New(errs)
	}

	return nil
}

func (e *EBPFCompiler) Close() {
	runtime.SetFinalizer(e, nil)
	C.delete_bpf_compiler(e.compiler)
	e.compiler = nil
}

func NewEBPFCompiler(verbose bool) *EBPFCompiler {
	ebpfCompiler := &EBPFCompiler{
		compiler: C.new_bpf_compiler(),
	}

	runtime.SetFinalizer(ebpfCompiler, func(e *EBPFCompiler) {
		e.Close()
	})

	return ebpfCompiler
}

func main() {
	compiler := NewEBPFCompiler(false)

	if err := compiler.CompileToObjectFile("/tmp/toto.c", "/tmp/toto.o", nil); err != nil {
		panic(err)
	}

	runtime.GC()
}
