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
	ProxyPoster(id string) string
	ProxyMovie(id string, fullPlot bool) string
	DetailPage(id string) string
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
func (u *urlGenImpl) DetailPage(id string) string {
	return "/movie/" + id
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

func (u *urlGenImpl) ProxyPoster(id string) string {
	return "/posterproxy?i=" + id
}
func (u *urlGenImpl) ProxyMovie(id string, fullPlot bool) string {
	plot := "short"
	if fullPlot {
		plot = "full"
	}
	return "/movieproxy?i=" + id + "&plot=" + plot
}

func (u *urlGenImpl) Auth() string {
	return "/auth"
}
