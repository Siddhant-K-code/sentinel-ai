package main

import (
	"log"

	"github.com/Siddhant-K-code/sentinel-ai/internal/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		log.Fatal(err)
	}
}
