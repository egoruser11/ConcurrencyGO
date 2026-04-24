// Паттерн: несколько продюсеров в один канал. Задача показывает, как несколько горутин одновременно пишут данные в общий канал, а главный поток безопасно дочитывает их до закрытия.
package main

import (
	"fmt"
	"sync"
)

func main() {
	numberCh := make(chan int)

	var producerGroup sync.WaitGroup
	producerGroup.Add(2)

	go writeRange(numberCh, 1, 5, &producerGroup)
	go writeRange(numberCh, 6, 100, &producerGroup)

	go func() {
		producerGroup.Wait()
		close(numberCh)
	}()

	for value := range numberCh {
		fmt.Println(value)
	}
}

func writeRange(output chan<- int, start, end int, producerGroup *sync.WaitGroup) {
	defer producerGroup.Done()

	for current := start; current <= end; current++ {
		output <- current
	}
}
