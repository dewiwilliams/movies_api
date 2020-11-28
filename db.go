package main

import (
    "os"
     "strconv"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-2"))

func getItem(id int) (*movie, error) {

    input := &dynamodb.GetItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Key: map[string]*dynamodb.AttributeValue{
            "movieid": {
                N: aws.String(strconv.Itoa(id)),
            },
        },
    }

    result, err := db.GetItem(input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, nil
    }

    movieResult := new(movie)
    err = dynamodbattribute.UnmarshalMap(result.Item, movieResult)
    if err != nil {
        return nil, err
    }

    return movieResult, nil
}

func putItem(movie *movie) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Item: map[string]*dynamodb.AttributeValue{
            "movieid": {
                N: aws.String(strconv.Itoa(movie.movieid)),
            },
            "Title": {
                S: aws.String(movie.Title),
            },
            "Year": {
                N: aws.String(strconv.Itoa(movie.Year)),
            },
        },
    }

    _, err := db.PutItem(input)
    return err
}