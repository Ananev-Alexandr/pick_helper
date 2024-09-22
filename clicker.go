package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	rg "github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const (
	keyEventStart = '+'
	keyEventStop  = '='
)

// Функция для выполнения кликов в горутинах
func runClicker(ctx context.Context, count *AtomicCounter, taskChan chan struct{}) {
	for {
		select {
		case <-ctx.Done(): // Завершение по таймеру
			return
		case _, ok := <-taskChan: // Выполняем задачу клика
			if ok {
				rg.Click("left", true) // Выполняем клик
				count.Inc()            // Увеличиваем счетчик кликов
			} else {
				return // Завершаем горутину, если канал закрыт
			}
		}
	}
}

// Запуск большого числа горутин для кликов с таймером на 3 секунды
func runClickerWithTimeout(count *AtomicCounter, keyEvents chan hook.Event, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	numGoroutines := 100
	taskChan := make(chan struct{}, numGoroutines) // Канал для передачи задач в воркеры

	// Запуск множества горутин для выполнения кликов
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runClicker(ctx, count, taskChan)
		}()
	}

	// Горутина для ручного завершения по клавише "="
	go func() {
		for {
			select {
			case e := <-keyEvents:
				if e.Keychar == keyEventStop && isRunning {
					fmt.Println("Остановлен вручную")
					cancel()          // Завершаем контекст
					isRunning = false // Останавливаем кликер
					return
				}
			case <-ctx.Done():
				isRunning = false // Останавливаем по истечению таймера
				return
			}
		}
	}()

	// Посылка задач воркерам для выполнения кликов
	for {
		select {
		case <-ctx.Done():
			close(taskChan) // Закрываем канал задач, когда все завершено
			wg.Wait()       // Ждем завершения всех горутин
			fmt.Println("Все клики завершены")
			return
		default:
			if isRunning {
				taskChan <- struct{}{} // Посылаем задачу для клика
			}
		}
	}
}

// Ожидание старта кликера
func waitForStart(keyEvents chan hook.Event, startSignal chan struct{}) {
	for {
		select {
		case e := <-keyEvents:
			if e.Keychar == keyEventStart && !isRunning {
				fmt.Println("Получен сигнал старта!")
				isRunning = true   // Устанавливаем флаг, что кликер запущен
				close(startSignal) // Сигнализируем о старте, после чего канал можно закрыть
				return
			}
		}
	}
}
