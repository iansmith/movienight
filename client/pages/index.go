// +build js
package main

import (
	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/client/uilib"
)

type mainPage struct {
	*uilib.PageWithFeedback
}

func (m *mainPage) Start() {
	m.PageWithFeedback = uilib.NewPageWithFeedback()
	m.PageWithFeedback.DisplayFeedback("howdy", uilib.Warning)
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
