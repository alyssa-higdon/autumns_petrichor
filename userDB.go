package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	// "os"
)

type User struct {
	fName string
	lName string
	email string
	pwd string
}

// Input  : None
// Returns: None
// Output : If the users table in the database.db doesn't already exist, create the table with
//          columns for a firstname, lastname, email, and password
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
      email TEXT,
      pwd   TEXT,
    )`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// TODO: Check if the email is already in the db
// TODO: Password encyption
func insertUser(user User) int {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	// Add password after figure out encryption
	res, err := db.Exec(`INSERT INTO users (fName, lName, email, pwd) VALUES(?,?,?,?);`,
		user.fName, user.lName, user.email)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		log.Fatal(err)
	}
	return int(id)
}
