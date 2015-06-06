// +build js

package uilib

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"
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

func (self *PageWithFeedback) HideFeedback() {
	Hide(self.feedbackrow)
}

func SetNewPage(url string) {
	js.Global.Get("document").Get("location").Set("href", url)
}

//
// UTILITY JQUERY ANIMATIONS
//
func SlideDown(id s5.HtmlId) {
	selector := id.TagName() + "#" + id.Id()
	if !s5.TestMode {
		jq := jquery.NewJQuery(selector)
		jq.Underlying().Call("slideDown")
	}
}
func SlideUp(id s5.HtmlId) {
	selector := id.TagName() + "#" + id.Id()
	if !s5.TestMode {
		jq := jquery.NewJQuery(selector)
		jq.Underlying().Call("slideUp")
	}
}
func Hide(id s5.HtmlId) {
	//xxx should be in the toolkit because tests should be able to model the
	//xxx visiblity state
	selector := id.TagName() + "#" + id.Id()
	if !s5.TestMode {
		jq := jquery.NewJQuery(selector)
		jq.Hide()
	}
}
