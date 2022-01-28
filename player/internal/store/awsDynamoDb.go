package store

import (
	"context"

	wdlConfig "aluance.io/wordleplayer/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	log "github.com/sirupsen/logrus"
)

func getAwsConfig() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(wdlConfig.CONFIG_AWS_REGION),
	)
	if err != nil {
		log.Error(ErrAwsConfig, err)
		return nil, ErrAwsConfig
	}

	return &cfg, nil
}

func getAwsDynamoService() (*dynamodb.Client, error) {
	cfg, err := getAwsConfig()
	if err != nil {
		return nil, err
	}

	svc := dynamodb.NewFromConfig(*cfg)

	return svc, nil
}
