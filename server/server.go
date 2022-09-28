package server

import (
	"encoding/json"
	"net/http"

	"github.com/irmatov/companies/service"
	"github.com/irmatov/companies/types"
	"github.com/julienschmidt/httprouter"
)

const (
	companiesPrefix = "/companies/"
)

var knownFields = []string{"id", "name", "code", "country", "website", "phone"}

type server struct {
	svc service.Companies
	mux http.Handler
}

func New(storage types.Storage) http.Handler {
	router := httprouter.New()
	s := &server{*service.New(storage), router}
	router.GET(companiesPrefix, s.getMany)
	router.GET(companiesPrefix+":id", s.getSingle)
	router.POST(companiesPrefix, s.create)
	router.DELETE(companiesPrefix+":id", s.delete)
	router.PUT(companiesPrefix+":id", s.update)
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJson(w http.ResponseWriter, statusCode int, content interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(content)
	if err != nil {
		// should be logged/ handled somehow in production
		panic(err)
	}
	w.Write(b)
}
