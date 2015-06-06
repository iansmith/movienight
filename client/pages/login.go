package main

import (
	"net/http"
	"strings"

	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"

	"github.com/iansmith/movienight/client/uilib"
	"github.com/iansmith/movienight/shared"
)

var (
	submitButton = s5.NewHtmlId("button", "submit")

	disabled = s5.NewCssClass("disabled")
)

type loginPage struct {
	*uilib.PageWithFeedback
	isValid s5.BooleanAttribute
	email   s5.FormElement
	pwd     s5.FormElement
}

func (l *loginPage) formValid(data map[string]s5.Equaler) s5.Equaler {
	email := strings.TrimSpace(data[l.email.Selector()].(s5.StringEqualer).S)
	pwd := strings.TrimSpace(data[l.pwd.Selector()].(s5.StringEqualer).S)

	return s5.BoolEqualer{email != "" && pwd != ""}
}

func (l *loginPage) attemptLogin() {
	email := l.email.Val()
	pwd := l.pwd.Val()
	auth := &uilib.PasswordAuthParameters{
		Username: email,
		Password: pwd,
		Op:       uilib.AUTH_OP_LOGIN,
	}
	loginContent, loginErrChan := s5.AjaxPost(auth, shared.URLGen.Auth())
	go func() {
		select {
		case <-loginContent:
			uilib.SetNewPage(shared.URLGen.Home())
		case err := <-loginErrChan:
			if err.StatusCode == http.StatusUnauthorized {
				l.PageWithFeedback.DisplayFeedback("That's probably not your password", uilib.Danger)
			} else {
				l.PageWithFeedback.DisplayFeedback(err.Message, uilib.Danger)
			}
		}
	}()
}

func (l *loginPage) Start() {
	l.PageWithFeedback = uilib.NewPageWithFeedback()

	l.isValid.Attach(s5.NewFormValidConstraint(l.formValid, l.email, l.pwd))
	submitButton.CssExistenceAttribute(disabled).Attach(
		s5.NewBooleanInverter(l.isValid))

	submitButton.Dom().On(s5.CLICK, func(e jquery.Event) {
		e.PreventDefault()
		l.attemptLogin()
	})

}

func newLoginPage() *loginPage {
	result := &loginPage{
		isValid: s5.NewBooleanSimple(false),
		email:   s5.NewInputTextId("email"),
		pwd:     s5.NewInputTextId("password"),
		//don't init the PageWithFeedback yet, DOM not loaded
	}

	return result
}

func main() {
	p := newLoginPage()
	s5.Main(p)
}
