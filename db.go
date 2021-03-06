package main

import (
    "os"
     "strconv"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("TABLE_REGION")))

func getItems() (*[]Movie, error) {
    input := &dynamodb.ScanInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
    }

    result, err := db.Scan(input)
    if err != nil {
        return nil, err
    }

    items := new([]Movie)
    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
    return items, nil
}

func getItem(id string) (*Movie, error) {

    input := &dynamodb.GetItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Key: map[string]*dynamodb.AttributeValue{
            "MovieID": {
                S: aws.String(id),
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

    movieResult := new(Movie)
    err = dynamodbattribute.UnmarshalMap(result.Item, movieResult)
    if err != nil {
        return nil, err
    }

    return movieResult, nil
}

func putItem(movie *Movie) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Item: map[string]*dynamodb.AttributeValue{
            "MovieID": {
                S: aws.String(movie.MovieID),
            },
            "Title": {
                S: aws.String(movie.Title),
            },
            "Released": {
                N: aws.String(strconv.Itoa(movie.Released)),
            },
        },
    }

    _, err := db.PutItem(input)
    return err
}
func updateItem(movie *Movie) error {
    input := &dynamodb.UpdateItemInput{
        TableName: aws.String(os.Getenv("TABLE_NAME")),
        Key: map[string]*dynamodb.AttributeValue{
            "MovieID": {
                S: aws.String(movie.MovieID),
            },
        },
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":title": {
                S: aws.String(movie.Title),
            },
            ":released": {
                N: aws.String(strconv.Itoa(movie.Released)),
            },
        },
        ReturnValues:     aws.String("UPDATED_NEW"),
        //UpdateExpression: aws.String("set Title = :title"),
        UpdateExpression: aws.String("set Title = :title, Released = :released"),
    }

    _, err := db.UpdateItem(input)
    return err
}
