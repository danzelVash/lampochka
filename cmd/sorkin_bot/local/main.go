package main

import (
	"context"
	"fmt"
	"log"
	"sorkin_bot/internal"
)

func main() {
	ctx := context.Background()
	fmt.Println("")
	log.Fatal(internal.NewApp(ctx).Run(ctx))
}
