//go:build wasm && js

package main

import (
	"fmt"

	"github.com/NuruProgramming/Nuru/evaluator"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/object"
	"github.com/NuruProgramming/Nuru/parser"

	"syscall/js"
)

func Read(args []js.Value) {
	code := args[0].String()

	if len(args) > 1 {
		// TODO: type check arg for array of strings
		stdinBufferJsArray := args[1]
		// fmt.Print(stdinBufferJsArray)
		// fmt.Print(stdinBufferJsArray.Type().String())

		bufferLength := stdinBufferJsArray.Length()
		evaluator.StinBuffer = make([]string, bufferLength)

		for i := 0; i < bufferLength; i++ {
			evaluator.StinBuffer[i] = stdinBufferJsArray.Index(i).String()
		}
	}

	jsOutputReceiverFunction := js.Global().Get("nuruOutputReceiver")

	env := object.NewEnvironment()

	l := lexer.New(code)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Kuna makosa yafuatayo:")
		jsOutputReceiverFunction.Invoke("Kuna makosa yafuatayo:", true)

		for _, msg := range p.Errors() {
			// fmt.Println("\t" + msg)
			jsOutputReceiverFunction.Invoke("\t"+msg, true)
		}

	}
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if evaluated.Type() != object.NULL_OBJ {
			jsOutputReceiverFunction.Invoke(evaluated.Inspect(), true)
		}
	}

}

func runCode(this js.Value, args []js.Value) interface{} {
	Read(args)
	return nil
}

func main() {
	fmt.Println("Go WASM initialized")
	js.Global().Set("runCode", js.FuncOf(runCode))
	<-make(chan bool)
}
