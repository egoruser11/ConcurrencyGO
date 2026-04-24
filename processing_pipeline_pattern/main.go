// Паттерн: многоэтапный конвейер обработки данных. Задача показывает, как несколько воркеров последовательно обрабатывают один объект и передают результат дальше по каналам.
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type ProcessingItem struct {
	SequenceNumber int
	RandomValue    int
	SquaredValue   int
	SquareRoot     float64
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sourceCh := make(chan *ProcessingItem)
	squaredCh := make(chan *ProcessingItem)
	resultCh := make(chan ProcessingItem)

	go calculateSquare(sourceCh, squaredCh)
	go calculateSquareRoot(squaredCh, resultCh)
	go printProcessedItems(resultCh)

	generateItems(ctx, sourceCh)
}

func generateItems(ctx context.Context, output chan<- *ProcessingItem) {
	defer close(output)

	itemNumber := 1
	for {
		time.Sleep(500 * time.Millisecond)

		newItem := &ProcessingItem{
			SequenceNumber: itemNumber,
			RandomValue:    rand.Intn(100),
		}
		itemNumber++

		select {
		case <-ctx.Done():
			return
		case output <- newItem:
		}
	}
}

func calculateSquare(input <-chan *ProcessingItem, output chan<- *ProcessingItem) {
	defer close(output)

	for item := range input {
		item.SquaredValue = item.RandomValue * item.RandomValue
		output <- item
	}
}

func calculateSquareRoot(input <-chan *ProcessingItem, output chan<- ProcessingItem) {
	defer close(output)

	for item := range input {
		item.SquareRoot = math.Sqrt(float64(item.RandomValue))
		output <- *item
	}
}

func printProcessedItems(input <-chan ProcessingItem) {
	for item := range input {
		fmt.Println(item)
	}
}
