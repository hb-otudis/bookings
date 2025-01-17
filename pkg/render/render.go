package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/hb-otudis/bookings/pkg/config"
	"github.com/hb-otudis/bookings/pkg/models"
)

var app *config.AppConfig

// New template sets config for new template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate render templates using html
func RenderTemplate(w http.ResponseWriter, tmpl string, td models.TemplateData) {
	// Get template cache from App config
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("err")
	}

	buf := new(bytes.Buffer)

	td = *AddDefaultData(&td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, err
}
