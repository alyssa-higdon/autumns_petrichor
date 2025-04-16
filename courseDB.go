package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
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

type Section struct {
	Title string
	Items []string
}

type SideBar struct {
	Sections []Section
}

type Home struct {
	LCourses ListCourses
	SBar     SideBar
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

// Input  : course Course : course to insert into the course table
// Returns: int : id of the course (index of the course in the table)
// Output : None
// Purpose: Insert course into courses table if it's not already in the table
func insertCourse(course Course) int {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO courses (courseName, courseID, university, instructor, quarter, link, contentType, visited, liked)
	SELECT ?, ?, ?, ?, ?, ?, ?, 0, 0
	WHERE NOT EXISTS (
		SELECT 1 
		FROM courses 
		WHERE courseName = ? AND university = ?
	);`,
		course.CourseName, course.CourseID, course.University, course.Instructor,
		course.Quarter, course.Link, course.ContentType, course.CourseName, course.University)

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
//
//	retList *[]string : A list to return all of the items of the specific property
//
// Returns: None
// Output : List of all of the items of the specific property
// Purpose: To find all items of a specific property in the courses table
// Ex: Get all university names
// This is good for the side navbar to choose only classes from certain universities
func getDistProperty(property string, retList *[]string) {
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
		*retList = append(*retList, property)
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

// Input  : reqs   : map of the requirements you're looking for (ie {{courseName: Life Science for Engineers}, {university: Cal Poly}})

// 	retList: *[]Course : a pointer to a list to return of all Courses with the desired reqs

// Returns:
// Output : Course : Course that you are querying for
// Purpose: Finds the Course based on the CourseName and University
func findCourse(reqs map[string]string, retList *[]Course) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Build query statement
	var query = `SELECT * from courses WHERE `
	var i = len(reqs)
	for param, arg := range reqs {
		query += param + "='" + arg + "'"
		i--
		if i != 0 {
			query += " AND "
		}
	}

	fmt.Println("", query)
	rows, err := db.Query(query)
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
		fmt.Println("", course.CourseName)
	}
}
