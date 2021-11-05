package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sjnam/dlx"
)

func main() {
	d, err := dlx.NewDLX(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range d.Dance() {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range sol {
			fmt.Println(opt)
		}
		fmt.Println()
	}
}
