package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"syscall/js"
)

func editionMode(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		js.Global().Call("alert", "Invalid no of arguments passed")
		return nil
	}

	inputTextArea := args[0]

	if inputTextArea.Get("classList").Call("contains", "success").Bool() {
		inputTextArea.Get("classList").Call("remove", "success")
	}

	if inputTextArea.Get("classList").Call("contains", "error").Bool() {
		inputTextArea.Get("classList").Call("remove", "error")
	}

	inputTextArea.Get("classList").Call("add", "info")

	return nil
}

func format(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		js.Global().Call("alert", "Invalid no of arguments passed")
		return nil
	}

	document := js.Global().Get("document")
	inputTextArea := args[0]
	inputTextArea.Get("classList").Call("remove", "info")
	text := inputTextArea.Get("value").String()
	outputTextArea := document.Call("getElementById", "output")
	outputTextArea.Set("value", "")
	indentSelector := document.Call("getElementById", "indent-selector")
	indentCount, _ := strconv.Atoi(indentSelector.Get("value").String())
	indent := strings.Repeat(" ", indentCount)

	var output bytes.Buffer

	err := json.Indent(&output, []byte(text), "", indent)

	if err != nil {
		if inputTextArea.Get("classList").Call("contains", "success").Bool() {
			inputTextArea.Get("classList").Call("remove", "success")
		}

		inputTextArea.Get("classList").Call("add", "error")

		if len(text) > 0 {
			outputTextArea.Set("value", err.Error())
		}

		return nil
	}

	outputTextArea.Set("value", output.String())

	if inputTextArea.Get("classList").Call("contains", "error").Bool() {
		inputTextArea.Get("classList").Call("remove", "error")
	}

	inputTextArea.Get("classList").Call("add", "success")

	return nil
}

func main() {
	js.Global().Set("format", js.FuncOf(format))
	js.Global().Set("editionMode", js.FuncOf(editionMode))

	c := make(chan int)

	<-c
}
