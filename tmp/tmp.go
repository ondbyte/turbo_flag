package main

import "fmt"

func a(data interface{}) {
	c, ok := data.(map[string]string)
	if !ok {
		fmt.Println("error")
	} else {
		fmt.Println(c)
	}

}

func main() {
	x := make(map[any]any)
	y := make(map[interface{}]interface{})
	x["yadu"] = "nandan"
	y["yadu"] = "nandan"
	//running for x
	a(x)
	//running for y
	a(y)
}
