package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	lambdaStarter "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/sergioyepes21/go-aws-lambda-layer-latest-macro/macromodels"
)

var (
	LAMBDA_LAYER_REGEX = `(^arn:(?:aws[a-zA-Z-]*)?:lambda:[a-z]{2}(?:(?:-gov)|(?:-is'
'o(?:b?)))?-[a-z]+-\d{1}:\d{12}:layer:([a-zA-Z0-9-_]+))`
	ErrInvalidEvent              = errors.New("invalid event")
	ErrLayerNameNotProvided      = errors.New("LayerName param was not provided")
	ErrLambdaLayerRegexFailed    = errors.New("error checking the Lambda Layer Regex")
	ErrLambdaLayerRegexIncorrect = errors.New("incorrect Lambda Layer Name")
	StatusSuccess                = "success"
	StatusFailure                = "failure"
)

func mapEventToStruct(event macromodels.MacroEventMap) (macromodels.Event, error) {
	eventStr, err := json.Marshal(event)
	if err != nil {
		return macromodels.Event{}, ErrInvalidEvent
	}

	macroEventS := macromodels.Event{}
	err = json.Unmarshal(eventStr, &macroEventS)
	if err != nil {
		return macromodels.Event{}, ErrInvalidEvent
	}
	return macroEventS, nil
}

func handler(ctx context.Context, event macromodels.MacroEventMap) (macromodels.ResponseMap, error) {
	macroEventStruct, err := mapEventToStruct(event)
	response := macromodels.ResponseMap{
		"requestId": macroEventStruct.RequestId,
		"status":    StatusFailure,
		"fragment":  "",
	}
	if err != nil {
		return response, err
	}

	layerName := macroEventStruct.Params.LayerName

	if layerName == "" {
		return response, ErrLayerNameNotProvided
	}
	matched, err := regexp.Match(LAMBDA_LAYER_REGEX, []byte(layerName))
	if err != nil {
		return response, ErrLambdaLayerRegexFailed
	}
	if !matched {
		return response, ErrLambdaLayerRegexIncorrect
	}
	awsRegion := os.Getenv("AWS_REGION")
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(awsRegion)}))
	client := lambda.New(sess)
	input := &lambda.ListLayerVersionsInput{
		LayerName: aws.String(layerName),
		MaxItems:  aws.Int64(1),
	}
	result, err := client.ListLayerVersions(input)

	if err != nil {
		return response, fmt.Errorf("error retrieving the layer version of %s", layerName)
	}
	layerVersions := result.LayerVersions

	if len(layerVersions) == 0 {
		return response, fmt.Errorf("no versions found for %s", layerName)
	}
	latestV := *layerVersions[len(layerVersions)-1].LayerVersionArn
	response["fragment"] = latestV
	response["status"] = StatusSuccess
	return response, nil
}

func main() {
	lambdaStarter.Start(handler)
}
