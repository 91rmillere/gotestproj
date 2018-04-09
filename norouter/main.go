package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

func main() {

	handler := &app{&userHandler{&timehandler{}}}

	http.ListenAndServe(":3030", handler)

}

func urlHelper(url, pathSeg string) (int, []string) {
	u := path.Clean(url)
	path := fmt.Sprintf("/%v", pathSeg)
	index := strings.Index(u, path)
	if index < 0 {
		return index, []string{}
	}
	segs := strings.Split(url[index+1:], "/")
	return index, segs
}

type app struct {
	UserHandler *userHandler
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.UserHandler.ServeHTTP(w, r)
}

type userHandler struct {
	TimeHandler *timehandler
}

func (u *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSeg := "user"
	url := path.Clean(r.URL.Path)
	index, segs := urlHelper(url, pathSeg)
	if index < 0 || segs[0] != pathSeg {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(segs) == 1 {
		w.Write([]byte("hello user"))
		return
	}

	switch segs[1] {
	case "time":
		u.TimeHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(404)
	}

}

type timehandler struct{}

func (t *timehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSeg := "time"
	url := path.Clean(r.URL.Path)
	index, segs := urlHelper(url, pathSeg)
	if index < 0 || segs[0] != pathSeg {

		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(segs) == 1 {
		w.Write([]byte(time.Now().String()))
		return
	}
	switch segs[1] {
	case "est":
		w.Write([]byte(time.Now().Local().String()))
	case "utc":
		w.Write([]byte(time.Now().UTC().String()))
	case "unix":
		w.Write([]byte(string(time.Now().Unix())))
	default:
		w.WriteHeader(404)
	}
}
