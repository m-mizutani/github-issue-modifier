package api_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/k0kubun/pp"
	"github.com/m-mizutani/github-issue-modifier/api"
	"github.com/stretchr/testify/require"
)

type config struct {
	SecretArn      string
	GithubEndpoint string
}

func loadConfig() config {
	path := "test.json"
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var conf config
	if err = json.Unmarshal(raw, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

func Test(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       issueCreatedBody,
		Path:       "/",
	}

	conf := loadConfig()

	args := api.Arguments{
		SecretArn:      conf.SecretArn,
		GithubEndpoint: conf.GithubEndpoint,
	}

	resp, err := api.Handler(request, args)
	require.NoError(t, err)
	pp.Println(resp)
}

const issueCreatedBody = `{
	"action": "opened",
	"issue": {
	  "url": "https://api.github.com/repos/m-mizutani/github-issue-modifier/issues/1",
	  "html_url": "https://github.com/m-mizutani/github-issue-modifier/issues/1",
	  "id": 448421921,
	  "number": 1,
	  "title": "test issue",
	  "user": {
		"login": "m-mizutani",
		"id": 605953,
		"avatar_url": "https://avatars0.githubusercontent.com/u/605953?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/m-mizutani",
		"html_url": "https://github.com/m-mizutani",
		"type": "User",
		"site_admin": false
	  },
	  "labels": [
	  ],
	  "state": "open",
	  "locked": false,
	  "assignee": null,
	  "assignees": [
	  ],
	  "milestone": null,
	  "comments": 0,
	  "created_at": "2019-05-25T05:31:07Z",
	  "updated_at": "2019-05-25T05:31:07Z",
	  "closed_at": null,
	  "author_association": "OWNER",
	  "body": "this is a test"
	},
	"repository": {
	  "id": 188520401,
	  "name": "github-issue-modifier",
	  "full_name": "m-mizutani/github-issue-modifier",
	  "private": false,
	  "owner": {
		"login": "m-mizutani",
		"id": 605953,
		"avatar_url": "https://avatars0.githubusercontent.com/u/605953?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/m-mizutani",
		"html_url": "https://github.com/m-mizutani",
		"type": "User",
		"site_admin": false
	  },
	  "html_url": "https://github.com/m-mizutani/github-issue-modifier",
	  "description": null,
	  "fork": false,
	  "url": "https://api.github.com/repos/m-mizutani/github-issue-modifier",
	  "created_at": "2019-05-25T04:27:23Z",
	  "updated_at": "2019-05-25T04:28:14Z",
	  "pushed_at": "2019-05-25T04:28:13Z",
	  "git_url": "git://github.com/m-mizutani/github-issue-modifier.git",
	  "ssh_url": "git@github.com:m-mizutani/github-issue-modifier.git",
	  "clone_url": "https://github.com/m-mizutani/github-issue-modifier.git",
	  "svn_url": "https://github.com/m-mizutani/github-issue-modifier",
	  "homepage": null,
	  "size": 0,
	  "stargazers_count": 0,
	  "watchers_count": 0,
	  "language": null,
	  "has_issues": true,
	  "has_projects": true,
	  "has_downloads": true,
	  "has_wiki": true,
	  "has_pages": false,
	  "forks_count": 0,
	  "mirror_url": null,
	  "archived": false,
	  "disabled": false,
	  "open_issues_count": 1,
	  "license": null,
	  "forks": 0,
	  "open_issues": 1,
	  "watchers": 0,
	  "default_branch": "master"
	},
	"sender": {
	  "login": "m-mizutani",
	  "id": 605953,
	  "avatar_url": "https://avatars0.githubusercontent.com/u/605953?v=4",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/m-mizutani",
	  "html_url": "https://github.com/m-mizutani",
	  "type": "User",
	  "site_admin": false
	}
  }
  `
