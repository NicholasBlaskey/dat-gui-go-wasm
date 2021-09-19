# dat-gui-go-wasm

A wasm wrapper around dat.gui.js. Not complete and a work in progress.

### Using library

An example of how to use the library can be seen in the ```example``` directory.

In order to run the example first build the guiExample.go file.

```
GOOS=js GOARCH=wasm go build -o o.wasm guiExample.go
```

Then run the ```main.go``` web server.

```
go run main.go
```

Finally visit the webpage at. 

```
http://127.0.0.1:8080/
```

View the console and change the values in the slider to see how the values are being printed upon change.

