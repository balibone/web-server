package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
)

// templates is used to store multiple templates that were parsed, producing 1 Template object.
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

// validPath is used to validate page titles submitted by users.
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Page represents a wiki page
type Page struct {
	Title string
	Body  []byte
}

//Save is a method on Page struct that will save the page into a text file.
func (page *Page) Save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filepath.Join("data", filename), page.Body, 0600)
}

//LoadPage loads a text file, reads it and creates a new Page literal from its content.
func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filepath.Join("data", filename))
	if err != nil {
		return nil, err
	}
	return &Page{
		Title: title,
		Body:  body,
	}, nil
}

func main() {
	//this is the routing of endpoint to handler. Like url mappings in web.xml.
	http.HandleFunc("/", frontPageHandler)
	http.HandleFunc("/view/", makeHandler(ViewHandler))
	http.HandleFunc("/edit/", makeHandler(EditHandler))
	http.HandleFunc("/save/", makeHandler(SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ViewHandler is a http.HandlerFunc serves a page requested from the "/view/*" endpoint.
// If page doesn't exist, then this handler sends a redirect to the edit page so that the page can be created.
func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", page)
}

// EditHandler is a http.HandlerFunc that serves an edit page.
func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := LoadPage(title)
	if err != nil {
		//if there's error loading page, then just return a page with the title = requested title.
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// SaveHandler is a http.HandlerFunc that saves the edit performed on /edit/ page.
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//get the form value which belongs to the field with name (or key) attribute "body"
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//after save, redirect to view page
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// renderTemplate executes a template that is specified by "tmpl" file name, and is already parsed into the "templates" variable.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//execute templates by passing in data (in this case page) that is required to fill up its values.
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// makeHandler returns a closure http.HandlerFunc that pre-checks the title string in the url.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//CASE CONTENTS
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Printf("Result of validPath.FindStringSubmatch is : %q\n", m)
		//use closure function to call the actual handler..
		fn(w, r, m[2])
	}
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	p, err := LoadPage("FrontPage")
	if err != nil {
		//if there's error loading page, then just return a page with the title = requested title.
		p = &Page{Title: "FrontPage"}
	}
	renderTemplate(w, "view", p)
}
