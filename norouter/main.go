package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"
)

func main() {

	handler := &app{
		APIHandler: &apiHandler{
			UserHandler: &userHandler{
				TimeHandler: &timehandler{},
			},
		},
	}
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(ln, handler)
	log.Println("Listening on port 8000")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("Shutting down server")
	if err := ln.Close(); err != nil {
		log.Fatalf("Could not shut down server: %v", err)
	}

}

// url helper takes a url string and a segent and returns the index of the segment
// starting with a slash and return the starting index and all url parts
func urlHelper(url, pathSeg string) []string {
	url = path.Clean(url + "/")
	if url == "/" {
		return []string{""}
	}
	segs := strings.Split(url, "/")

	for i, seg := range segs {
		if seg == pathSeg {
			return segs[i:]
		}

	}
	return []string{}

}

type app struct {
	APIHandler *apiHandler
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSeg := ""
	url := path.Clean(r.URL.Path)
	segs := urlHelper(url, pathSeg)
	if len(segs) <= 0 {
		w.WriteHeader(404)
		return
	}
	if len(segs) == 1 {
		w.Write([]byte("Index"))
		return
	}
	switch segs[1] {
	case "api":
		a.APIHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(404)
	}
}

type apiHandler struct {
	UserHandler *userHandler
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSeg := "api"
	url := path.Clean(r.URL.Path)
	segs := urlHelper(url, pathSeg)
	if len(segs) <= 1 {
		w.WriteHeader(404)
		return
	}
	if segs[0] != pathSeg {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch segs[1] {
	case "user":
		a.UserHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(404)
	}
}

type userHandler struct {
	TimeHandler *timehandler
}

func (u *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSeg := "user"
	url := path.Clean(r.URL.Path)
	segs := urlHelper(url, pathSeg)
	if len(segs) <= 0 {
		w.WriteHeader(404)
		return
	}
	if segs[0] != pathSeg {
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
	segs := urlHelper(url, pathSeg)
	if segs[0] != pathSeg {

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
