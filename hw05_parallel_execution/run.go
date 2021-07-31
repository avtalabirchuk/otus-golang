package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorsCounter struct {
	mu         sync.Mutex
	counter    int
	errorLimit int
}

func (ec *errorsCounter) reachedLimit() bool {
	defer ec.mu.Unlock()
	ec.mu.Lock()
	return ec.counter > ec.errorLimit
}

func (ec *errorsCounter) increase() {
	defer ec.mu.Unlock()
	ec.mu.Lock()
	ec.counter++
}

type Worker struct {
	wg            *sync.WaitGroup
	tasksChannel  chan Task
	errorsCounter *errorsCounter
}

// обработка канала пока он не будет закрыт или не будет достигнут предел ошибок

func (w Worker) Working() {
	for {
		if w.errorsCounter.reachedLimit() {
			break
		}
		task, channelIsOpen := <-w.tasksChannel
		if !channelIsOpen {
			break
		}
		taskError := task()
		if taskError != nil {
			w.errorsCounter.increase()
		}
	}
	w.wg.Done()
}

// помещает все задачи по одной в наш канал и затем закрывает его
func putToChannel(tasks []Task, channel chan<- Task) {
	for _, task := range tasks {
		channel <- task
	}
	close(channel)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorsCounter := &errorsCounter{
		counter:    0,
		errorLimit: m,
	}
	tasksChannel := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	wg.Add(n)
	// запуск заданного кол-ва воркеров n
	for i := 0; i < n; i++ {
		worker := Worker{
			wg:            &wg,
			tasksChannel:  tasksChannel,
			errorsCounter: errorsCounter,
		}
		go worker.Working()
	}
	// кладем все задачи в канал и ждем когда они обработаются
	go putToChannel(tasks, tasksChannel)
	wg.Wait()
	if errorsCounter.reachedLimit() {
		return ErrErrorsLimitExceeded
	}
	return nil
}
