package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/m-mizutani/github-issue-modifier/api"
)

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		args := api.Arguments{
			SecretArn:      os.Getenv("SecretArn"),
			GithubEndpoint: os.Getenv("GithubEndpoint"),
			Action:         os.Getenv("Action"),
			Region:         os.Getenv("AWS_REGION"),
		}
		return api.Handler(request, args)
	})
}
