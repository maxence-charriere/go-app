package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/maxence-charriere/go-app/v10/pkg/errors"
)

func get(ctx app.Context, path string) ([]byte, error) {
	url := path
	if !strings.HasPrefix(url, "http") {
		u := *ctx.Page().URL()
		u.Path = path
		url = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New("creating request failed").
			WithTag("path", path).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("getting document failed").
			WithTag("path", path).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.New(res.Status).WithTag("path", path)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("reading document failed").
			WithTag("path", path).
			Wrap(err)
	}
	return b, nil
}
