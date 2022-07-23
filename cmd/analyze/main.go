package main

import (
	"code"
	"log"
)

func main() {
	c, err := code.Analyze()
	if err != nil {
		log.Fatal(err)
	}

	c.PrintStructs()
	c.PrintFuncs()
}
