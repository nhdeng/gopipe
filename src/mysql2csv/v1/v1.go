package v1

import (
	"encoding/csv"
	"fmt"
	"gopipe/src/mysql2csv"
	"os"
	"strconv"
)

const sql = "select * from books order by book_id limit ? offset ?"

type Book struct {
	BookId   int    `gorm:"column:book_id"`
	BookName string `gorm:"column:book_name"`
}

type BookList struct {
	Data []*Book
	Page int
}

// ReadData 读取数据
func ReadData() {
	page := 1
	pageSize := 1000

	for {
		bookList := &BookList{Data: make([]*Book, 0), Page: page}
		db := mysql2csv.InitDB().Raw(sql, pageSize, (page-1)*pageSize).Find(&bookList.Data)
		if db.Error != nil || db.RowsAffected == 0 {
			break
		}
		err := SaveData(bookList)
		if err != nil {
			fmt.Println(err)
		}
		page++
	}
}

// SaveData 保存数据到csv文件中
func SaveData(list *BookList) error {
	filePath := fmt.Sprintf("./src/mysql2csv/csv/%d.csv", list.Page)
	csvFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	header := []string{"book_id", "book_name"}
	export := [][]string{header}

	for _, data := range list.Data {
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
