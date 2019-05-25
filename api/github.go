package api

import (
	"context"
	"net/url"
	"strconv"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type githubIssue struct {
	URL       string       `json:"url"`
	ID        int64        `json:"id"`
	Assignees []githubUser `json:"assignees"`
}

type githubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Type  string `json:"type"`
}

type githubRepository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type githubWebHookPayload struct {
	Action     string           `json:"action"`
	Issue      githubIssue      `json:"issue"`
	Sender     githubUser       `json:"sender"`
	Repository githubRepository `json:"repository"`
}

type actionArgument struct {
	Action     string            `json:"action"`
	Conditions []actionCondition `json:"conditions"`
	Value      string            `json:"value"`
}

type actionCondition struct {
	Repositories []string `json:"repositories,omitempty"`
	Assignees    []string `json:"assingees,omitempty"`
}

func (x actionCondition) inRepositories(fullName string) bool {
	for _, repo := range x.Repositories {
		if fullName == repo {
			return true
		}
	}

	return false
}

func (x actionCondition) inAssignees(userName string) bool {
	for _, assignee := range x.Assignees {
		if assignee == userName {
			return true
		}
	}

	return false
}

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

func addIssueToProjectColumn(client *github.Client, payload githubWebHookPayload, action actionArgument) error {
	if payload.Action != "opened" {
		return nil
	}

	id, err := strconv.ParseInt(action.Value, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "Invalid parameter (CloumnID must be integer): %s", action.Value)
	}

	ctx := context.Background()
	card, resp, err := client.Projects.CreateProjectCard(ctx, id, &github.ProjectCardOptions{
		ContentID:   payload.Issue.ID,
		ContentType: "Issue",
	})

	Logger.WithFields(logrus.Fields{
		"card":  card,
		"resp":  resp,
		"error": err,
	}).Info("Done github operation")

	if err != nil {
		return errors.Wrap(err, "Fail to create project card")
	}

	return nil
}
