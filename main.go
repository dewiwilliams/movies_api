package main

import (
    "encoding/json"
    "log"
    "strconv"
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

type Movie struct {
    movieid     int `json:"movie"`
    Title       string `json:"title"`
    Year        int `json:"year"`
}

func showOne(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    id, err := strconv.Atoi(req.PathParameters["id"]);
    if err != nil {
        return serverError(err)
    }
    movie, err := getItem(id)
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

func showAll(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    movies, err := getItems()
    if err != nil {
        return serverError(err)
    }

    js, err := json.Marshal(movies)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    /*js, _ := json.Marshal(req)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil*/

    if req.HTTPMethod == "GET" && req.Resource=="/movies" {
        return showAll(req)
    }
    if req.HTTPMethod == "GET" && req.Resource=="/movies/{id}" {
        return showOne(req)
    }
    return clientError(http.StatusMethodNotAllowed)
}

func main() {
    lambda.Start(router)
}

