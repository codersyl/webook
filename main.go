package main

import "fmt"

func main() {
	var a any = nil
	b, ok := a.(int64)
	fmt.Println(b, ok)
}
