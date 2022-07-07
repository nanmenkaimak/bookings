package handlers

import (
	"net/http"

	"github.com/nanmenkaimak/bookings/pkg/config"
	"github.com/nanmenkaimak/bookings/pkg/models"
	"github.com/nanmenkaimak/bookings/pkg/render"
)


var Repo *Repository // repository used by handlers

type Repository struct { // repository type
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository { //creates new repository
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) { // sets repository for the handlers
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello, meshok"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}


