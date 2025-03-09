package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	// "os"
)

type PageContent struct {
	Title string
	Data  any
}

type Course struct {
	Id          int
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

// Input  : None
// Returns: None
// Output : Creates db if doesn't exist and creates course table if doesn't exist
// Purpose: Creates db if doesn't exist and creates course table if doesn't exist
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
	  liked INTEGER
    )`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// TODO: Check if the course is already in the db
// Input  : course Course : course to insert into the course table
// Returns: int : id of the course (index of the course in the table)
// Output : None
// Purpose: Insert course into courses table
func insertCourse(course Course) int {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO courses (courseName, courseID, university, instructor, quarter, link, contentType, visited, liked)
		VALUES(?,?,?,?,?,?,?,0,0);`,
		course.CourseName, course.CourseID, course.University, course.Instructor,
		course.Quarter, course.Link, course.ContentType)
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

// Input  : property string  : Property of course table searching for
//        retList []string : A list to return all of the items of the specific property
// Returns: None
// Output : List of all of the items of the specific property
// Purpose: To find all items of a specific property in the courses table
// Ex: Get all university names
// This is good for the side navbar to choose only classes from certain universities
func getDistProperty(property string, retList []string) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT DISTINCT " + property + " from courses;")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var property string
		if err := rows.Scan(&property); err != nil {
			log.Fatal(err)
		}
		retList = append(retList, property)
	}
}

// Input  : retList *[]Course : a pointer to a list to return of all Courses in the courses table
// Returns: None
// Output : a pointer to a list to return of all Courses in the courses table
// Purpose: get all Courses in the course table
func getAllCourses(retList *[]Course) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT * from courses`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer rows.Close()
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.Id, &course.CourseName, &course.CourseID, &course.University, &course.Instructor,
			&course.Quarter, &course.Link, &course.ContentType, &course.Visited, &course.Liked); err != nil {
			log.Fatal(err)
		}
		*retList = append(*retList, course)
	}
}

// Input  : courseName string : name of the course searching for
// Returns: Course : Course that you are querying for
// Output : Course : Course that you are querying for
// Purpose: Finds the Course based on the CourseName and University
func findCourse(courseName string, university string) Course {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT * from courses WHERE courseName=? AND university=?`, courseName, university)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer rows.Close()
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.Id, &course.CourseName, &course.CourseID, &course.University, &course.Instructor,
			&course.Quarter, &course.Link, &course.ContentType, &course.Visited, &course.Liked); err != nil {
			log.Fatal(err)
		}
		return course
	}
	return Course{}
}
