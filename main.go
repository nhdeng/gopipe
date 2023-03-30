package main

import (
	"fmt"
	"gopipe/src/even/v1"
	"gopipe/src/even/v2"
	"gopipe/src/even/v3"
	mysql2csv1 "gopipe/src/mysql2csv/v1"
	mysql2csv2 "gopipe/src/mysql2csv/v2"
	"time"
)

func test(v string) {
	nums := []int{1, 2, 3, 4, 6, 8, 9, 11, 32}
	start := time.Now().Unix()
	if v == "v1" {
		v1.Test(nums)
	}
	if v == "v2" {
		v2.Test(nums)
	}
	if v == "v3" {
		v3.Test(nums)
	}
	end := time.Now().Unix()
	fmt.Printf("总计耗时：%d秒 \n", end-start)
}

func testReadData() {
	start := time.Now().Unix()
	mysql2csv1.ReadData()
	end := time.Now().Unix()
	fmt.Printf("总计耗时：%d秒 \n", end-start)
}

func testPipeReadData() {
	start := time.Now().Unix()
	mysql2csv2.Test()
	end := time.Now().Unix()
	fmt.Printf("总计耗时：%d秒 \n", end-start)
}
func main() {
	//test("v1")
	//test("v2")
	//test("v3")
	//testReadData()
	testPipeReadData()
}
