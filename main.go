package main // import "argc.in/shrt"

import (
	"log"

	"argc.in/shrt/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
