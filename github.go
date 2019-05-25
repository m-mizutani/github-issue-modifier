package modifyissue

import (
	"context"
	"net/url"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)


func newGithubClient(endpoint, token string) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid URL format: %s", endpoint)
	}

	client.BaseURL = url
	return client, nil
}
