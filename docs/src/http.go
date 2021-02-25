package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/errors"
)

func get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("creating request failed").
			Tag("path", path).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("getting document failed").
			Tag("path", path).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.New(res.Status).Tag("path", path)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("reading document failed").
			Tag("path", path).
			Wrap(err)
	}
	return b, nil
}
