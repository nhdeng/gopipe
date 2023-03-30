package v1

import "fmt"

// Events 找出偶数
func Events(input []int) []int {
	out := make([]int, 0)
	for i := 0; i < len(input); i++ {
		if input[i]%2 == 0 {
			out = append(out, input[i])
		}
	}
	return out
}

// M2 数字乘2
func M2(input []int) []int {
	out := make([]int, 0)
	for i := 0; i < len(input); i++ {
		out = append(out, input[i]*2)
	}
	return out
}

// M5 数字乘5
func M5(input []int) []int {
	out := make([]int, 0)
	for i := 0; i < len(input); i++ {
		out = append(out, input[i]*5)
	}
	return out
}

// Cmd 定义管道函数的参数类型
type Cmd func(list []int) (ret []int)

// Pipe 管道函数
func Pipe(nums []int, f1 Cmd, f2 Cmd, f3 Cmd) []int {
	return f3(f2(f1(nums)))
}

// Test 示例调用管道函数
func Test(nums []int) {
	res := Pipe(nums, Events, M2, M5)
	for _, val := range res {
		fmt.Printf("%d ", val)
	}
}
