package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var err500 string = "500 Internal Server Error"
var err404 string = "404 This page not found"
var err400 string = "400 Bad Request"
var selectedBanner string = "standard"
var firstRun int = 0

type Page struct {
	Banner string
	Body   []byte
	Output []byte
}

func (p *Page) save() error {
	filename := "test.txt"

	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(banner string) (*Page, error) {
	filename := "test.txt"
	body, err := ioutil.ReadFile(filename)
	if firstRun == 0 {
		body = []byte("")
	}
	output, err := ioutil.ReadFile("output.txt")
	firstRun = 1
	if err != nil {
		return nil, err
	}
	return &Page{Banner: banner, Body: body, Output: output}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// firstRun = 0
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	banner := r.FormValue("banners")
	selectedBanner = banner
	p, err := loadPage(selectedBanner)

	if err != nil {
		fmt.Println(err)
		p = &Page{}
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	banner := r.FormValue("banners")
	selectedBanner = banner
	p := &Page{Banner: selectedBanner, Body: []byte(body)}
	p.save()
	input := body[:(len(body))]
	if !isValid(input) {
		errorHandler(w, r, 400)
		return
	}
	out, err := asciify(input, banner)
	if err != nil {
		errorHandler(w, r, 500)
		return
	}
	ioutil.WriteFile("output.txt", []byte(out), 0600)
	http.Redirect(w, r, "/", http.StatusFound)
}

func isValid(s string) bool {
	for _, letter := range s {
		if letter < 32 || letter > 126 {
			return false
		}
	}
	return true
}

func asciify(args ...string) (string, error) {
	cmd := exec.Command("./ascii", args...)
	cmd.Dir = "ascii"

	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(data))
	return string(data), err
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	t, _ := template.ParseFiles("error.html")
	p := &Page{}
	if status == http.StatusNotFound {
		p = &Page{Banner: err404}
	} else if status == 500 {
		p = &Page{Banner: err500}
	} else if status == 400 {
		p = &Page{Banner: err400}
	}
	t.Execute(w, p)
}

// Left to handle:

// 5. +Is there a test file for this code?
// 6. +Are the tests checking each possible case?
// 7. +Are the instructions in the website clear?
// 8. +Does the project run using an API?

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
