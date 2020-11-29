# Dewi's movie database lambda handler

This is the code repository for the lambda handler portion of Dewi's movie database.

## Requirements
To compile and run this project you will need,
* Go compiler (I used go1.15.5, but I'm sure other versions will work as well)

Install required dependencies:
* `go get github.com/aws/aws-lambda-go/events`
* `go get github.com/aws/aws-lambda-go/lambda`
* `go get github.com/google/uuid`
* `go get github.com/aws/aws-sdk-go/aws`
* `go get github.com/aws/aws-sdk-go/aws/session`
* `go get github.com/aws/aws-sdk-go/service/dynamodb`
* `go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute`

## Compiling and bundling
To compile and bundle this project, open a command line sessions, and navigate go the directory containing the project. Run:
```
env GOOS=linux GOARCH=amd64 go build -o /tmp/main .
zip -j ./main.zip /tmp/main
```
The generated `main.zip` output can then be provided deployed to an AWS lambda handler.