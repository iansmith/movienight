package wire

import (
	"time"
)

type UserRecord struct {
	UserUdid                 string `qbs:"pk"`
	FirstName                string
	LastName                 string
	EmailAddr                string
	Password                 string
	Disabled, Admin, Manager bool
}

type Movie struct {
	Id        int64 `qbs:"pk"`
	Title     string
	ImdbId    string
	PosterUrl string
}

type Love struct {
	Id       int64  `qbs:"pk"`
	MovieId  int64  `qbs:"fk:Movie"`
	UserUdid string `qbs:"fk:User"`
	Movie    *Movie
	User     *UserRecord
}

type Hate struct {
	Id       int64  `qbs:"pk"`
	MovieId  int64  `qbs:"fk:Movie"`
	UserUdid string `qbs:"fk:User"`
	Movie    *Movie
	User     *UserRecord
}

type Comment struct {
	Id       int64  `qbs:"pk"`
	MovieId  int64  `qbs:"fk:Movie"`
	UserUdid string `qbs:"fk:User"`
	Movie    *Movie
	User     *UserRecord
	Comment  string
	Updated  time.Time
}

type SessionData struct {
	LoggedInUser *UserRecord
}
