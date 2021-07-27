package render

import (
	"encoding/gob"
	"github.com/GitCMDR/go-bookings/internal/config"
	"github.com/GitCMDR/go-bookings/internal/models"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	session = scs.New()              // enable sessions support
	session.Lifetime = 24 * time.Hour // kill session in 24 hours
	session.Cookie.Persist = true     // allow session to persist after browse closure
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session // share session to all packages in the app via config file

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct {}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write (b []byte) (int, error) {
	length := len(b)
	return length, nil
}