package datGUI

import (
	"reflect"
	"syscall/js"
)

type GUI struct {
	JSGUI js.Value
}

func New() *GUI {
	return &GUI{JSGUI: js.Global().Get("dat").Get("GUI").New()}
}

func (g *GUI) AddFolder(name string) *GUI {
	subJSGUI := g.JSGUI.Call("addFolder", name)
	return &GUI{JSGUI: subJSGUI}
}

func (g *GUI) Open() {
	g.JSGUI.Call("open")
}

func (g *GUI) Close() {
	g.JSGUI.Call("close")
}

func (g *GUI) Add(obj interface{}, fieldName string) *controller {
	// Get our field value.
	structVal := reflect.Indirect(reflect.ValueOf(obj))
	fieldVal := structVal.FieldByName(fieldName)
	fieldType := fieldVal.Type()

	// Make our field into a javascript object.
	jsObjMap := make(map[string]interface{})
	switch fieldType.Name() {
	case "float32":
		jsObjMap[fieldName] = fieldVal.Float()
	case "float64":
		jsObjMap[fieldName] = fieldVal.Float()
	case "int":
		jsObjMap[fieldName] = fieldVal.Int()
	case "int64":
		jsObjMap[fieldName] = fieldVal.Int()
	case "bool":
		jsObjMap[fieldName] = fieldVal.Bool()
	case "string":
		jsObjMap[fieldName] = fieldVal.String()
	case "": // Func type.
		jsObjMap[fieldName] = js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				fieldVal.Call(nil)
				return nil
			})
	default:
		panic("Tried to set unsupported type of " +
			fieldType.Name() + " in object ")
	}

	jsObj := js.ValueOf(jsObjMap)

	// We need this changeFunction since we will need to set the value again with
	// the new javascript controller when calls happen of both .min and .max for
	// some reason.
	c := &controller{}
	c.JSController = g.JSGUI.Call("add", jsObj, fieldName)
	c.changeFunc = func(jsController js.Value) {
		jsController.Call("onChange", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				// Set the value from the Javascript object to our Go object.
				changed := jsObj.Get(fieldName)
				f := reflect.ValueOf(obj).Elem().FieldByName(fieldName)

				switch fieldType.Name() {
				case "float32":
					f.SetFloat(changed.Float())
				case "float64":
					f.SetFloat(changed.Float())
				case "int":
					f.SetInt(int64(changed.Int()))
				case "bool":
					f.SetBool(changed.Bool())
				case "string":
					f.SetString(changed.String())
				case "": // Func type. Do nothing since we don't need to assign anything back.
				default:
					panic("We should never reach this since we check the type above.")
				}

				if c.listenerFunc != nil {
					c.listenerFunc()
				}
				return nil
			}))
	}
	c.changeFunc(c.JSController)

	return c
}

type controller struct {
	JSController js.Value
	changeFunc   func(js.Value)
	listenerFunc func()
}

func (c *controller) Min(x int) *controller {
	c.JSController = c.JSController.Call("min", x)
	c.changeFunc(c.JSController)
	return c
}

func (c *controller) Max(x int) *controller {
	c.JSController = c.JSController.Call("max", x)
	c.changeFunc(c.JSController)
	return c
}

func (c *controller) Step(x int) *controller {
	c.JSController = c.JSController.Call("max", x)
	c.changeFunc(c.JSController)
	return c
}

func (c *controller) Name(x string) *controller {
	c.JSController = c.JSController.Call("name", x)
	c.changeFunc(c.JSController)
	return c
}

func (c *controller) OnChange(fun func()) *controller {
	c.listenerFunc = fun

	return c
}
