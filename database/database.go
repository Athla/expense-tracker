package database

import (
	"log"

	"github.com/Athla/expense-tracker.git/expense"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Conn() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "../db/expenses.db")
	//TODO Add context for better error handling
	if err != nil {
		log.Fatalln(err)
	}
	return db, nil
}

func (db *Database) Add(e *expense.Expense) {
	conn, err := Conn()
	//TODO Add context for better error handling
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	tx, err := conn.Begin()
	//TODO Add context for better error handling
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := tx.Exec(
		"INSERT INTO DATA(TYPE, NAME, VALUE, DESCRIPTION) VALUES(?, ?, ?, ?)",
		e.Type,
		e.Name,
		e.Value,
		e.Description,
	); err != nil {
		//TODO Add context for better error handling
		log.Fatalln(err)
	}
	log.Println("Data sucess")
}
