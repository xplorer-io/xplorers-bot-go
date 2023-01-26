package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xplorer-io/xplorers-bot-go/internal/slack"
)

func main() {
	lambda.Start(slack.ProcessRequest)
}
