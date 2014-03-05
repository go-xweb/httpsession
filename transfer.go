package httpsession

import (
	"net/http"
	"net/url"
	"time"
)

// Transfer provide and set sessionid
type Transfer interface {
	Get(req *http.Request) (Id, error)
	Set(rw http.ResponseWriter, id Id)
	Clear(rw http.ResponseWriter)
}

// CookieRetriever provide sessionid from cookie
type CookieTransfer struct {
	name   string
	expire time.Duration
}

func NewCookieTransfer(name string) *CookieTransfer {
	return &CookieTransfer{name: name}
}

func (transfer *CookieTransfer) Get(req *http.Request) (Id, error) {
	cookie, err := req.Cookie(transfer.name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	return Id(cookie.Value), nil
}

func (transfer *CookieTransfer) Set(rw http.ResponseWriter, id Id) {
	cookie := http.Cookie{
		Name:     transfer.name,
		Path:     "/",
		Value:    url.QueryEscape(string(id)),
		HttpOnly: true,
		Expires:  time.Now().Add(transfer.expire),
		MaxAge:   int(transfer.expire),
		Secure:   true,
	}

	http.SetCookie(rw, &cookie)
}

func (transfer *CookieTransfer) Clear(rw http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     transfer.name,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
	http.SetCookie(rw, &cookie)
}

var _ Transfer = NewCookieTransfer("test")

// CookieRetriever provide sessionid from url
/*type UrlTransfer struct {
}

func NewUrlTransfer() *UrlTransfer {
	return &UrlTransfer{}
}

func (transfer *UrlTransfer) Get(req *http.Request) (string, error) {
	return "", nil
}

func (transfer *UrlTransfer) Set(rw http.ResponseWriter, id Id) {

}

var (
	_ Transfer = NewUrlTransfer()
)
*/
