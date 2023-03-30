package v2

import (
	"encoding/csv"
	"fmt"
	"gochan/src/mysql2csv"
	"os"
	"strconv"
	"sync"
	"time"
)

// TODO 使用管道及多路复用的方式优化数据库数据导出
const sql = "select * from books order by book_id limit ? offset ?"

type BookList struct {
	BookId   int    `gorm:"column:book_id"`
	BookName string `gorm:"column:book_name"`
}

type QueryResult struct {
	Data []*BookList
	Page int
}

type Result struct {
	Page int
	Err  error
}

type InChan chan *QueryResult
type OutChan chan *Result

// GetData 从数据库获取数据
func GetData() InChan {
	page := 1
	pageSize := 1000
	in := make(InChan, 0)
	go func() {
		defer close(in)
		for {
			result := &QueryResult{[]*BookList{}, page}
			db := mysql2csv.InitDB().Raw(sql, pageSize, (page-1)*pageSize).Find(&result.Data)
			if db.Error != nil || db.RowsAffected == 0 {
				break
			}
			in <- result
			page++
		}
	}()

	return in
}

// Data2Csv 将数据保存到csv文件中
func Data2Csv(res *QueryResult) error {
	time.Sleep(time.Millisecond * 500)
	path := fmt.Sprintf("./src/mysql2csv/csv/%d.csv", res.Page)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	header := []string{"book_id", "book_name"}
	export := [][]string{header}
	for _, data := range res.Data {
		content := []string{strconv.Itoa(data.BookId), data.BookName}
		export = append(export, content)
	}
	err = w.WriteAll(export)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func WriteData(inChan InChan) OutChan {
	out := make(chan *Result, 0)
	go func() {
		defer close(out)
		for i := range inChan {
			out <- &Result{i.Page, Data2Csv(i)}
		}
	}()
	return out
}

type CmdFunc func() InChan
type PipeCmdFunc func(inChan InChan) OutChan

// Pipe 管道函数
func Pipe(cmd CmdFunc, pipeCmdFuncs ...PipeCmdFunc) OutChan {
	wg := sync.WaitGroup{}
	out := make(OutChan, 0)
	in := cmd()
	for _, pcf := range pipeCmdFuncs {
		res := pcf(in)
		wg.Add(1)
		go func(input OutChan) {
			defer wg.Done()
			for i := range input {
				out <- i
			}
		}(res)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}

func Test() {
	out := Pipe(GetData, WriteData, WriteData)
	for d := range out {
		log := fmt.Sprintf("%d.csv文件导出完毕，错误：%v", d.Page, d.Err)
		fmt.Println(log)
	}
}
