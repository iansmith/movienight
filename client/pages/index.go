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
	movieParent = s5.NewHtmlId("div", "movie-parent")

	row     = s5.NewCssClass("row")
	odd     = s5.NewCssClass("odd")
	even    = s5.NewCssClass("even")
	vcenter = s5.NewCssClass("vcenter")
	col2    = s5.NewCssClass("col-md-2")
	col4    = s5.NewCssClass("col-md-4")
)

type movieModel struct {
	s5.ModelName
	posterURL s5.StringAttribute
	title     s5.StringAttribute
	plot      s5.StringAttribute
	nominator s5.StringAttribute
}

func newMovieModel(detail *wire.IMDBDetail, name string) *movieModel {
	result := &movieModel{
		posterURL: s5.NewStringSimple(shared.URLGen.Poster(detail.ImdbID)),
		title:     s5.NewStringSimple(detail.Title),
		plot:      s5.NewStringSimple(detail.Plot),
		nominator: s5.NewStringSimple(name),
	}
	result.ModelName = s5.NewModelName(result)
	return result
}

//we are the same if we have the same ID
func (m movieModel) Equal(e s5.Equaler) bool {
	if e == nil {
		return false
	}
	other := e.(*movieModel)
	return m.Id() == other.Id()
}

type mainPage struct {
	*uilib.PageWithFeedback
	*s5.Collection
}

func (m *mainPage) loggedIn(self *wire.UserRecord) {
	loginParent.Dom().SetText(self.FirstName)
	newMovie.Dom().SetText("New Movie")
	newMovie.Dom().SetCss("cursor", "pointer")
	newMovie.Dom().On(s5.CLICK, func(e jquery.Event) {
		e.PreventDefault()
		uilib.SetNewPage(shared.URLGen.NewMovie())
	})
	var movies []*wire.Movie
	content, errCh := s5.AjaxIndex(&movies, shared.URLGen.MovieResource())
	go func() {
		select {
		case <-content:
			for _, movie := range movies {
				var detail wire.IMDBDetail
				detailCh := make(chan interface{})
				detailErrCh := make(chan s5.AjaxError)

				if err := s5.AjaxRawChannels(&detail, "", detailCh, detailErrCh, "GET",
					shared.URLGen.MovieDetails(movie.ImdbId, true), nil); err != nil {
					m.PageWithFeedback.DisplayFeedback("Unable to connect?", uilib.Danger)
					return
				}
				go func() {
					select {
					case <-detailCh:
						if detail.Response == "False" {
							return
						}
						model := newMovieModel(&detail, movie.Nominator.FirstName+" "+movie.Nominator.LastName)
						m.Collection.Add(model)
					case err := <-detailErrCh:
						m.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
					}
				}()

			}
		case err := <-errCh:
			m.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
		}
	}()

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

	uilib.GetSelf(m.loggedIn, m.notLoggedIn)
}

func newMainPage() *mainPage {
	result := &mainPage{
	//don't init the PageWithFeedback yet, DOM not loaded
	}
	result.Collection = s5.NewList(result)
	return result
}

func (self *mainPage) Add(length int, newObj s5.Model) {
	model := newObj.(*movieModel)
	bkgrd := odd
	if length%2 == 0 {
		bkgrd = even
	}
	tree := s5.DIV(
		s5.Class(row),
		s5.Class(vcenter),
		s5.Class(bkgrd),
		s5.DIV(
			s5.Class(col2),
			s5.H4(
				s5.TextEqual(model.title),
			),
		),
		s5.DIV(
			s5.Class(col2),
			s5.IMG(
				s5.HtmlAttrConstant(s5.HEIGHT, "200"),
				s5.HtmlAttrConstant(s5.WIDTH, "auto"),
				s5.HtmlAttrEqual(s5.SRC, model.posterURL),
			),
		),
		s5.DIV(
			s5.Class(col4),
			s5.H6(
				s5.TextEqual(model.plot),
			),
		),
		s5.DIV(
			s5.Class(col2),
			s5.H6(
				s5.TextEqual(model.nominator),
			),
		),
	).Build()
	movieParent.Dom().Append(tree)
	print("add reached", length, model)
}
func (self *mainPage) Remove(ignored int, oldObj s5.Model) {
	panic("remove not implemented")
}

func main() {
	p := newMainPage()
	s5.Main(p)
}
