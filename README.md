# Juno Go

This is a library that provides you with helper methods for interfacing with the microservices framework, [juno](https://github.com/bytesonus/juno).

## How to use:

There is a lot of flexibility provided by the library, in terms of connection options and encoding protocol options. However, in order to use the library, none of that is required.

In case you are planning to implement a custom connection option, you will find an example in `connection/unix_socket_connection.go`.

For all other basic needs, you can get away without worrying about any of that.

### A piece of code is worth a thousand words

```go
package main

import (
    "fmt"
    juno "github.com/bytesonus/juno-go"
)

func main() {
	module := juno.Default("./path/to/juno.sock")
	channel, err := module.Initialize("module-name", "1.0.0", nil)
	if err != nil {
		panic(err)
	}
	_ = <-channel // Wait for initialize to complete
	fmt.Println("Initialized!")

	channel, err = module.DeclareFunction("printHello", func (args map[string]interface{}) interface{} {
		fmt.Println("Hello")
		return nil
	})
	if err != nil {
		panic(err)
	}
	_ = <-channel // Wait for declaration to complete
	fmt.Println("Declaration!")

	channel, err = module.CallFunction("module2.printHelloWorld", nil)
	if err != nil {
		panic(err)
	}
	response := <-channel // Wait for declaration to complete
	fmt.Println(response)
}
```
