package main

import (
	"github.com/GitCMDR/go-bookings/pkg/config"
	"github.com/GitCMDR/go-bookings/pkg/handlers"
	"github.com/GitCMDR/go-bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig // create variable for app config
var session *scs.SessionManager // create a variable of type pointer to scs.SessionManager

// main is the main application function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()              // enable sessions support
	session.Lifetime = 24 * time.Hour // kill session in 24 hours
	session.Cookie.Persist = true     // allow session to persist after browse closure
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session // share session to all packages in the app via config file

	tc, err := render.CreateTemplateCache() // create template cache when app runs
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc // store template cache in app config
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	log.Printf("Starting application on port %v", portNumber)
	//_ = http.ListenAndServe(portNumber, nil) // if error throw error away

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

// 🔥 in order to run the app as a whole do go run *.go instead of go run main.go
