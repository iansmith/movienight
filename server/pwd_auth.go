package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/coocood/qbs"
	s5 "github.com/seven5/seven5"

	"github.com/iansmith/movienight/wire"
)

type movienightValidatingSessionManager struct {
	*s5.SimpleSessionManager
}

// newMovienightValidatingSessionManager creates a new movienightValidatingSessionManager
// but returns it as an s5.ValidatingSessionManager to insure that we meet the
// interface it requires.
func newMovienightValidatingSessionManager() s5.ValidatingSessionManager {
	result := &movienightValidatingSessionManager{}
	result.SimpleSessionManager = s5.NewSimpleSessionManager(result)
	return result
}

//
// Generate converts a unique id, previosuly generated by this session manager
// into a user data record.
//
func (self *movienightValidatingSessionManager) Generate(uniqId string) (interface{}, error) {
	q, err := qbs.GetQbs()
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var ur wire.UserRecord
	ur.UserUdid = uniqId
	if err := q.Find(&ur); err != nil {
		return nil, err
	}
	return &wire.SessionData{&ur}, nil
}

// Check that a username and password are as we have them in the database.
// If they match, we return the user's UDID as the uniq value for the session
// plus the user data record.
func (self *movienightValidatingSessionManager) ValidateCredentials(username, pwd string) (string, interface{}, error) {
	q, err := qbs.GetQbs()
	if err != nil {
		return "", nil, err
	}
	defer q.Close()

	var ur wire.UserRecord
	u := strings.TrimSpace(username)
	p := strings.TrimSpace(pwd)
	if len(u) == 0 || len(p) == 0 {
		return "", nil, err
	}
	cond := qbs.NewEqualCondition("email_addr", u).AndEqual("password", p)
	if err := q.Condition(cond).Find(&ur); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error trying to validate credentials: %v", err)
			return "", nil, err
		}
		//normal case of bad pwd
		return "", nil, nil
	}
	//return the udid as uniq part,then the rest of the object as user data
	return ur.UserUdid, &wire.SessionData{&ur}, nil
}

// SendUserDetails is responsible for filtering out fields that we may not wish
// to send to the client side of the wire when returning a user record.
func (self *movienightValidatingSessionManager) SendUserDetails(i interface{}, w http.ResponseWriter) error {
	ur := i.(*wire.SessionData).LoggedInUser
	ur.Password = ""
	return s5.SendJson(w, ur)
}

func (self *movienightValidatingSessionManager) UseResetRequest(userId string, requestId string, newpwd string) (bool, error) {
	panic("UseResetRequest not implemented")
}

func (self *movienightValidatingSessionManager) GenerateResetRequest(userId string) (string, error) {
	panic("GenerateResetRequest not implemented")
}
