package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	s5 "github.com/seven5/seven5"
	//trick godeps to putting things in the Godep directory that the client
	//side needs to compile
	_ "github.com/gopherjs/gopherjs/js"
	_ "github.com/gopherjs/gopherjs/nosync"
	_ "github.com/gopherjs/jquery"
	//this is one is actually used by the resulting binary
	_ "github.com/lib/pq"

	"github.com/iansmith/movienight/server/resource"
	"github.com/iansmith/movienight/shared"
	"github.com/iansmith/movienight/wire"
)

const (
	Name = "movienight" //mostly for env vars
)

//shoved in by the linker
var version string

type movienightConfig struct {
	serveMux *s5.ServeMux
	base     *s5.BaseDispatcher
	pwd      *s5.SimplePasswordHandler
	heroku   *s5.HerokuDeploy
	matcher  *s5.SimpleComponentMatcher
	cm       s5.CookieMapper
	sm       s5.ValidatingSessionManager
}

//returns the config information, makes testing easier
func config() *movienightConfig {
	result := &movienightConfig{}

	result.heroku = s5.NewHerokuDeploy("iggy-movienight", Name)
	result.cm = s5.NewSimpleCookieMapper(Name)
	result.sm = newMovienightValidatingSessionManager()
	result.base = s5.NewBaseDispatcher(result.sm, result.cm)
	result.serveMux = s5.NewServeMux()
	result.pwd = s5.NewSimplePasswordHandler(result.sm, result.cm)

	//what do we do if given empty URL, note we use the SINGULAR here
	homepage := s5.ComponentResult{
		Status: http.StatusOK,
		Path:   "/en/web/index.html",
	}

	log.Printf("[TEST MODE] %v", result.heroku.IsTest())
	result.matcher = s5.NewSimpleComponentMatcher(result.cm, result.sm, "static",
		homepage, result.heroku.IsTest(), &movieComponent{})

	return result
}

func addResourcesAndAuth(conf *movienightConfig) {
	//auth, session management, and "self" are closely tied up
	conf.serveMux.HandleFunc(shared.URLGen.Me(), conf.pwd.MeHandler)
	conf.serveMux.HandleFunc(shared.URLGen.Auth(), conf.pwd.AuthHandler)

	store := conf.heroku.GetQbsStore()
	//this is a GET to get the movie list
	conf.base.ResourceSeparate("movie",
		&wire.Movie{},
		s5.QbsWrapIndex(&resource.MovieResource{}, store),
		nil, //find
		s5.QbsWrapPost(&resource.MovieResource{}, store), //post
		nil, //put
		nil) //delete

	conf.serveMux.Dispatch("/rest/", conf.base)

	//add static files
	conf.serveMux.Handle("/", conf.matcher)
}

func movieproxy(w http.ResponseWriter, req *http.Request) {
	omdb, err := url.Parse("http://www.omdbapi.com")
	if err != nil {
		http.Error(w, "bad url", http.StatusInternalServerError)
		return
	}
	q := omdb.Query()
	q.Add("i", req.URL.Query().Get("i"))
	q.Add("plot", req.URL.Query().Get("plot"))
	q.Add("r", "json")
	log.Printf("[MOVIEPROXY] proxying request for %s-> %v",
		req.URL.Query().Get("i"), "http://www.omdbapi.com?"+q.Encode())

	resp, err := http.Get("http://www.omdbapi.com?" + q.Encode())
	if err != nil {
		http.Error(w, "unable to reach omdb", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}
	var detail wire.IMDBDetail
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&detail); err != nil {
		http.Error(w, "unable to decode body:"+err.Error(), http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(detail); err != nil {
		http.Error(w, "unable to encode body:"+err.Error(), http.StatusInternalServerError)
		return
	}
	//200
}

func posterproxy(w http.ResponseWriter, req *http.Request) {
	omdb, err := url.Parse("http://img.omdbapi.com")
	if err != nil {
		http.Error(w, "bad url", http.StatusInternalServerError)
		return
	}
	q := omdb.Query()
	q.Add("i", req.URL.Query().Get("i"))
	q.Add("apikey", os.Getenv("OMDB_API_KEY"))
	log.Printf("[POSTERPROXY] proxying request for %s-> %v",
		req.URL.Query().Get("i"), "http://img.omdbapi.com?"+q.Encode())
	resp, err := http.Get("http://img.omdbapi.com?" + q.Encode())
	if err != nil {
		http.Error(w, "unable to reach omdb", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}
	for k, v := range resp.Header {
		w.Header().Add(k, v[0])
	}
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("[POSTERPROXY] failed to copy: %v", err)
	}
	//200
}

func main() {
	log.Printf("[VERSION] %s", version)
	log.Printf("[DATABASE] %s", os.Getenv("DATABASE_URL"))
	conf := config()
	addResourcesAndAuth(conf)

	//get the port from the environment
	port := conf.heroku.Port()

	//convenient way to get the version
	conf.serveMux.HandleFunc("/version", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(version))
	})
	conf.serveMux.HandleFunc("/movieproxy", movieproxy)
	conf.serveMux.HandleFunc("/posterproxy", posterproxy)
	//serve up the content forever
	log.Printf("[SERVE] waiting on :%d", port)
	log.Fatalf("%s", http.ListenAndServe(fmt.Sprintf(":%d", port), conf.serveMux))

}
