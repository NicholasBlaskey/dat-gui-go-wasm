package main

import (
	"fmt"
	"syscall/js"

	"github.com/nicholasblaskey/dat-gui-go-wasm/datGUI"
)

type testType struct {
	X   int
	Y   bool
	Z   float32
	W   string
	Fun func()
}

func main() {
	// Create struct we want to bind to our GUI.
	obj := testType{1, true, 3, "String",
		func() {
			js.Global().Call("alert", "alerted")
		},
	}

	// Create GUI and a function we want to go off everytime a value changes.
	gui := datGUI.New()
	printObject := func() {
		fmt.Println("Printing object", obj)
	}

	// Add our fields with different options that can be applied.
	gui.Add(&obj, "X").Min(1).Max(10).Step(3).Name("x").OnChange(printObject)
	gui.Add(&obj, "Y").Name("This value is y").OnChange(printObject)
	gui.Add(&obj, "Z").Min(100).Max(1000).OnChange(printObject)
	gui.Add(&obj, "W").Name("This value is a string")

	// Add a folder and add a field to it..
	folder := gui.AddFolder("folder")
	folder.Open()
	folder.Add(&obj, "W").Name("Function in a folder")

	<-make(chan bool) // Prevent program from exiting
}
