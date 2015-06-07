package main

import (
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/client/uilib"
	"github.com/iansmith/movienight/shared"
	"github.com/iansmith/movienight/wire"
)

var (
	submitButton = s5.NewHtmlId("a", "submit")
	imdbId       = s5.NewHtmlId("input", "imdb_id")
	posterParent = s5.NewHtmlId("span", "poster-parent")
	posterTitle  = s5.NewHtmlId("span", "poster-title")
	disabled     = s5.NewCssClass("disabled")
)

type newMoviePage struct {
	*uilib.PageWithFeedback
	havePoster s5.BooleanAttribute
	movieId    string
	userUdid   string
}

func NewNewMoviePage() *newMoviePage {
	return &newMoviePage{
		havePoster: s5.NewBooleanSimple(false),
	}
}

func (n *newMoviePage) fetchMovie(id string) {
	var detail wire.IMDBDetail
	contentCh := make(chan interface{})
	errCh := make(chan s5.AjaxError)

	if err := s5.AjaxRawChannels(&detail, "", contentCh, errCh, "GET", "/movieproxy?i="+id, nil); err != nil {
		n.PageWithFeedback.DisplayFeedback("Unable to connect?", uilib.Danger)
		return
	}
	go func() {
		select {
		case <-contentCh:
			img := detail.Poster
			if detail.Response == "False" {
				img = "/fixed/questionmark.jpg"
				n.havePoster.Set(false)
			} else {
				n.havePoster.Set(true)
				n.movieId = id
			}
			posterParent.Dom().Clear()
			posterParent.Dom().Append(
				s5.IMG(
					s5.HtmlAttrConstant(s5.HEIGHT, "200"),
					s5.HtmlAttrConstant(s5.WIDTH, "auto"),
					s5.HtmlAttrConstant(s5.SRC, img),
				).Build())
			posterTitle.Dom().SetText(detail.Title)

		case err := <-errCh:
			n.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
		}
	}()
}

func (n *newMoviePage) notLoggedIn(int, string) {
	uilib.SetNewPage(shared.URLGen.Home())
}

func (n *newMoviePage) submit() {
	var movie wire.Movie
	movie.ImdbId = n.movieId
	movie.NominatedBy = n.userUdid
	content, errCh := s5.AjaxPost(&movie, shared.URLGen.MovieResource())
	go func() {
		select {
		case <-content:
			uilib.SetNewPage(shared.URLGen.Home())
		case err := <-errCh:
			n.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
		}
	}()
}

func (n *newMoviePage) loggedIn(user *wire.UserRecord) {
	n.userUdid = user.UserUdid
}

func (n *newMoviePage) Start() {
	n.PageWithFeedback = uilib.NewPageWithFeedback()
	uilib.GetSelf(n.loggedIn, n.notLoggedIn)

	submitButton.CssExistenceAttribute(disabled).Attach(
		s5.NewBooleanInverter(n.havePoster))

	imdbId.Dom().On(s5.CHANGE, func(jquery.Event) {
		n.fetchMovie(imdbId.Dom().Val())
	})

	submitButton.Dom().On(s5.CLICK, func(jquery.Event) {
		n.submit()
	})
}

func main() {
	s5.Main(NewNewMoviePage())
}
