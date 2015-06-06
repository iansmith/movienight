package main

import (
	"fmt"
	"log"
	"net/http"
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
		nil, //post
		nil, //put
		nil) //delete

	//add static files
	conf.serveMux.Handle("/", conf.matcher)
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
	//serve up the content forever
	log.Printf("[SERVE] waiting on :%d", port)
	log.Fatalf("%s", http.ListenAndServe(fmt.Sprintf(":%d", port), conf.serveMux))

}
