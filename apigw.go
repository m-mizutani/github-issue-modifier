package modifyissue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

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
