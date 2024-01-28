package uri_builder

import (
	"strings"
)

type UrlBuilder struct {
	ssl bool
	uri string
}

func (u *UrlBuilder) Ssl(ssl bool) {
	u.ssl = ssl
}

func (u *UrlBuilder) Uri(uri string) {
	u.uri = uri
}

func (u *UrlBuilder) GetUri() string {
	return u.uri
}

func (u *UrlBuilder) GetSsl() bool {
	return u.ssl
}

func (u *UrlBuilder) BuildURI() (bool, string) {

	buildUri := u.uri
	// Check if the URL has a scheme
	if !strings.Contains(u.uri, "://") {
		if u.ssl {
			buildUri = "https://" + buildUri
		} else {
			buildUri = "http://" + buildUri
		}
	}

	// Return the rebuilt URL
	return true, buildUri
}
