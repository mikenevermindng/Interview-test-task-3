package my_http

import (
	"fmt"
	"net/http"
)

type MyHttp struct {
	Client *http.Client
}

func New(maxRedirect int) *MyHttp {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > maxRedirect {
				return fmt.Errorf("stopped after %d redirects", maxRedirect)
			}
			// Allow the redirect
			return nil
		},
	}

	return &MyHttp{
		Client: client,
	}
}
