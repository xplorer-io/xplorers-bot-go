package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/xplorer-io/xplorers-bot-go/internal/chatgpt"
)

func main() {
	lambda.Start(chatgpt.ProcessSlackRequest)
}
