package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

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
	output, err := ioutil.ReadFile("output.txt")
	if err != nil {
		return nil, err
	}
	return &Page{Banner: banner, Body: body, Output: output}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	banner := r.FormValue("banners")
	p, err := loadPage(banner)
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
	p := &Page{Banner: banner, Body: []byte(body)}
	p.save()
	out := asciify(body[:(len(body))], banner)
	ioutil.WriteFile("output.txt", []byte(out), 0600)
	http.Redirect(w, r, "/", http.StatusFound)
}

func asciify(args ...string) string {
	cmd := exec.Command("./ascii", args...)
	cmd.Dir = "ascii"

	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(data))
	return string(data)
}

// Left to handle:
// 1. Check the input for correctness
// 2. Avoid 404 status
// 3. Avoid 400 status
// 4. Avoid 500 status
// 5. +Is there a test file for this code?
// 6. +Are the tests checking each possible case?
// 7. +Are the instructions in the website clear?
// 8. +Does the project run using an API?
// 9. Make it empty for the very first run
// 10. Make it remember the choice for banner

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
