package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    guuid "github.com/google/uuid"
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
    MovieID     string `json:"movieid"`
    Title       string `json:"title"`
    Released    int `json:"released"`
}

func showOne(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    movie, err := getItem(req.PathParameters["id"])
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

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }

    movie := new(Movie)
    err := json.Unmarshal([]byte(req.Body), movie)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }
    if (movie.Title == "" || movie.Released == 0) {
        return clientError(http.StatusBadRequest)
    }
    movie.MovieID = guuid.New().String()

    err = putItem(movie)
    if err != nil {
        return serverError(err)
    }

    js, err := json.Marshal(movie)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Body:       string(js),
    }, nil
}
func update(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }
    existingMovie, err := getItem(req.PathParameters["id"])
    if err != nil {
        return serverError(err)
    }
    if existingMovie == nil {
        return clientError(http.StatusNotFound)
    }

    newMovieAttributes := new(Movie)
    err = json.Unmarshal([]byte(req.Body), newMovieAttributes)
    if err != nil {
        return clientError(http.StatusUnprocessableEntity)
    }
    if (newMovieAttributes.Title == "" || newMovieAttributes.Released == 0) {
        return clientError(http.StatusBadRequest)
    }

    existingMovie.Title = newMovieAttributes.Title
    existingMovie.Released = newMovieAttributes.Released

    err = updateItem(existingMovie)
    if err != nil {
        return serverError(err)
    }

    js, err := json.Marshal(existingMovie)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 202,
        Body:       string(js),
    }, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.HTTPMethod == "GET" && req.Resource=="/movies" {
        return showAll(req)
    }
    if req.HTTPMethod == "POST" && req.Resource=="/movies" {
        return create(req)
    }
    if req.HTTPMethod == "GET" && req.Resource=="/movies/{id}" {
        return showOne(req)
    }
    if req.HTTPMethod == "PUT" && req.Resource=="/movies/{id}" {
        return update(req)
    }
    return clientError(http.StatusMethodNotAllowed)
}

func main() {
    lambda.Start(router)
}

