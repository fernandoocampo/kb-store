package dynamodb

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
)

const kbsTable = "kbs"

var (
	updateKBExpression = aws.String("set user_id = :userid, username = :username, event_id = :eventid, content = :content, update_date = :updatedate")
)

var (
	errLoadingAWSConfig = errors.New("unable to load aws config")
	errCreatingDynamodb = errors.New("unable to connect to DynamoDB")
	errSavingKB         = errors.New("unable to save kb")
	errUpdatingKB       = errors.New("unable to update kb")
	errDeletingKB       = errors.New("unable to delete kb")
	errGettingKB        = errors.New("unable to get kb")
	errBuildingKBKey    = errors.New("unable to build kb key")
)

// Setup contains dynamodb settings.
type Setup struct {
	Logger   *slog.Logger
	Region   string
	Endpoint string
}

// Client defines logic for dynamodb repository.
type Client struct {
	client *dynamodb.Client
	logger *slog.Logger
}

func NewClient(ctx context.Context, setup Setup) (*Client, error) {
	newDynamodb := new(Client)
	newDynamodb.logger = setup.Logger

	awsconfig, err := newDynamodb.getConfig(ctx, setup.Region, setup.Endpoint)
	if err != nil {
		return nil, errCreatingDynamodb
	}

	newDynamodb.client = dynamodb.NewFromConfig(awsconfig)

	return newDynamodb, nil
}

func (c *Client) getConfig(ctx context.Context, region, endpoint string) (aws.Config, error) {
	cfg, err := c.loadAWSConfig(ctx, region, endpoint)
	if err != nil {
		c.logger.Error("unable to load aws config", "error", err)

		return cfg, errLoadingAWSConfig
	}

	return cfg, nil
}

func (c *Client) loadAWSConfig(ctx context.Context, region, endpoint string) (aws.Config, error) {
	optFns := make([]func(*config.LoadOptions) error, 0)

	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	if endpoint != "" {
		customResolver := newEndpointResolver(endpoint)
		optFns = append(optFns, config.WithEndpointResolverWithOptions(customResolver))
		optFns = append(optFns, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("d", "d", "")))
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return awsConfig, errLoadingAWSConfig
	}

	return awsConfig, nil
}

func newEndpointResolver(endpoint string) aws.EndpointResolverWithOptionsFunc {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}

		return aws.Endpoint{}, nil
	})
}

func (c *Client) QueryByID(ctx context.Context, kbID kbs.KBID) (*kbs.KB, error) {
	kbKey, err := c.buildTableKey("id", kbID.String())
	if err != nil {
		return nil, errGettingKB
	}

	data, err := c.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(kbsTable),
		Key:       kbKey,
	})
	if err != nil {
		c.logger.Error("unable to get kb", "error", err)

		return nil, errGettingKB
	}

	if data.Item == nil {
		return nil, nil
	}

	var item KB

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		c.logger.Error("unable to unmarshal kb", "error", err)

		return nil, errGettingKB
	}

	kb := item.toRepositoryKB()

	return &kb, nil
}

func (c *Client) Save(ctx context.Context, newKB kbs.KB) error {
	akb := transformKB(newKB)

	data, err := attributevalue.MarshalMap(akb)
	if err != nil {
		c.logger.Error("unable to marshal new kb", "error", err)

		return errSavingKB
	}

	_, err = c.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(kbsTable),
		Item:      data,
	})
	if err != nil {
		c.logger.Error("unable to persist kb", "error", err)

		return errSavingKB
	}

	return nil
}

func (c *Client) Update(ctx context.Context, kb kbs.UpdateKB) error {
	kbKey, err := c.buildTableKey("id", kb.ID.String())
	if err != nil {
		return errDeletingKB
	}

	_, err = c.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:        aws.String(kbsTable),
		Key:              kbKey,
		UpdateExpression: updateKBExpression,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userid":     &types.AttributeValueMemberS{Value: kb.UserID.String()},
			":username":   &types.AttributeValueMemberS{Value: kb.UserName},
			":eventid":    &types.AttributeValueMemberS{Value: kb.EventID.String()},
			":content":    &types.AttributeValueMemberS{Value: kb.Content},
			":updatedate": &types.AttributeValueMemberN{Value: kb.UpdateDateString()},
		},
	})
	if err != nil {
		c.logger.Error("unable to update kb",
			slog.String("id", kb.ID.String()),
			"error", err)

		return errUpdatingKB
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, kb kbs.KB) error {
	kbKey, err := c.buildTableKey("id", kb.ID.String())
	if err != nil {
		return errDeletingKB
	}

	_, err = c.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(kbsTable),
		Key:       kbKey,
	})
	if err != nil {
		c.logger.Error("unable to delete kb from store", "error", err)

		return errDeletingKB
	}

	return nil
}

// https://stackoverflow.com/questions/70019358/how-do-i-get-pagination-working-with-exclusivestartkey-for-dynamodb-aws-sdk-go-v
// https://github.com/aws/aws-sdk-go-v2/issues/1724
// https://docs.aws.amazon.com/code-library/latest/ug/go_2_dynamodb_code_examples.html
func (c *Client) Query(ctx context.Context, filter kbs.QueryFilter) (kbs.SearchKBsResult, error) {
	var result kbs.SearchKBsResult

	keyEx := expression.Key("event_id").Equal(expression.Value(filter.EventID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return result, errGettingKB
	}

	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(kbsTable),
		Limit:                     aws.Int32(int32(filter.RowsPerPage)),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	data, err := c.client.Query(ctx, &queryInput)
	if err != nil {
		c.logger.Error("unable to get kb", "error", err)

		return result, errGettingKB
	}

	if len(data.Items) == 0 {
		return result, nil
	}

	items := make([]KB, len(data.Items))

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		c.logger.Error("unable to unmarshal kbs", "error", err)

		return result, errGettingKB
	}

	result.KBs = make([]kbs.KB, len(items))

	for indx, item := range items {
		result.KBs[indx] = item.toRepositoryKB()
	}

	return result, nil
}

func (c *Client) DatasetStatus(ctx context.Context) error {
	return nil
}

func (c *Client) Count() int {
	return 1
}

func (c *Client) buildTableKey(fieldKey, value string) (map[string]types.AttributeValue, error) {
	selectedKeys := map[string]string{
		fieldKey: value,
	}

	key, err := attributevalue.MarshalMap(selectedKeys)
	if err != nil {
		c.logger.Error("unable to marshal kb keys", "error", err)

		return nil, errBuildingKBKey
	}

	return key, nil
}
