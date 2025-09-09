package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Employee struct {
	Id         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float32 `db:"salary"`
}

func queryEmployees(db *sqlx.DB) ([]Employee, error) {
	var e []Employee
	err := db.Select(&e, "SELECT * FROM Employee where department = ?", "技术部")
	if err != nil {
		return nil, err
	}
	return e, nil
}

func queryMaxSalary(db *sqlx.DB) (*Employee, error) {
	e := &Employee{}
	err := db.Get(e, "SELECT * FROM Employee order by salary desc limit 1")
	if err != nil {
		return nil, err
	}
	return e, nil
}

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"` // 使用 float64 更适合货币类型
}

func queryExpensiveBooks(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book
	err := db.Select(&books, "SELECT id, title, author, price FROM books WHERE price > ?", minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}
	return books, nil
}
