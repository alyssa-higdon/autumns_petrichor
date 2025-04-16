package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Input  : w http.ResponceWriter
//
//	tmpl string     : name of template to render (ex: Index.html)
//	data any        : data that needs to be passed to fill the template (Ex: Header names, course data to fill pages)
//
// Returns: None
// Output : Renders the desired template with the header and footer
// Purpose: Render out the desired page with the header and footer. Layout.html describes how to combine them
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

// Input  : w http.ResponseWriter
//
//	r *http.Request
//
// Returns: None
// Output : The home page
// Purpose: This function is a wrapper to renderTemplate, where all of the data for
//
//	the home page is processed before sending it to renderTemplate
func homePage(w http.ResponseWriter, r *http.Request) {
	courses := []Course{}
	universities := []string{}

	getAllCourses(&courses)
	getDistProperty("university", &universities)

	// Check for a university filter in the query parameters
	selectedUniversity := r.URL.Query().Get("university")
	retList := []Course{}
	findCourses(map[string]string{"university": selectedUniversity}, &retList)

	if selectedUniversity != "" {
		courses = retList
	}
	data := PageContent{
		Title: "Home",
		Data: Home{
			LCourses: ListCourses{
				Title:   "All Courses",
				Courses: courses,
			},
			SBar: SideBar{
				Sections: []Section{
					{Title: "university",
						Items: universities},
				},
			},
		},
	}

	if searchCourse := r.FormValue("searchCourse"); searchCourse == "" {
		renderTemplate(w, "Index.html", data)
	}

}

func homePagePOST(w http.ResponseWriter, r *http.Request) {
	if searchCourse := r.FormValue("searchCourse"); searchCourse != "" {
		var searchCourseSp = strings.Split(searchCourse, " - ")
		var searchCourseName = searchCourseSp[0]
		var searchCourseUni = searchCourseSp[1]

		var reqs = map[string]string{
			"courseName": searchCourseName,
			"university": searchCourseUni}

		retList := []Course{}
		findCourses(reqs, &retList)
		if len(retList) > 0 {
			http.Redirect(w, r, "/course/"+strconv.Itoa(retList[0].Id), http.StatusSeeOther)
		} else {
			http.Error(w, "Course not found", http.StatusNotFound)
		}

	}
}

// Input  : w http.ResponseWriter
//
//	r *http.Request
//
// Returns: None
// Output : The home page
// Purpose: This function is a wrapper to renderTemplate, where all of the data for
//
//	a course page is processed before sending it to renderTemplate
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
	//coursesStarterData()
	// allUniversities := []string{}
	// getDistProperty("university", &allUniversities)

	// Serve static assets (ie. css)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/submit", homePagePOST)
	http.HandleFunc("/course/", coursePage)

	// Start the server
	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
