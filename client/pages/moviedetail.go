package main

import (
	"net/http"
	"strings"

	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/client/uilib"
	"github.com/iansmith/movienight/shared"
	"github.com/iansmith/movienight/wire"
)

type movieDetailPage struct {
	*uilib.PageWithFeedback
}

func newMovieDetailPage() *movieDetailPage {
	return &movieDetailPage{}
}
func (m *movieDetailPage) loggedIn(user *wire.UserRecord, detail *wire.IMDBDetail) {

}
func (m *movieDetailPage) notLoggedIn(code int, msg string, detail *wire.IMDBDetail) {
	if code != http.StatusUnauthorized {
		m.PageWithFeedback.DisplayFeedback(msg, uilib.Danger)
		return
	}
	//it's just unauthed
}

func (m *movieDetailPage) Start() {
	m.PageWithFeedback = uilib.NewPageWithFeedback()

	p := uilib.CurrentURLPath()
	parts := strings.Split(p, "/")
	if len(parts) != 3 || parts[1] != "movie" {
		m.PageWithFeedback.DisplayFeedback("Can't understand URL: "+p, uilib.Danger)
		return
	}

	var detail wire.IMDBDetail
	contentCh := make(chan interface{})
	errCh := make(chan s5.AjaxError)

	if err := s5.AjaxRawChannels(&detail, "", contentCh, errCh, "GET", shared.URLGen.ProxyMovie(parts[2], true), nil); err != nil {
		m.PageWithFeedback.DisplayFeedback("Unable to connect?", uilib.Danger)
		return
	}
	go func() {
		select {
		case <-contentCh:
			if detail.Response != "True" {
				m.PageWithFeedback.DisplayFeedback("Movie Not Found", uilib.Danger)
				return
			}
		case err := <-errCh:
			m.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
			return
		}
	}()

	uilib.GetSelf(func(user *wire.UserRecord) {
		m.loggedIn(user, &detail)
	}, func(code int, msg string) {
		m.notLoggedIn(code, msg, &detail)
	})
}

func main() {
	s5.Main(newMovieDetailPage())
}
