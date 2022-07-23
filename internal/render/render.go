package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/nanmenkaimak/bookings/internal/config"
	"github.com/nanmenkaimak/bookings/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package 
func NewRenderer(a *config.AppConfig){
	app = a
}

// add data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData{
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// render templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// create template cache 
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("cannot get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)
	// render the template
	_,err := buf.WriteTo(w)
	if err != nil{
		log.Println(err)
		return err
	}

	return nil
}

// creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error){

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates)) // to get all files named *.page.html
	if err != nil{
		return myCache, err
	}

	// range through all files names *.pages.html
	for _, page := range pages{
		name := filepath.Base(page) // return name of file 
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil{
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil{
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil{
				return myCache, err
			}
		}
		myCache[name] = ts;
	}
	return myCache, nil
}