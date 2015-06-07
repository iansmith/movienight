package resource

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/coocood/qbs"
	s5 "github.com/seven5/seven5"

	"github.com/iansmith/movienight/wire"
)

type MovieResource struct {
	//stateless
}

var (
	NYI = errors.New("not yet implemented")
)

func (self *MovieResource) IndexQbs(pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	var movies []*wire.Movie
	if err := q.FindAll(&movies); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, "reading db:"+err.Error())
	}
	for _, m := range movies {
		m.Nominator.Admin, m.Nominator.Disabled, m.Nominator.Manager = false, false, false
		m.Nominator.Password = ""
	}
	return movies, nil
}

func (self *MovieResource) PostQbs(i interface{}, pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	movie := i.(*wire.Movie)
	resp, err := http.Get(fmt.Sprintf("http://www.omdbapi.com?i=%s&r=json", movie.ImdbId))
	if err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, "unable to reach omdb")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("omdb error: %d", resp.StatusCode))
	}
	var detail wire.IMDBDetail
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&detail); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, fmt.Sprintf("unable to decode body: %v", err))
	}
	if detail.Response != "True" {
		return nil, s5.HTTPError(http.StatusBadRequest, "cant find "+movie.ImdbId)
	}
	session := pb.Session().UserData().(*wire.SessionData)
	if movie.NominatedBy != session.LoggedInUser.UserUdid {
		return nil, s5.HTTPError(http.StatusBadRequest, "wrong user nominating "+movie.NominatedBy)
	}
	movie.Nominator = nil //just to be safe
	err = q.WhereEqual("imdb_id", movie.ImdbId).Find(movie)
	if err != nil && err != sql.ErrNoRows {
		return nil, s5.HTTPError(http.StatusInternalServerError, "failed on reading db:"+err.Error())
	}
	if err == nil {
		return nil, s5.HTTPError(http.StatusBadRequest, "already exists:"+movie.ImdbId)
	}
	if _, err = q.Save(movie); err != nil {
		return nil, s5.HTTPError(http.StatusInternalServerError, "failed on writing to db:"+err.Error())
	}
	return movie, nil
}
