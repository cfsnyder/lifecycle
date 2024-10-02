package simplelife

import "sync"

type TaskFunc func() Stopper

type StopFunc func()

type Task interface {
	Start() Stopper
}

type Stopper interface {
	Stop()
}

func (self TaskFunc) Start() Stopper {
	return self()
}

func (self StopFunc) Stop() {
	self()
}

func Seq(tasks ...Task) Task {
	var task TaskFunc = func() Stopper {
		var stoppers []Stopper
		for _, task := range tasks {
			stoppers = append(stoppers, task.Start())
		}
		var stopper StopFunc = func() {
			for i := len(stoppers) - 1; i >= 0; i-- {
				stoppers[i].Stop()
			}
		}
		return stopper
	}
	return task
}

func Parallel(tasks ...Task) Task {
	var task TaskFunc = func() Stopper {
		ch := make(chan Stopper, len(tasks))
		for _, task := range tasks {
			go func(t Task) {
				ch <- t.Start()
			}(task)
		}
		var stoppers []Stopper
		for i := 0; i < len(tasks); i++ {
			stoppers = append(stoppers, <-ch)
		}
		var stopper StopFunc = func() {
			var wg sync.WaitGroup
			for _, stopper := range stoppers {
				wg.Add(1)
				go func(s Stopper) {
					s.Stop()
					wg.Done()
				}(stopper)
			}
			wg.Wait()
		}
		return stopper
	}
	return task
}
