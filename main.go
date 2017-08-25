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

const version = 4.1

func main() {
	//p := profile.Start(profile.CPUProfile, profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)

	log.Infof("Starting version: %f", version)
	router := routing.New()

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
	router.Get("/users/<id>", ParseID, ConnKeepAlive, User)
	router.Get("/users/<id>/visits", ParseID, VisitsByUser)
	router.Post("/users/new", ConnClose, UserCreate)
	router.Post("/users/<id>", ParseID, ConnClose, UserUpdate)

	router.Get("/locations/<id>/avg", ParseID, ConnKeepAlive, LocationAvg)
	router.Get("/locations/<id>", ParseID, ConnKeepAlive, Location)
	router.Post("/locations/new", ConnClose, LocationCreate)
	router.Post("/locations/<id>", ParseID, ConnClose, LocationUpdate)

	router.Get("/visits/<id>", ParseID, ConnKeepAlive, Visit)
	router.Post("/visits/new", ConnClose, VisitCreate)
	router.Post("/visits/<id>", ParseID, ConnClose, VisitUpdate)

	router.NotFound(ConnKeepAlive, NotFound)
}
