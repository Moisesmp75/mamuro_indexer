package main

import (
	"fmt"
	"mamuro_indexer/helpers"
	"time"
)

func main() {
	startTime := time.Now()
	path := "full\\of\\path"
	helpers.Indexer(path)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("El programa tardó %s en ejecutarse.\n", duration)
}