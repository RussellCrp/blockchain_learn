package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	goPrintNum(10)

	schedul()
}

func goPrintNum(num int) {
	var done sync.WaitGroup
	done.Add(2)
	go printEvenNum(num, &done)
	go printOddNum(num, &done)
	done.Wait()
}

func printEvenNum(num int, done *sync.WaitGroup) {
	defer done.Done()
	for i := 1; i <= num; i++ {
		if i%2 == 0 {
			fmt.Println("打印偶数：", i)
			time.Sleep(time.Millisecond * 2)
		}
	}
}

func printOddNum(num int, done *sync.WaitGroup) {
	defer done.Done()
	for i := 1; i <= num; i++ {
		if i%2 == 1 {
			fmt.Println("打印奇数：", i)
			time.Sleep(time.Millisecond * 2)
		}
	}
}

// ==================================================================

func schedul() {
	num := 1
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	task1 := &PriceAddTask{&num}
	task2 := &SliceMul2Task{&slice}
	tasks := []Task{task1, task2}
	scheduler := &scheduler{tasks}
	scheduler.run()

	time.Sleep(time.Second * 2)
}

type Task interface {
	Exec()
	TaskName() string
}

type PriceAddTask struct {
	num *int
}

func (t *PriceAddTask) Exec() {
	ptrAdd(t.num)
	time.Sleep(time.Millisecond * 100)
}

func (t *PriceAddTask) TaskName() string {
	return "指针加十"
}

type SliceMul2Task struct {
	ptrSlice *[]int
}

func (t *SliceMul2Task) Exec() {
	ptrSlice(t.ptrSlice)
	time.Sleep(time.Millisecond * 300)
}

func (t *SliceMul2Task) TaskName() string {
	return "切片*2"
}

type scheduler struct {
	tasks []Task
}

func (s *scheduler) run() {
	for _, task := range s.tasks {
		go func() {
			taskName := task.TaskName()
			fmt.Println(taskName, " 开始执行。。。")
			start := time.Now()
			task.Exec()
			end := time.Now()
			fmt.Printf("任务 %s 结束，耗时 %v\n", taskName, end.Sub(start))
		}()
	}
}
