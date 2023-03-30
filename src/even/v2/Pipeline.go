package v2

import (
	"fmt"
	"sync"
	"time"
)

// 以channel的方式进行优化，不会造成堵塞

// Events 从切片中找出偶数
func Events(input []int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < len(input); i++ {
			if input[i]%2 == 0 {
				ch <- input[i]
			}
		}
	}()
	return ch
}

// M2 将偶数乘以2
func M2(input chan int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range input {
			time.Sleep(time.Second * 2)
			ch <- i * 2
		}
	}()
	return ch
}

// M5 将偶数乘以5
func M5(input chan int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range input {
			time.Sleep(time.Second * 2)
			ch <- i * 5
		}
	}()
	return ch
}

type Cmd func([]int) chan int

// PipeCmd 定义管道函数的参数类型
type PipeCmd func(ch chan int) (ret chan int)

// Pipe 管道函数
func Pipe(nums []int, f1 Cmd, f2 PipeCmd, f3 PipeCmd) chan int {
	return f3(f2(f1(nums)))
}

// Test 示例调用管道函数
func Test(nums []int) {
	wg := sync.WaitGroup{}
	res := Pipe(nums, Events, M2, M5)
	for v := range res {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%d ", v)
		}()
	}
	wg.Wait()
}
