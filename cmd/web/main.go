package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hb-otudis/bookings/pkg/config"
	"github.com/hb-otudis/bookings/pkg/handlers"
	"github.com/hb-otudis/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig

var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("can't create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	fmt.Printf("Starting Application on %s", portNumber)

	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
