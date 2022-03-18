package main

import (
	"log"

	"github.com/bekzourdk/orders/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

//a
