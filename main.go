package main

import (
	"fmt"
	"os"
	"time"

	hook "github.com/robotn/gohook"
)

var isRunning bool // Глобальная переменная для управления состоянием кликера

func main() {
	count := &AtomicCounter{}
	keyEvents := make(chan hook.Event, 10)
	startSignal := make(chan struct{}, 1)
	keyEventStart := "+"
	keyEventStop := "-"

	// Запускаем регистрацию событий
	go func() {
		hook.Register(hook.KeyDown, []string{keyEventStop}, func(e hook.Event) {
			keyEvents <- e
		})
		hook.Register(hook.KeyDown, []string{keyEventStart}, func(e hook.Event) {
			keyEvents <- e
		})
		s := hook.Start()
		<-hook.Process(s)
	}()

	// Ожидаем сигнал старта кликера
	waitForStart(keyEvents, startSignal)

	// Запускаем кликер с таймером 3 секунды
	runClickerWithTimeout(count, keyEvents, 3*time.Second)

	fmt.Printf("Количество кликов: %d\n", count.Get())
	os.Exit(0)
}
