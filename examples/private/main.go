package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"

	"github.com/m-mizutani/modifyissue"
)

func addIssueToProjectColumn(client *github.Client, event *github.IssuesEvent) error {
	if event.Action != nil && *event.Action != "opened" {
		return nil
	}

	ColumnID := int64(5462673)

	ctx := context.Background()
	_, _, err := client.Projects.CreateProjectCard(ctx, ColumnID, &github.ProjectCardOptions{
		ContentID:   *event.Issue.ID,
		ContentType: "Issue",
	})

	if err != nil {
		return errors.Wrap(err, "Fail to create project card")
	}

	return nil
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		args := modifyissue.Arguments{
			SecretArn:          os.Getenv("SecretArn"),
			IssueEventCallback: addIssueToProjectColumn,
		}
		return modifyissue.Handler(request, args)
	})
}
