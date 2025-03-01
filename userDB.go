package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	// "os"
)

// add password after I figure out encryption
type User struct {
	fName string
	lName string
	email string
}

func initUserDB() {
	// Create database if doesn't already exist
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Add password once figure out encryption
	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      fName TEXT,
	  lName TEXT,
      email TEXT
    )`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func insertUser(user User) (int, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	// Add password after figure out encryption
	res, err := db.Exec(`INSERT INTO users (fName, lName, email) VALUES(?,?,?);`,
		user.fName, user.lName, user.email)
	defer db.Close()
	if err != nil {
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}
