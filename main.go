package main

import (
	"html/template"
	"net/http"
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

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("./components/Index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := ListCourses{
			Title: "All Courses",
			Courses: []Course{
				{CourseName: "Life Science for Engineers", CourseID: "BIO 213", University: "Cal Poly", Instructor: "Babu",
					Quarter: "", Link: "https://drive.google.com/drive/folders/1BEZBngPP2oWIeY5ZtsdQw-Wfj60zTLt8?usp=sharing",
					ContentType: "Notes"},
				{CourseName: "Mechanics of Materials", CourseID: "CE 208", University: "Cal Poly", Instructor: "Elghandour",
					Quarter: "", Link: "https://drive.google.com/drive/folders/1XzEniwg4VlQzx4QnPLHwj8MEWO7npUcn?usp=sharing",
					ContentType: "Notes"},
				{CourseName: "General Chemistry for Physical Science and Engineering I", CourseID: "CHEM 124", University: "Cal Poly", Instructor: "Campbell",
					Quarter: "Winter 2020", Link: "https://drive.google.com/drive/folders/1IrwmB22AvKrlfnCX7_uptpVNYQnq_t80?usp=sharing",
					ContentType: "Notes"},
				{CourseName: "Data Structures", CourseID: "CPE 202", University: "Cal Poly", Instructor: "Parkinson",
					Quarter: "", Link: "https://drive.google.com/drive/u/2/folders/16Oko_STFETBLTgFUexa9iPFnrh4Uew2N",
					ContentType: "Notes"},
				{CourseName: "Computer Architecture", CourseID: "CSC 315", University: "Cal Poly", Instructor: "Seng",
					Quarter: "Spring 2022", Link: "https://photos.app.goo.gl/fmDA2vogjGGSeR9YA",
					ContentType: "Notes"},
				{CourseName: "Linear Analysis I", CourseID: "MATH 244", University: "Cal Poly", Instructor: "Choboter",
					Quarter: "", Link: "https://drive.google.com/drive/folders/1wCnOaarXPRB7Jxfa85I34iSQ7dpvMt20?usp=sharing",
					ContentType: "Notes"},
				{CourseName: "Ecology I: The Earth System", CourseID: "1.018J", University: "MIT", Instructor: "DeLong, Chisholm",
					Quarter: "", Link: "https://ocw.mit.edu/courses/1-018j-ecology-i-the-earth-system-fall-2009/",
					ContentType: "Entire Course"},
			},
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":80", nil)
}
