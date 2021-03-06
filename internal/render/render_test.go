package render

import (
	"github.com/GitCMDR/go-bookings/internal/models"
	"net/http"
	"testing"
)

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser when cache is turned off")
	}

	err = RenderTemplate(&ww, r, "some-other-page-that-doesn't-exist.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that doesn't exist 👻")
	}

	app.UseCache = true

	err = RenderTemplate(&ww, r, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser when cache is turned on")
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) // create session on fake request context

	r = r.WithContext(ctx)
	return r, nil
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}


}