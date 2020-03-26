package storage_test

import (
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/convox/console/pkg/cycle"
	"github.com/convox/console/pkg/storage"
	"github.com/stretchr/testify/require"
)

type Item struct {
	ID string `dynamo:"id"`
}

type Items []Item

func TestDynamoGet(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.GetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"Key":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Response: []byte(`{"Item":{"id":{"S":"foo"}}}`),
	})

	var item Item

	err := d.Get("items", "foo", &item)

	require.NoError(t, err)
	require.NotNil(t, item)
	require.Equal(t, "foo", item.ID)
}

func TestDynamoGetError(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.GetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"Key":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Code:     400,
		Response: []byte(`{"__type":"com.amazonaws.dynamodb.v20120810#ResourceNotFoundException","message":"test"}`),
	})

	var item Item

	err := d.Get("items", "foo", &item)

	require.EqualError(t, err, "ResourceNotFoundException: test\n\tstatus code: 400, request id: ")
}

func TestDynamoGetInvalidType(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.GetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"Key":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Response: []byte(`{"Item":{"id":{"S":"foo"}}}`),
	})

	var item struct {
		ID int `dynamo:"id"`
	}

	err := d.Get("items", "foo", item)

	require.EqualError(t, err, "item must be a pointer")

	err = d.Get("items", "foo", &item)

	require.EqualError(t, err, `"id" is not a number`)
}

func TestDynamoGetNotFound(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.GetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"Key":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Response: []byte(`{"Item":{}}`),
	})

	var item Item

	err := d.Get("items", "foo", &item)

	require.EqualError(t, err, "not found")
	require.IsType(t, storage.NotFoundError{}, err)
}

func TestDynamoGetBatch(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.BatchGetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"RequestItems":{"test-items":{"ConsistentRead":true,"Keys":[{"id":{"S":"foo"}}]}}}`),
		Response: []byte(`{"Responses":{"test-items":[{"id":{"S":"foo"}}]}}`),
	})

	var items Items

	err := d.GetBatch("items", []string{"foo"}, &items)

	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "foo", items[0].ID)
}

func TestDynamoGetBatchInvalidType(t *testing.T) {
	d, ct := mockDynamo()

	var items1 struct{}
	err := d.GetBatch("items", []string{}, &items1)
	require.EqualError(t, err, "items must be a pointer to a slice")

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.BatchGetItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"RequestItems":{"test-items":{"ConsistentRead":true,"Keys":[{"id":{"S":"foo"}}]}}}`),
		Response: []byte(`{"Responses":{"test-items":[{"id":{"S":"foo"}}]}}`),
	})

	var items2 []struct {
		ID int `dynamo:"id"`
	}

	err = d.GetBatch("items", []string{"foo"}, &items2)

	require.EqualError(t, err, `"id" is not a number`)
}

func TestDynamoGetIndex(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Query"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ExpressionAttributeNames":{"#id":"id"},"ExpressionAttributeValues":{":id":{"S":"foo"}},"IndexName":"idx","KeyConditionExpression":"#id=:id","TableName":"test-items"}`),
		Response: []byte(`{"Items":[{"id":{"S":"foo"}}]}`),
	})

	var items Items

	err := d.GetIndex("items", "idx", map[string]string{"id": "foo"}, &items)

	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "foo", items[0].ID)
}

func TestDynamoGetIndexError(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Query"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ExpressionAttributeNames":{"#id":"id"},"ExpressionAttributeValues":{":id":{"S":"foo"}},"IndexName":"idx","KeyConditionExpression":"#id=:id","TableName":"test-items"}`),
		Code:     400,
		Response: []byte(`{"__type":"com.amazonaws.dynamodb.v20120810#ResourceNotFoundException","message":"test"}`),
	})

	var items Items

	err := d.GetIndex("items", "idx", map[string]string{"id": "foo"}, &items)

	require.EqualError(t, err, "ResourceNotFoundException: test\n\tstatus code: 400, request id: ")
}

func TestDynamoGetIndexInvalidType(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Query"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ExpressionAttributeNames":{"#id":"id"},"ExpressionAttributeValues":{":id":{"S":"foo"}},"IndexName":"idx","KeyConditionExpression":"#id=:id","TableName":"test-items"}`),
		Response: []byte(`{"Items":[{"id":{"S":"foo"}}]}`),
	})

	var items []struct {
		ID int `dynamo:"id"`
	}

	err := d.GetIndex("items", "idx", map[string]string{"id": "foo"}, items)

	require.EqualError(t, err, "items must be a pointer to a slice")

	err = d.GetIndex("items", "idx", map[string]string{"id": "foo"}, &items)

	require.EqualError(t, err, `"id" is not a number`)
}

func TestDynamoGetIndexNotFound(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Query"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ExpressionAttributeNames":{"#id":"id"},"ExpressionAttributeValues":{":id":{"S":"foo"}},"IndexName":"idx","KeyConditionExpression":"#id=:id","TableName":"test-items"}`),
		Response: []byte(`{"Items":[]}`),
	})

	var items Items

	err := d.GetIndex("items", "idx", map[string]string{"id": "foo"}, &items)

	require.Nil(t, err)
	require.Len(t, items, 0)
}

func TestDynamoList(t *testing.T) {
	d, ct := mockDynamo()

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Scan"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"TableName":"test-items"}`),
		Response: []byte(`{"Items":[{"id":{"S":"foo"}}]}`),
	})

	var items Items

	err := d.List("items", &items)

	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "foo", items[0].ID)
}

func TestDynamoListInvalidType(t *testing.T) {
	d, ct := mockDynamo()

	var items1 struct{}

	err := d.List("items", &items1)

	require.EqualError(t, err, "items must be a pointer to a slice")

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.Scan"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"ConsistentRead":true,"TableName":"test-items"}`),
		Response: []byte(`{"Items":[{"id":{"S":"foo"}}]}`),
	})

	var items2 []struct {
		ID int `dynamo:"id"`
	}

	err = d.List("items", &items2)

	require.EqualError(t, err, `"id" is not a number`)
}

func TestDynamoPut(t *testing.T) {
	d, ct := mockDynamo()

	item := Item{ID: "foo"}

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.PutItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"Item":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Response: []byte(`{}`),
	})

	err := d.Put("items", &item)

	require.NoError(t, err)
}

func TestDynamoPutError(t *testing.T) {
	d, ct := mockDynamo()

	item := Item{ID: "foo"}

	ct.Register(cycle.Cycle{
		Headers:  map[string]string{"X-Amz-Target": "DynamoDB_20120810.PutItem"},
		Method:   "POST",
		Path:     "/",
		Request:  []byte(`{"Item":{"id":{"S":"foo"}},"TableName":"test-items"}`),
		Code:     400,
		Response: []byte(`{"__type":"com.amazonaws.dynamodb.v20120810#ResourceNotFoundException","message":"test"}`),
	})

	err := d.Put("items", &item)

	require.EqualError(t, err, "ResourceNotFoundException: test\n\tstatus code: 400, request id: ")
}

func TestDynamoPutInvalidType(t *testing.T) {
	d, _ := mockDynamo()

	err := d.Put("items", Items{})

	require.EqualError(t, err, "item must be a pointer")

	item := struct {
		ID chan string `dynamo:"id"`
	}{ID: nil}

	err = d.Put("items", &item)

	require.EqualError(t, err, `could not marshal ID into "id" as json: json: unsupported type: chan string`)
}

func mockDynamo() (*storage.Dynamo, *cycle.Tester) {
	ct := &cycle.Tester{}

	ts := httptest.NewServer(ct)

	d := storage.NewDynamo("test", "key")

	d.DB = dynamodb.New(session.New(), &aws.Config{
		Credentials: credentials.NewStaticCredentials("foo", "bar", "baz"),
		Endpoint:    aws.String(ts.URL),
		Region:      aws.String("us-test-1"),
	})

	return d, ct
}
