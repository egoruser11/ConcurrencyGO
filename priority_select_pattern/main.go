// Паттерн: приоритетное чтение через вложенный select. Задача показывает, как сначала пытаться читать из более приоритетных каналов, а затем переходить к менее приоритетным.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	highPriorityCh := make(chan int)
	mediumPriorityCh := make(chan int)
	lowPriorityCh := make(chan int)

	go startPriorityGenerator(highPriorityCh, 200*time.Millisecond)
	go startPriorityGenerator(mediumPriorityCh, 500*time.Millisecond)
	go startPriorityGenerator(lowPriorityCh, 800*time.Millisecond)

	readWithPriority(ctx, highPriorityCh, mediumPriorityCh, lowPriorityCh)
}

func readWithPriority(ctx context.Context, highPriorityCh, mediumPriorityCh, lowPriorityCh <-chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Чтение с приоритетом завершено")
			return
		case value := <-highPriorityCh:
			fmt.Println("Высокий приоритет:", value)
		default:
			select {
			case value := <-mediumPriorityCh:
				fmt.Println("Средний приоритет:", value)
			default:
				select {
				case value := <-lowPriorityCh:
					fmt.Println("Низкий приоритет:", value)
				default:
					select {
					case <-ctx.Done():
						fmt.Println("Чтение с приоритетом завершено")
						return
					case value := <-highPriorityCh:
						fmt.Println("Высокий приоритет:", value)
					case value := <-mediumPriorityCh:
						fmt.Println("Средний приоритет:", value)
					case value := <-lowPriorityCh:
						fmt.Println("Низкий приоритет:", value)
					}
				}
			}
		}
	}
}

func startPriorityGenerator(output chan<- int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		output <- rand.Intn(100) + 1
	}
}
