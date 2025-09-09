package main

import (
	"fmt"
	"time"
)

func main() {
	channelCommunication()
	bufferChannel()
}

func channelCommunication() {
	nums := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			nums <- i
			fmt.Println("通道中放入参数", i)
			time.Sleep(time.Millisecond * 300)
		}
		close(nums)
	}()

	go func() {
		for {
			num, ok := <-nums
			if ok {
				fmt.Println("通道中接收到参数", num)
			} else {
				break
			}

		}
	}()

	select {
	case <-time.After(time.Second * 5):
		return
	}
}

func bufferChannel() {
	ints := make(chan int, 10)
	go func() {
		for i := 0; i < 100; i++ {
			ints <- i
			fmt.Println("生产者生产: ", i)
		}
		close(ints)
	}()

	go func() {
		for num := range ints {
			fmt.Println("消费者消费: ", num)
			time.Sleep(time.Millisecond * 30)
		}
	}()

	select {
	case <-time.After(time.Second * 6):
		return
	}
}
