package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	// "os"
)

type Course struct {
	CourseName  string
	CourseID    string
	University  string
	Instructor  string
	Quarter     string
	Link        string
	ContentType string
	Visited     int
	Liked       int
}

type ListCourses struct {
	Title   string
	Courses []Course
}

func initCourseDB() {
	// Create database if doesn't already exist
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create courses table if doesn't exist
	sqlStmt := `CREATE TABLE IF NOT EXISTS courses (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      courseName TEXT,
      courseID TEXT,
      university TEXT,
      instructor TEXT,
      quarter TEXT,
      link TEXT,
	  contentType TEXT,
	  visited INTEGER,
	  like INTEGER
    )`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func insertCourse(course Course) (int, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO courses (courseName, courseID, university, instructor, quarter, link, contentType)
		VALUES(?,?,?,?,?,?,?);`,
		course.CourseName, course.CourseID, course.University, course.Instructor,
		course.Quarter, course.Link, course.ContentType)
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
