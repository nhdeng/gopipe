package v3

import (
	"fmt"
	"sync"
	"time"
)

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

type Cmd func([]int) chan int
type PipeCmd func(ch chan int) (ret chan int)

// Pipe 管道函数
func Pipe(nums []int, f1 Cmd, ps ...PipeCmd) chan int {
	wg := sync.WaitGroup{}
	evench := f1(nums) // 找偶数
	out := make(chan int)
	for _, p := range ps {
		getChan := p(evench)
		wg.Add(1)
		go func(ch chan int) {
			defer wg.Done()
			for c := range ch {
				out <- c
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}

func Test(nums []int) {
	res := Pipe(nums, Events, M2, M2)
	for v := range res {
		fmt.Printf("%d ", v)
	}
}
