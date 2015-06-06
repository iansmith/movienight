package resource

import (
	"errors"

	"github.com/coocood/qbs"
	s5 "github.com/seven5/seven5"

	//"github.com/iansmith/movienight/wire"
)

type MovieResource struct {
	//stateless
}

var (
	NYI = errors.New("not yet implemented")
)

func (self *MovieResource) IndexQbs(pb s5.PBundle, q *qbs.Qbs) (interface{}, error) {
	return nil, NYI
}
