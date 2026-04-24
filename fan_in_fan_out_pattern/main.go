// Паттерн: fan-out + fan-in. Задача показывает, как раздать обработку нескольким параллельным пайплайнам и затем собрать результаты обратно в один поток.
package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx := context.Background()
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)
		for _, value := range []int{1, 2, 3, 4, 5, 6} {
			inputCh <- value
		}
	}()

	workerOutputs := distributeWork(inputCh, 3, func(value int) int {
		return value * 10
	})

	mergedOutput := mergeChannels(ctx, workerOutputs)
	for result := range mergedOutput {
		fmt.Println(result)
	}
}

func distributeWork(input <-chan int, workerCount int, transform func(int) int) []<-chan int {
	outputs := make([]<-chan int, 0, workerCount)
	for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
		outputs = append(outputs, startPipelineWorker(input, transform))
	}
	return outputs
}

func startPipelineWorker(input <-chan int, transform func(int) int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for value := range input {
			output <- transform(value)
		}
	}()

	return output
}

func mergeChannels(ctx context.Context, inputs []<-chan int) <-chan int {
	output := make(chan int)

	var mergeGroup sync.WaitGroup
	mergeGroup.Add(len(inputs))

	for _, input := range inputs {
		go func(source <-chan int) {
			defer mergeGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, isOpen := <-source:
					if !isOpen {
						return
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}(input)
	}

	go func() {
		mergeGroup.Wait()
		close(output)
	}()

	return output
}
