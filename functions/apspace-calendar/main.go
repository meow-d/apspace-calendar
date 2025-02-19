package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/meow-d/apspace-calendar/src/calendar"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	intake := request.QueryStringParameters["intake"]
	if intake == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing required parameter: intake"}, nil
	}

	titleFormat := request.QueryStringParameters["title"]
	if titleFormat == "" {
		titleFormat = "module_name"
	}

	icsData, err := calendar.FetchAndConvert(intake, titleFormat)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to process calendar"}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":        "text/calendar",
			"Content-Disposition": "attachment; filename=calendar.ics",
		},
		Body: icsData,
	}, nil
}

func main() {
	lambda.Start(handler)
}
