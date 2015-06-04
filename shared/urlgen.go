package shared

type URLGenerator interface {
	Me() string
	Auth() string
}

type urlGenImpl struct {
}

var URLGen URLGenerator

func init() {
	URLGen = &urlGenImpl{}
}

func (u *urlGenImpl) Me() string {
	return "/me"
}

func (u *urlGenImpl) Auth() string {
	return "/auth"
}