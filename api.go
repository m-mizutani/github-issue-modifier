package modifyissue

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v25/github"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
)

// Logger can be modified from external
var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}

type issueEventCallback func(client *github.Client, event *github.IssuesEvent) error

type Arguments struct {
	SecretArn      string
	GithubEndpoint string

	// Callbacks
	IssueEventCallback issueEventCallback
}

type secretValues struct {
	GithubToken string `json:"github_token"`
}

// Handler is a main procedue of the API.
func Handler(request events.APIGatewayProxyRequest, args Arguments) (events.APIGatewayProxyResponse, error) {
	Logger.WithField("request", request).Info("Start Handler")

	var secrets secretValues
	if err := getSecretValues(args.SecretArn, &secrets); err != nil {
		return replySystemError(err, "Fail to get secret values: %s", args.SecretArn)
	}

	githubEvent, ok := request.Headers["X-Github-Event"]
	pp.Println(request)
	if !ok {
		return replyUserError(nil, "Missing X-Github-Event in request header")
	}

	event, err := github.ParseWebHook(githubEvent, []byte(request.Body))
	if err != nil {
		return replyUserError(err, "Fail to parse JSON payload in request: %s", request.Body)
	}

	client, err := newGithubClient(args.GithubEndpoint, secrets.GithubToken)
	if err != nil {
		return replyUserError(err, "Fail to create github client")
	}

	switch event := event.(type) {
	case *github.IssuesEvent:
		if args.IssueEventCallback != nil {
			if err := args.IssueEventCallback(client, event); err != nil {
				return replySystemError(err, "Fail to github operation")
			}
		}
	}

	return replySuccess()
}
