// Паттерн: конвейер с context и фильтрацией. Задача показывает, как безопасно строить цепочку генерации, преобразования, фильтрации и финального сбора данных.
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	numbers := []int{1, 2, 2, 12, 24, 27}

	numberStream := generateSequence(ctx, numbers)
	multipliedStream := multiplyValues(ctx, numberStream, 10)
	filteredStream := filterValues(ctx, multipliedStream, isNotGreaterThanTwenty)
	result := collectValues(ctx, filteredStream)

	fmt.Println(result)
}

func generateSequence(ctx context.Context, numbers []int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for _, number := range numbers {
			select {
			case <-ctx.Done():
				return
			case output <- number:
			}
		}
	}()

	return output
}

func multiplyValues(ctx context.Context, input <-chan int, factor int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for value := range input {
			select {
			case <-ctx.Done():
				return
			case output <- value * factor:
			}
		}
	}()

	return output
}

func filterValues(ctx context.Context, input <-chan int, predicate func(int) bool) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for value := range input {
			if !predicate(value) {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case output <- value:
			}
		}
	}()

	return output
}

func collectValues(ctx context.Context, input <-chan int) []int {
	result := make([]int, 0)

	for {
		select {
		case <-ctx.Done():
			return result
		case value, isOpen := <-input:
			if !isOpen {
				return result
			}
			result = append(result, value)
		}
	}
}

func isNotGreaterThanTwenty(value int) bool {
	return value <= 20
}
