// Паттерн: graceful shutdown через context. Задача показывает, как остановить фоновую горутину по сигналу отмены или по внутреннему условию времени.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interruptSignalCh := make(chan os.Signal, 1)
	signal.Notify(interruptSignalCh, os.Interrupt)

	go printCurrentTimeUntilCanceled(ctx)
	go stopAtSecondMultipleOfFifteen(cancel)

	select {
	case <-interruptSignalCh:
		fmt.Println("Программа остановлена через Ctrl+C")
	case <-ctx.Done():
		fmt.Println("Программа остановлена по условию времени")
	}
}

func printCurrentTimeUntilCanceled(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Фоновый вывод времени остановлен")
			return
		default:
			fmt.Println(time.Now())
			time.Sleep(1 * time.Second)
		}
	}
}

func stopAtSecondMultipleOfFifteen(cancel context.CancelFunc) {
	for {
		if time.Now().Second()%15 == 0 {
			cancel()
			return
		}
	}
}
