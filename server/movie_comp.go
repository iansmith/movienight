package main

import (
	"net/http"
	"path/filepath"

	s5 "github.com/seven5/seven5"
)

type movieComponent struct {
	//statelesss
}

func (mc *movieComponent) Page(pb s5.PBundle, path []string, slashTerminated bool) s5.ComponentResult {
	if len(path) != 1 {
		return s5.ComponentResult{
			Status:  http.StatusNotFound,
			Message: "not found",
		}
	}

	return s5.ComponentResult{
		Status: http.StatusOK,
		Path:   "/" + filepath.Join("moviedetail.html"),
	}

}

//part of the fixed URL space, not including preceding slash
func (mc *movieComponent) UrlPrefix() string {
	return "movie"
}
