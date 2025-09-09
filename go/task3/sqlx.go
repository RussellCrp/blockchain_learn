package main

import (
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
