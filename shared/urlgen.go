//things in this package are compiled *BOTH* for client and server
package shared

//URLGenerator prevents people from using hard-coded strings in URLs or, worse
//fmt.Sprintf("/blah/%d/%s/%d")
type URLGenerator interface {
	Me() string
	Auth() string
	Home() string
	Login() string
	NewMovie() string
	MovieResource() string
}

type urlGenImpl struct {
}

var URLGen URLGenerator

func init() {
	URLGen = &urlGenImpl{}
}

func (u *urlGenImpl) Login() string {
	return "/login.html"
}
func (u *urlGenImpl) MovieResource() string {
	return "/rest/movie"
}

func (u *urlGenImpl) NewMovie() string {
	return "/newmovie.html"
}
func (u *urlGenImpl) Me() string {
	return "/me"
}

func (u *urlGenImpl) Home() string {
	return "/"
}

func (u *urlGenImpl) Auth() string {
	return "/auth"
}
