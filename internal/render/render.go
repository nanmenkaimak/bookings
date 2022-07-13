package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/nanmenkaimak/bookings/internal/config"
	"github.com/nanmenkaimak/bookings/internal/models"
)

//var functions = template.FuncMap{}

var app *config.AppConfig
// NewTemplates sets the config for the template package 
func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData{
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData){
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
		log.Fatal("Could not get template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)
	// render the template
	_,err := buf.WriteTo(w)
	if err != nil{
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error){

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html") // to get all files named *.page.html
	if err != nil{
		return myCache, err
	}

	// range through all files names *.pages.html
	for _, page := range pages{
		name := filepath.Base(page) // return name of file 
		ts, err := template.New(name).ParseFiles(page)
		if err != nil{
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil{
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil{
				return myCache, err
			}
		}
		myCache[name] = ts;
	}
	return myCache, nil
}