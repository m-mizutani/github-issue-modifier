package api

import (
	"github.com/aws/aws-lambda-go/events"
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
}

// Handler is a main procedue of the API.
func Handler(request events.APIGatewayProxyRequest, args Arguments) (events.APIGatewayProxyResponse, error) {
	Logger.WithField("request", request).Info("Start Handler")

	return events.APIGatewayProxyResponse{
		Body:       `{"message":"ok"}`,
		StatusCode: 200,
	}, nil
}
