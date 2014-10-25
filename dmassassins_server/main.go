package main

import (
	"fmt"
)

// Theres actually a main function here
func main() {
	fmt.Println("Load Config")
	LoadConfig()
	fmt.Println("Server")
	StartServer()

}
