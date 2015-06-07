//these are data structures EXCHANGED between client and server and thus
//are compiled on both sides of the wire.  thus must match exactly to the
//tables created in migrate/migrate.go
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
	Id          int64 `qbs:"pk"`
	ImdbId      string
	NominatedBy string `qbs:"fk:Nominator"`
	Nominator   *UserRecord
}

type Love struct {
	Id       int64 `qbs:"pk"`
	MovieId  int64 `qbs:"fk:Movie"`
	UserUdid string
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

type IMDBDetail struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Metascore  string
	ImdbRating string
	ImdbVotes  string
	ImdbID     string
	Type       string
	Response   string
}
