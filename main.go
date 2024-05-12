package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/izaakdale/lib/listener"
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	AirbyteURL          string `envconfig:"AIRBYTE_URL"`
	AirbyteConnectionID string `envconfig:"AIRBYTE_CONNECTION_ID"`
	AirbyteAuth         string `envconfig:"AIRBYTE_AUTH"`
	SQSURL              string `envconfig:"SQS_URL"`
	CredentialsFile     string `envconfig:"CREDENTIALS_FILE"`
}

func main() {
	var spec specification
	err := envconfig.Process("", &spec)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{spec.CredentialsFile},
		),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	err = listener.Initialise(cfg, spec.SQSURL)
	if err != nil {
		log.Fatal(err)
	}
	errCh := make(chan error)

	go listener.Listen(
		context.Background(),
		customersUpdated(
			spec.AirbyteURL,
			spec.AirbyteConnectionID,
			"Basic "+spec.AirbyteAuth,
		), errCh)

	log.Printf("blocking waiting for errors\n")
	err = <-errCh
	if err != nil {
		panic(err)
	}
}
