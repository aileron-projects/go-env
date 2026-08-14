package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aileron-projects/go-env"
)

var (
	file = flag.String("env", "env.txt", "env file path to load")
)

func main() {
	flag.Parse()
	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}
	kvs, err := env.Load(*file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed values:", *file)
	fmt.Println("---------------------------")
	fmt.Println("Number of variables:", len(kvs))
	for k, v := range kvs {
		fmt.Println("")
		fmt.Println(">> KEY:", k)
		if strings.Contains(v, "\n") {
			fmt.Println(">> VALUE: ⬎")
			fmt.Println(v)
		} else {
			fmt.Println(">> VALUE:", v)
		}
	}
}
