package main

import (
	"encoding/gob"
	"github.com/GitCMDR/go-bookings/internal/config"
	"github.com/GitCMDR/go-bookings/internal/handlers"
	"github.com/GitCMDR/go-bookings/internal/helpers"
	"github.com/GitCMDR/go-bookings/internal/models"
	"github.com/GitCMDR/go-bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig // create variable for app config
var session *scs.SessionManager // create a variable of type pointer to scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}

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

func run() error {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// create loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()              // enable sessions support
	session.Lifetime = 24 * time.Hour // kill session in 24 hours
	session.Cookie.Persist = true     // allow session to persist after browse closure
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session // share session to all packages in the app via config file

	tc, err := render.CreateTemplateCache() // create template cache when app runs

	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc // store template cache in app config
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}
