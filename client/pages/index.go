// +build js
package main

import (
	"net/http"

	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/client/uilib"
	"github.com/iansmith/movienight/shared"
	"github.com/iansmith/movienight/wire"
)

var (
	loginParent = s5.NewHtmlId("span", "login-parent")
	newMovie    = s5.NewHtmlId("span", "new-movie")
)

type mainPage struct {
	*uilib.PageWithFeedback
}

func (m *mainPage) loggedIn(self *wire.UserRecord) {
	loginParent.Dom().SetText(self.FirstName)
	newMovie.Dom().SetText("New Movie")
	newMovie.Dom().SetCss("cursor", "pointer")
	newMovie.Dom().On(s5.CLICK, func(e jquery.Event) {
		e.PreventDefault()
		uilib.SetNewPage(shared.URLGen.NewMovie())
	})
}

func (m *mainPage) notLoggedIn(code int, message string) {
	loginParent.Dom().SetText("Log In")
	loginParent.Dom().SetCss("cursor", "pointer")
	loginParent.Dom().On(s5.CLICK, func(e jquery.Event) {
		e.PreventDefault()
		uilib.SetNewPage(shared.URLGen.Login())
	})
	if code != http.StatusUnauthorized {
		m.PageWithFeedback.DisplayFeedback(message, uilib.Danger)
	}
}

func (m *mainPage) Start() {
	m.PageWithFeedback = uilib.NewPageWithFeedback()
	m.PageWithFeedback.DisplayFeedback("howdy", uilib.Warning)

	uilib.GetSelf(m.loggedIn, m.notLoggedIn)
}

func newMainPage() *mainPage {
	return &mainPage{
	//don't init the PageWithFeedback yet, DOM not loaded
	}
}

func main() {
	p := newMainPage()
	s5.Main(p)
}
