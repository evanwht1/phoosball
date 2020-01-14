package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  template.HTML
}

var templates = template.Must(template.ParseFiles("template.html"))

func setContentType(w http.ResponseWriter, title string) string {
	fileType := append(strings.Split(title, "."), "")[1]
	switch fileType {
	case "css":
		w.Header().Set("content-type", "text/css")
		return "text/css"
	case "svg":
		w.Header().Set("content-type", "image/svg+xml")
		return "image/svg+xml"
	case "jpg":
		w.Header().Set("content-type", "image/jpg")
		return "image/jpg"
	case "js":
		w.Header().Set("content-type", "text/javascript")
		return "text/javascript"
	default:
		w.Header().Set("content-type", "text/html")
		return "text/html"
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadHtml(fileName string) (*Page, error) {
	var (
		body []byte
		err  error
	)
	if templateName := strings.Split(fileName, ".")[0] + "_template.html"; fileExists(templateName) {
		body, err = ioutil.ReadFile(templateName)
	} else if fileExists(fileName) {
		body, err = ioutil.ReadFile(fileName)
	}
	if err != nil {
		return nil, err
	}
	return &Page{Title: fileName, Body: template.HTML(body)}, nil
}

func renderTemplate(w http.ResponseWriter, p *Page) {
	err := templates.ExecuteTemplate(w, "template.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	if len(title) == 0 {
		title = "index"
	}
	switch fileType := setContentType(w, title); fileType {
	case "text/html", "text/plain":
		p, err := loadHtml(title + "_template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			renderTemplate(w, p)
		}
	default:
		p, err := ioutil.ReadFile(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "%s", p)
		}
	}
}

func addNewPlayer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.PostForm) > 0 {
		p, err := template.ParseFiles("account.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var buff bytes.Buffer
		if err = p.Execute(&buff, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		renderTemplate(w, &Page{Title: "Account", Body: template.HTML(buff.String())})
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/save_player", addNewPlayer)
	// http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	log.Fatal(http.ListenAndServe(":3032", nil))
}
