package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func coursesStarterData() {
	insertCourse(Course{CourseName: "Life Science for Engineers", CourseID: "BIO 213", University: "Cal Poly", Instructor: "Babu",
		Quarter: "", Link: "https://drive.google.com/drive/folders/1BEZBngPP2oWIeY5ZtsdQw-Wfj60zTLt8?usp=sharing",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Mechanics of Materials", CourseID: "CE 208", University: "Cal Poly", Instructor: "Elghandour",
		Quarter: "", Link: "https://drive.google.com/drive/folders/1XzEniwg4VlQzx4QnPLHwj8MEWO7npUcn?usp=sharing",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "General Chemistry for Physical Science and Engineering I", CourseID: "CHEM 124", University: "Cal Poly", Instructor: "Campbell",
		Quarter: "Winter 2020", Link: "https://drive.google.com/drive/folders/1IrwmB22AvKrlfnCX7_uptpVNYQnq_t80?usp=sharing",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Data Structures", CourseID: "CPE 202", University: "Cal Poly", Instructor: "Parkinson",
		Quarter: "", Link: "https://drive.google.com/drive/u/2/folders/16Oko_STFETBLTgFUexa9iPFnrh4Uew2N",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Data Structures", CourseID: "CPE 202", University: "Cal Poly", Instructor: "Parkinson",
		Quarter: "", Link: "https://drive.google.com/drive/u/2/folders/16Oko_STFETBLTgFUexa9iPFnrh4Uew2N",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Computer Architecture", CourseID: "CSC 315", University: "Cal Poly", Instructor: "Seng",
		Quarter: "Spring 2022", Link: "https://photos.app.goo.gl/fmDA2vogjGGSeR9YA",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Linear Analysis I", CourseID: "MATH 244", University: "Cal Poly", Instructor: "Choboter",
		Quarter: "", Link: "https://drive.google.com/drive/folders/1wCnOaarXPRB7Jxfa85I34iSQ7dpvMt20?usp=sharing",
		ContentType: "Notes"})
	insertCourse(Course{CourseName: "Ecology I: The Earth System", CourseID: "1.018J", University: "MIT", Instructor: "DeLong, Chisholm",
		Quarter: "", Link: "https://ocw.mit.edu/courses/1-018j-ecology-i-the-earth-system-fall-2009/",
		ContentType: "Entire Course"})

}

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// Define paths to all your templates
	layoutPath := "./components/Layout.html"
	headerPath := "./components/Header.html"
	footerPath := "./components/Footer.html"
	pagePath := "./components/" + tmpl

	// Parse the templates (you can add more templates to the list as needed)
	t, err := template.ParseFiles(layoutPath, headerPath, footerPath, pagePath)
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}

	// Render the template with the provided data
	err = t.ExecuteTemplate(w, "Layout", data) // "Layout" is the name of the main template
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	courses := []Course{}
	getAllCourses(&courses)
	data := PageContent{
		Title: "Home",
		Data: ListCourses{
			Title:   "All Courses",
			Courses: courses,
		},
	}
	renderTemplate(w, "Index.html", data)
}

func coursePage(w http.ResponseWriter, r *http.Request) {
	// Get the course ID from the URL (e.g., /course/1)
	id := strings.TrimPrefix(r.URL.Path, "/course/")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal("Failed casting course/id: ", err)
	}

	courses := []Course{}
	getAllCourses(&courses)
	// Find the course by ID
	var selectedCourse *Course
	for _, c := range courses {
		if c.Id == idInt {
			selectedCourse = &c
			break
		}
	}

	if selectedCourse == nil {
		http.NotFound(w, r)
		return
	}

	// Render the course detail page
	data := PageContent{
		Title: selectedCourse.CourseName,
		Data:  selectedCourse}
	renderTemplate(w, "Course.html", data)
}

func main() {
	initCourseDB()
	initUserDB()
	// coursesStarterData()
	allUniversities := []string{}
	getDistProperty("university", allUniversities)

	// database()
	// Serve static assets (ie. css)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/course/", coursePage)
	// http.HandleFunc("/about", aboutPage)

	// Start the server
	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
