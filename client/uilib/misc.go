// +build js

package uilib

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"
)

func (self *PageWithFeedback) HideFeedback() {
	Hide(self.feedbackrow)
}

func SetNewPage(url string) {
	js.Global.Get("document").Get("location").Set("href", url)
}
func CurrentURLPath() string {
	return js.Global.Get("document").Get("location").Get("pathname").String()
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
