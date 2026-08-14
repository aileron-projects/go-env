package main

import (
	"fmt"
	"os"

	_ "github.com/aileron-projects/go-env/autoload"
)

func main() {
	fmt.Println("FOO=", os.Getenv("FOO"))
	fmt.Println("BAR=", os.Getenv("BAR"))
	fmt.Println("BAZ=", os.Getenv("BAZ"))
}
