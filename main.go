package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

type movie struct {
    movieid     int `json:"movie"`
    Title       string `json:"title"`
    Year        int `json:"year"`
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    movie, err := getItem(100)
    if err != nil {
        return serverError(err)
    }
    if movie == nil {
        return clientError(http.StatusNotFound)
    }

    js, err := json.Marshal(movie)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func main() {
    lambda.Start(show)
}

