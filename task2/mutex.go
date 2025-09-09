package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type counter struct {
	mu  sync.Mutex
	num int64
}

func (c *counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.num++
}

func (c *counter) atomicIncrement() {
	atomic.AddInt64(&c.num, 1)
}

func mutex() {
	var done sync.WaitGroup
	counter := &counter{}
	for i := 0; i < 10; i++ {
		go func() {
			done.Add(1)
			for i := 0; i < 1000; i++ {
				counter.increment()
			}
			done.Done()
		}()
	}
	done.Wait()
	fmt.Println("互斥锁方式,共享的计数器值: ", counter.num)
}

func atomicAdd() {
	var done sync.WaitGroup
	counter := &counter{}
	for i := 0; i < 10; i++ {
		go func() {
			done.Add(1)
			for i := 0; i < 1000; i++ {
				counter.atomicIncrement()
			}
			done.Done()
		}()
	}
	done.Wait()
	fmt.Println("原子方式,共享的计数器值: ", counter.num)
}

func main() {
	mutex()
	atomicAdd()
}
