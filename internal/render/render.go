package render

import (
	"bytes"
	"github.com/GitCMDR/go-bookings/internal/config"
	"github.com/GitCMDR/go-bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) { // define a function to render // templates, this will be used by all page handlers

	var tc map[string]*template.Template

	if app.UseCache { // if dev mode use cache is set to false
		// pull template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl] // check for key with ok notation
	if !ok {
		log.Fatal("could not get template from template cache") // if key not found stop the app
	}

	buf := new(bytes.Buffer) // create a buffer

	td = AddDefaultData(td, r) // add default data to template data map

	_ = t.Execute(buf, td)   // execute template to check for errors and write to buffer
	_, err := buf.WriteTo(w) // write the template to writer
	if err != nil {
		log.Println(err) // print error to logs
	}

	if err != nil { // check for errors
		log.Println("error is", err)
		return
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{} // define a map of string:pointer2template

	pages, err := filepath.Glob("./templates/*.page.gohtml") // go to folder and find matches to string, get a slice of strings

	if err != nil { // check for errors
		return templateCache, err
	}

	for _, page := range pages { // iterate through all the templates
		name := filepath.Base(page) // instead of getting whole file path, just get file name
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		matches, err := filepath.Glob("./templates/*.layout.gohtml") // check for layouts
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = ts
	}

	return templateCache, nil

}
