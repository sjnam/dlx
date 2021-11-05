package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sjnam/go-dlx/dlx"
)

func main() {
	dx, err := dlx.NewDancer(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range dx.Dance() {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range sol {
			fmt.Println(opt)
		}
		fmt.Println()
	}
}
