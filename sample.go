package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sjnam/go-dlx/dlx"
)

func main() {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	mcc := dlx.NewMCC()
	solStream, err := mcc.Dance(ctx, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range solStream {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range sol {
			fmt.Println(opt)
		}
		fmt.Println()
	}
}
