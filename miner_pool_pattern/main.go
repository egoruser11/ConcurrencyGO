// Паттерн: пул воркеров с отменой через context. Задача показывает, как несколько шахтёров параллельно добывают ресурс, отправляют его в общий канал и корректно завершают работу.
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var totalCoal atomic.Int32

	workdayCtx, cancelWorkday := context.WithCancel(context.Background())
	defer cancelWorkday()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Рабочий день завершён")
		cancelWorkday()
	}()

	for minedCoal := range startMinerPool(workdayCtx, 3) {
		totalCoal.Add(int32(minedCoal))
	}

	fmt.Println("Всего добыто угля:", totalCoal.Load())
}

func startMinerPool(ctx context.Context, minerCount int) <-chan int {
	coalTransferCh := make(chan int)

	var minerGroup sync.WaitGroup
	minerGroup.Add(minerCount)

	for minerID := 1; minerID <= minerCount; minerID++ {
		go mineCoal(ctx, minerID, minerID*10, coalTransferCh, &minerGroup)
	}

	go func() {
		minerGroup.Wait()
		close(coalTransferCh)
	}()

	return coalTransferCh
}

func mineCoal(ctx context.Context, minerID, miningPower int, output chan<- int, minerGroup *sync.WaitGroup) {
	defer minerGroup.Done()

	for {
		fmt.Println("Шахтёр", minerID, "начал добычу")

		select {
		case <-ctx.Done():
			fmt.Println("Шахтёр", minerID, "завершил работу")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("Шахтёр", minerID, "добыл уголь:", miningPower)
		}

		select {
		case <-ctx.Done():
			fmt.Println("Шахтёр", minerID, "завершил работу")
			return
		case output <- miningPower:
			fmt.Println("Шахтёр", minerID, "передал уголь:", miningPower)
		}
	}
}
