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
			SecretArn: os.Getenv("SecretArn"),
		}
		return api.Handler(request, args)
	})
}
