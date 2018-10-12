package datastructures

import (
	"fmt"
	"io/ioutil"
)

//Page represents a wiki page
type Page struct {
	Title string
	Body  []byte
}

//Save is a method on Page struct that will save the page into a text file.
func (page *Page) Save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}

//LoadPage loads a text file, reads it and creates a new Page literal from its content.
func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title: title,
		Body:  body,
	}, nil
}

func main() {
	page1 := &Page{
		Title: "First Page",
		Body:  []byte("This is my first Page."),
	}
	page1.Save()
	page2, err := LoadPage(page1.Title)
	if err != nil {
		fmt.Println("There was an error in LoadPage func when reading page1.")
	}
	fmt.Println("Title:", page2.Title)
	fmt.Println("Body:", string(page2.Body))
}
