// +build js

package uilib

import (
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/shared"
	"github.com/iansmith/movienight/wire"
)

var (
	bgSuccess = s5.NewCssClass("bg-success")
	bgDanger  = s5.NewCssClass("bg-danger")
	bgWarning = s5.NewCssClass("bg-warning")
)

type FeedbackType int

const (
	Success FeedbackType = iota
	Warning
	Danger
)

type PageWithFeedback struct {
	feedbackrow  s5.HtmlId
	feedbackcol  s5.HtmlId
	feedbackText s5.HtmlId
	closer       s5.HtmlId
}

func NewPageWithFeedback() *PageWithFeedback {
	result := &PageWithFeedback{
		feedbackrow:  s5.NewHtmlId("div", "feedback-row"),
		feedbackcol:  s5.NewHtmlId("div", "feedback-col"),
		feedbackText: s5.NewHtmlId("span", "feedback-text"),
		closer:       s5.NewHtmlId("button", "err-closer"),
	}
	result.closer.Dom().On(s5.CLICK, func(evt jquery.Event) {
		SlideUp(result.feedbackrow)
	})

	return result
}

func (self *PageWithFeedback) DisplayFeedback(msg string, ft FeedbackType) {
	self.feedbackcol.Dom().RemoveClass(bgDanger.ClassName())
	self.feedbackcol.Dom().RemoveClass(bgWarning.ClassName())
	self.feedbackcol.Dom().RemoveClass(bgSuccess.ClassName())
	switch ft {
	case Danger:
		self.feedbackcol.Dom().AddClass(bgDanger.ClassName())
	case Warning:
		self.feedbackcol.Dom().AddClass(bgWarning.ClassName())
	case Success:
		self.feedbackcol.Dom().AddClass(bgSuccess.ClassName())
	}
	self.feedbackText.Dom().SetText(msg)
	SlideDown(self.feedbackrow)

}

func GetSelf(successFn func(rec *wire.UserRecord), failFn func(int, string)) {
	content, errorChan := s5.AjaxGet(&wire.UserRecord{}, shared.URLGen.Me())
	go func() {
		select {
		case raw := <-content:
			result := raw.(*wire.UserRecord)
			successFn(result)
		case err := <-errorChan:
			failFn(err.StatusCode, err.Message)
		}
	}()
}
