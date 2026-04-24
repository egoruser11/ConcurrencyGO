// Паттерн: базовый pipeline. Задача показывает минимальный конвейер из генератора, обработчика и потребителя данных.
package main

import "fmt"

func main() {
	sourceCh := make(chan int)
	resultCh := make(chan int)

	go generateNumbers(sourceCh, 1, 10)
	go doubleStream(sourceCh, resultCh)

	for result := range resultCh {
		fmt.Println(result)
	}
}

func generateNumbers(output chan<- int, start, end int) {
	defer close(output)

	for current := start; current <= end; current++ {
		output <- current
	}
}

func doubleStream(input <-chan int, output chan<- int) {
	defer close(output)

	for value := range input {
		output <- value * 2
	}
}
