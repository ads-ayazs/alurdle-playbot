package store

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (sm *oneGameSM) Save(in interface{}) error {
	svc, err := getAwsDynamoService()
	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(in)
	if err != nil {
		log.Error(ErrOneSmMarshalFailure, err)
		return err
	}
	log.Infof("marshalled struct: %+v", av)

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(CONST_AWS_DYNAMODB_TABLENAME),
	})
	if err != nil {
		return err
	}

	return nil
}

//////////////////

var (
	CONST_AWS_DYNAMODB_TABLENAME = "PlayerGame"
)

type oneGameSM struct{}

func createOneGameSM() (StoreManager, error) {
	sm := new(oneGameSM)
	return sm, nil
}

func (sm *oneGameSM) initStore() error {
	svc, err := getAwsDynamoService()
	if err != nil {
		return err
	}

	ctx := context.TODO()

	out, err := svc.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("playerName"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("gameId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("playerName"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("gameId"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(CONST_AWS_DYNAMODB_TABLENAME),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(out)

	w := dynamodb.NewTableExistsWaiter(svc)
	err = w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(CONST_AWS_DYNAMODB_TABLENAME),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return err
	}

	return nil
}

var (
	TEST_AWS_DYNAMODB_TABLENAME = "TEST_" + CONST_AWS_DYNAMODB_TABLENAME
)

func useTestMode() error {
	CONST_AWS_DYNAMODB_TABLENAME = TEST_AWS_DYNAMODB_TABLENAME

	return nil
}

func cleanupTestMode() error {
	svc, err := getAwsDynamoService()
	if err != nil {
		return err
	}

	ctx := context.TODO()

	out, err := svc.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(TEST_AWS_DYNAMODB_TABLENAME),
	})
	if err != nil {
		return err
	}

	w := dynamodb.NewTableNotExistsWaiter(svc)
	err = w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(CONST_AWS_DYNAMODB_TABLENAME),
		},
		2*time.Minute,
		func(o *dynamodb.TableNotExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return err
	}

	log.Debug(out)
	return nil
}
