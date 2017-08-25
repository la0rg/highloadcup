package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/la0rg/highloadcup/store"
	"github.com/la0rg/highloadcup/util"
	"github.com/mailru/easyjson"
	"github.com/qiangxue/fasthttp-routing"

	"github.com/la0rg/highloadcup/model"
)

var emptyObject = []byte("{}")
var ErrParse = errors.New("Could not parse Id from request")

const (
	FromDate   = "fromDate"
	ToDate     = "toDate"
	FromAge    = "fromAge"
	ToAge      = "toAge"
	Gender     = "gender"
	Country    = "country"
	ToDistance = "toDistance"
)

// User returns a user by id
func User(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}
	user, ok := dataStore.GetUserByID(id)
	if ok {
		err = writeStructAsJSON(ctx, user)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

// UserUpdate update user entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func UserUpdate(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}

	bytes := ctx.PostBody()
	errNull := util.ContainsNull(bytes)
	if errNull {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	var user model.User
	errParse := easyjson.Unmarshal(bytes, &user)
	if errParse != nil {
		user = model.User{}
	}
	err = dataStore.UpdateUserByID(id, user)
	// 404 is a higher priority than 400
	if err != nil {
		if err == store.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return nil
	}
	if errParse != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// UserCreate create user entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func UserCreate(ctx *routing.Context) error {
	var user model.User
	err := easyjson.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	err = dataStore.AddUser(user)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// Location returns a location by id
func Location(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}
	location, ok := dataStore.GetLocationByID(id)
	if ok {
		err = writeStructAsJSON(ctx, location)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

// Location returns a location by id
func LocationAvg(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}

	var fromDate, toDate, fromAge, toAge *int64
	var gender *string

	args := ctx.QueryArgs()
	if args.Has(FromDate) {
		i64, err := strconv.ParseInt(string(args.Peek(FromDate)), 10, 64)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		fromDate = &i64
	}
	if args.Has(ToDate) {
		i64, err := strconv.ParseInt(string(args.Peek(ToDate)), 10, 64)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		toDate = &i64
	}
	if args.Has(FromAge) {
		i, err := strconv.Atoi(string(args.Peek(FromAge)))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		date := now.AddDate(-i, 0, 0).Unix()
		fromAge = &date
	}
	if args.Has(ToAge) {
		i, err := strconv.Atoi(string(args.Peek(ToAge)))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		date := now.AddDate(-i, 0, 0).Unix()
		toAge = &date
	}
	if args.Has(Gender) {
		str := string(args.Peek(Gender))
		gender = &str
		if !util.IsGender(*gender) {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
	}

	avg, ok := dataStore.GetLocationAvg(id, fromDate, toDate, fromAge, toAge, gender)
	if ok {
		//avg = util.RoundPlus(avg, 5)
		err = writeStructAsJSON(ctx, model.Avg{Value: avg})
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

// LocationUpdate update location entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func LocationUpdate(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}

	bytes := ctx.PostBody()
	errNull := util.ContainsNull(bytes)
	if errNull {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	var location model.Location
	errParse := easyjson.Unmarshal(bytes, &location)
	if errParse != nil {
		location = model.Location{}
	}

	err = dataStore.UpdateLocationByID(id, location)
	if err != nil {
		if err == store.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return nil
	}
	if errParse != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// LocationCreate create location entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func LocationCreate(ctx *routing.Context) error {
	bytes := ctx.PostBody()

	var location model.Location
	err := easyjson.Unmarshal(bytes, &location)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	err = dataStore.AddLocation(location)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// Visit returns a visit by id
func Visit(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}
	visit, ok := dataStore.GetVisitByID(id)
	if ok {
		err = writeStructAsJSON(ctx, visit)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

func VisitsByUser(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}

	var fromDate, toDate *int64
	var country *string
	var toDistance *int32
	args := ctx.QueryArgs()
	if args.Has(FromDate) {
		i64, err := strconv.ParseInt(string(args.Peek(FromDate)), 10, 64)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		fromDate = &i64
	}
	if args.Has(ToDate) {
		i64, err := strconv.ParseInt(string(args.Peek(ToDate)), 10, 64)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		toDate = &i64
	}
	if args.Has(Country) {
		str := string(args.Peek(Country))
		country = &(str)
		if !util.OnlyLetters(*country) {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
	}
	if args.Has(ToDistance) {
		i64, err := strconv.ParseInt(string(args.Peek(ToDistance)), 10, 32)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return nil
		}
		i32 := int32(i64)
		toDistance = &i32
	}

	visits, ok := dataStore.GetVisitsByUserID(id, fromDate, toDate, country, toDistance)
	if ok {
		err = writeStructAsJSON(ctx, visits)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
		return nil
	}
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

// VisitUpdate update user entity
// success - 200 with body {}
// id is not found - 404
// incorrect request - 400
func VisitUpdate(ctx *routing.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return nil
	}

	bytes := ctx.PostBody()
	errNull := util.ContainsNull(bytes)
	if errNull {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}

	var visit model.Visit
	errParse := easyjson.Unmarshal(bytes, &visit)
	if errParse != nil {
		visit = model.Visit{}
	}
	err = dataStore.UpdateVisitByID(id, visit)
	// 404 is a higher priority than 400
	if err != nil {
		if err == store.ErrDoesNotExist {
			ctx.SetStatusCode(http.StatusNotFound)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		return nil
	}
	if errParse != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// VisitCreate create location entity
// success - 200 with body {}
// already exist - 400
// incorrect request - 400
func VisitCreate(ctx *routing.Context) error {
	bytes := ctx.PostBody()

	var visit model.Visit
	err := easyjson.Unmarshal(bytes, &visit)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	err = dataStore.AddVisit(visit)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return nil
	}
	ctx.SetBody(emptyObject)
	return nil
}

// NotFound custom request handler for non-found requests
func NotFound(ctx *routing.Context) error {
	ctx.SetStatusCode(http.StatusNotFound)
	return nil
}

func parseID(ctx *routing.Context) (int32, error) {
	id64, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	id := int32(id64)
	if err != nil {
		return 0, ErrParse
	}
	return id, nil
}

func writeStructAsJSON(ctx *routing.Context, object easyjson.Marshaler) error {
	b, err := easyjson.Marshal(object)
	if err != nil {
		return err
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))

	ctx.SetBody(b)
	return err
}

func ConnKeepAlive(ctx *routing.Context) error {
	ctx.Response.Header.Set("Connection", "Keep-Alive")
	return nil
}

func ConnClose(ctx *routing.Context) error {
	ctx.Response.Header.Set("Connection", "Close")
	return nil
}
