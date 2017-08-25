package main

import (
	"time"

	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	"github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var dataStore = store.NewStore()
var now time.Time

const version = 4.0

func main() {
	//p := profile.Start(profile.CPUProfile, profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)

	log.Infof("Starting version: %f", version)

	router := routing.New()
	//router.NotFoundHandler = http.HandlerFunc(NotFound)

	now = util.ImportCurrentTimestamp()

	// import static data
	start := time.Now()
	err := util.ImportDataFromZip(dataStore)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Time to load data.zip: %v", time.Since(start))

	// set up routes
	setRouting(router)

	// start http server
	log.Fatal(fasthttp.ListenAndServe(":80", router.HandleRequest))
}

func setRouting(router *routing.Router) {
	router.Get("/users/<id>", ConnKeepAlive, User)
	router.Get("/users/<id>/visits", VisitsByUser)
	router.Post("/users/new", ConnClose, UserCreate)
	router.Post("/users/<id>", ConnClose, UserUpdate)

	router.Get("/locations/<id>/avg", ConnKeepAlive, LocationAvg)
	router.Get("/locations/<id>", ConnKeepAlive, Location)
	router.Post("/locations/new", ConnClose, LocationCreate)
	router.Post("/locations/<id>", ConnClose, LocationUpdate)

	router.Get("/visits/<id>", ConnKeepAlive, Visit)
	router.Post("/visits/new", ConnClose, VisitCreate)
	router.Post("/visits/<id>", ConnClose, VisitUpdate)
}
