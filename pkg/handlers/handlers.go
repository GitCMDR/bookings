package handlers

import (
	"github.com/GitCMDR/go-bookings/pkg/config"
	"github.com/GitCMDR/go-bookings/pkg/models"
	"github.com/GitCMDR/go-bookings/pkg/render"
	"net/http"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) { // declare handler
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{}) // each handler will be mapped to a single gohtml template
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) { // declare handler
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["Test"] = "Hello, I'm context data"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
