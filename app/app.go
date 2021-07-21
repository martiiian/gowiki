package app

import (
	"fmt"
	"gowiki/consts"
	"gowiki/models"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles(getViewFileName("edit"), getViewFileName("view"), getViewFileName("home"), getViewFileName("404")))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getViewFileName(name string) string {
	return fmt.Sprintf("%s/%s.html", consts.ViewDir, name)
}

func loadPage(title string) (*models.Page, error) {
	filename := fmt.Sprintf("%s/%s.txt", consts.ArticlesDir, title)

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &models.Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, templateName string, page *models.Page) {
	err := templates.ExecuteTemplate(w, templateName + ".html", page)
	if err != nil {
		showHttpError(w, err)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", page)
}

func EditHandler(w http.ResponseWriter, _ *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &models.Page{Title: title}
	}
	renderTemplate(w, "edit", page)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &models.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		showHttpError(w, err)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func MainHandler(w http.ResponseWriter, _ *http.Request) {
	var emptyData interface{}
	err := templates.ExecuteTemplate(w, "home.html", emptyData)
	if err != nil {
		showHttpError(w, err)
	}
}

func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	var emptyData interface{}
	err := templates.ExecuteTemplate(w, "404.html", emptyData)
	if err != nil {
		showHttpError(w, err)
	}
}

func showHttpError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func MakeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/not_found", http.StatusNotFound)
			return
		}

		fn(w, r, m[2])
	}
}
