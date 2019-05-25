package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logger can be modified from external
var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}

type Arguments struct {
	SecretArn      string
	GithubEndpoint string
	Action         string
	Region         string
}

type secretValues struct {
	GithubToken string `json:"github_token"`
}

func replySuccess() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       `{"message":"ok"}`,
		StatusCode: 200,
	}, nil
}

func replySystemError(err error, f string, v ...interface{}) (events.APIGatewayProxyResponse, error) {
	msg := fmt.Sprintf(f, v...)
	Logger.WithError(err).WithField("message", msg).Error("System Error")
	return events.APIGatewayProxyResponse{
		Body:       `{"message":"system error"}`,
		StatusCode: 500,
	}, errors.Wrap(err, msg)
}

func replyUserError(err error, f string, v ...interface{}) (events.APIGatewayProxyResponse, error) {
	msg := fmt.Sprintf(f, v...)
	Logger.WithError(err).WithField("message", msg).Error("User Error")

	reply := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}
	raw, err := json.Marshal(reply)
	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(raw),
		StatusCode: 400,
	}, errors.Wrap(err, msg)
}

// Handler is a main procedue of the API.
func Handler(request events.APIGatewayProxyRequest, args Arguments) (events.APIGatewayProxyResponse, error) {
	Logger.WithField("request", request).Info("Start Handler")

	var secrets secretValues
	if err := getSecretValues(args.SecretArn, &secrets); err != nil {
		return replySystemError(err, "Fail to get secret values: %s", args.SecretArn)
	}

	var payload githubWebHookPayload
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		return replyUserError(err, "Invalid JSON format payload of request: %s", request.Body)
	}

	client, err := newGithubClient(args.GithubEndpoint, secrets.GithubToken)
	if err != nil {
		return replyUserError(err, "Fail to create github client")
	}

	var action actionArgument
	if err := json.Unmarshal([]byte(args.Action), &action); err != nil {
		return replyUserError(err, "Invalid JSON format payload for action: %s", args.Action)
	}

	if err := addIssueToProjectColumn(client, payload, action); err != nil {
		return replySystemError(err, "Fail to github operation")
	}

	return replySuccess()
}
