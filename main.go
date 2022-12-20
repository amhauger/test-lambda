package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	shutdownFunc, err := CreateTracerAndBatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("error instantiating otel monitoring")
	}
	defer shutdownFunc()

	lambda.Start(otellambda.InstrumentHandler(handleLambda))
}

func handleLambda(ctx context.Context, sqsEvent events.SQSEvent) error {
	tracer := otel.Tracer("testLayers.handleLambda")
	_, handleSpan := tracer.Start(
		ctx,
		"app.handleLambda",
		trace.WithAttributes(
			attribute.String("hello", "world"),
		),
	)
	defer handleSpan.End()

	log.Info().Strs("environmentVariables", os.Environ()).Msg("hello from lambda")
	fmt.Println("hello from lambda")
	return nil
}
