package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

func main() {
	lambda.Start(otellambda.InstrumentHandler(handleLambda))
}

func handleLambda(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Info().Msg("hello from lambda")
	fmt.Println("hello from lambda")
	os.Environ()
	return nil
}
