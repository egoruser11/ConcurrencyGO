// Паттерн: fan-in. Задача показывает, как объединить несколько каналов в один общий поток и централизованно обработать результаты.
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	firstSourceCh := make(chan int)
	secondSourceCh := make(chan int)
	thirdSourceCh := make(chan int)
	mergedResultCh := make(chan int)

	var workerGroup sync.WaitGroup
	workerGroup.Add(3)

	go writeRandomNumbers(1, firstSourceCh, &workerGroup)
	go writeRandomNumbers(2, secondSourceCh, &workerGroup)
	go writeRandomNumbers(3, thirdSourceCh, &workerGroup)

	var mergeGroup sync.WaitGroup
	mergeGroup.Add(3)
	go mergeAndSquareValues(firstSourceCh, mergedResultCh, &mergeGroup)
	go mergeAndSquareValues(secondSourceCh, mergedResultCh, &mergeGroup)
	go mergeAndSquareValues(thirdSourceCh, mergedResultCh, &mergeGroup)

	go func() {
		workerGroup.Wait()
		mergeGroup.Wait()
		close(mergedResultCh)
	}()

	for result := range mergedResultCh {
		fmt.Printf("Результат после объединения: %d\n", result)
	}
}

func writeRandomNumbers(writerID int, output chan<- int, workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()
	defer close(output)

	for iteration := 0; iteration < 3; iteration++ {
		randomValue := rand.Intn(9) + 1
		fmt.Printf("Писатель %d отправил: %d\n", writerID, randomValue)
		output <- randomValue
	}
}

func mergeAndSquareValues(input <-chan int, mergedOutput chan<- int, mergeGroup *sync.WaitGroup) {
	defer mergeGroup.Done()

	for value := range input {
		mergedOutput <- value * value
	}
}
