package main

import (
	"bytes"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var dataStore = store.NewStore()
var now time.Time

const idLabel = "id"
const version = 6.0

func main() {
	//defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	debug.SetGCPercent(50)
	log.Infof("Starting version: %f", version)

	now = util.ImportCurrentTimestamp()
	// import static data
	start := time.Now()
	err := util.ImportDataFromZip(dataStore)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Time to load data.zip: %v", time.Since(start))
	runtime.GC()

	// start http server
	log.Fatal(fasthttp.ListenAndServe(":80", manualRouting()))
}

func manualRouting() fasthttp.RequestHandler {
	delim := []byte{'/'}
	users := []byte("users")
	locations := []byte("locations")
	visits := []byte("visits")
	new := []byte("new")
	avg := []byte("avg")

	return func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		if path[0] != delim[0] {
			NotFound(ctx)
			return
		}
		path = path[1:]

		ConnKeepAlive(ctx)
		ctx.SetStatusCode(http.StatusOK)
		if ctx.IsGet() {
			//parts := bytes.Split(path, delim)
			l := bytes.Count(path, delim)
			var part1, part2, part3 []byte
			if l == 1 {
				i1 := bytes.IndexByte(path, delim[0])
				part1 = path[:i1]
				part2 = path[i1+1:]
			} else if l == 2 {
				i1 := bytes.IndexByte(path, delim[0])
				i2 := i1 + 1 + bytes.IndexByte(path[i1+1:], delim[0])
				part1 = path[:i1]
				part2 = path[i1+1 : i2]
				part3 = path[i2+1:]
			} else {
				NotFound(ctx)
				return
			}

			if bytes.Equal(part1, users) {
				ctx.SetUserValue(idLabel, part2)
				if l == 2 && bytes.Equal(part3, visits) {
					VisitsByUser(ctx)
				} else {
					User(ctx)
				}
			} else if bytes.Equal(part1, locations) {
				ctx.SetUserValue(idLabel, part2)
				if l == 2 && bytes.Equal(part3, avg) {
					LocationAvg(ctx)
				} else {
					Location(ctx)
				}
			} else if bytes.Equal(part1, visits) {
				if l == 2 {
					NotFound(ctx)
					return
				}
				ctx.SetUserValue(idLabel, part2)
				Visit(ctx)
			}
			return
		}
		if ctx.IsPost() {
			ConnClose(ctx)
			l := bytes.Count(path, delim)
			var part1, part2 []byte
			if l == 1 {
				i1 := bytes.IndexByte(path, delim[0])
				part1 = path[:i1]
				part2 = path[i1+1:]
			} else {
				NotFound(ctx)
				return
			}

			if bytes.Equal(part1, users) {
				if bytes.Equal(part2, new) {
					UserCreate(ctx)
				} else {
					ctx.SetUserValue(idLabel, part2)
					UserUpdate(ctx)
				}
			} else if bytes.Equal(part1, locations) {
				if bytes.Equal(part2, new) {
					LocationCreate(ctx)
				} else {
					ctx.SetUserValue(idLabel, part2)
					LocationUpdate(ctx)
				}
			} else if bytes.Equal(part1, visits) {
				if bytes.Equal(part2, new) {
					VisitCreate(ctx)
				} else {
					ctx.SetUserValue(idLabel, part2)
					VisitUpdate(ctx)
				}
			}
			return
		}

		NotFound(ctx)
	}
}
