package storage

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/convox/console/pkg/dynamique"
	"github.com/pkg/errors"
)

type Dynamo struct {
	DB dynamodbiface.DynamoDBAPI

	key    string
	prefix string
}

func NewDynamo(prefix, key string) *Dynamo {
	return &Dynamo{
		DB:     dynamodb.New(session.New()),
		key:    key,
		prefix: prefix,
	}
}

func (d *Dynamo) Delete(table, id string) error {
	_, err := d.DB.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": &dynamodb.AttributeValue{S: aws.String(id)}},
		TableName: aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *Dynamo) Get(table, id string, item interface{}) error {
	log := log.At("dynamo.get").Append("table=%s id=%s", table, id).Start()

	if r := reflect.TypeOf(item); r.Kind() != reflect.Ptr {
		return errors.WithStack(log.Error(fmt.Errorf("item must be a pointer")))
	}

	res, err := d.DB.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		Key:            map[string]*dynamodb.AttributeValue{"id": &dynamodb.AttributeValue{S: aws.String(id)}},
		TableName:      aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
	})
	if err != nil {
		return errors.WithStack(log.Error(err))
	}
	if len(res.Item) == 0 {
		return NotFoundError{fmt.Errorf("not found")}
	}

	v := reflect.New(reflect.TypeOf(item).Elem())

	if err := dynamique.Unmarshal(res.Item, v.Interface(), d.key); err != nil {
		return errors.WithStack(log.Error(err))
	}

	reflect.ValueOf(item).Elem().Set(v.Elem())

	if d, ok := item.(Defaulter); ok {
		d.Defaults()
	}

	return log.Success()
}

func (d *Dynamo) GetBatch(table string, ids []string, items interface{}) error {
	log := log.At("dynamo.getbatch").Append("table=%s ids=%d", table, len(ids)).Start()

	if r := reflect.TypeOf(items); r.Kind() != reflect.Ptr || r.Elem().Kind() != reflect.Slice {
		return errors.WithStack(log.Error(fmt.Errorf("items must be a pointer to a slice")))
	}

	for _, page := range paginate(ids, 100) {
		keys := []map[string]*dynamodb.AttributeValue{}

		for _, id := range page {
			keys = append(keys, map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: aws.String(id)},
			})
		}

		req := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				fmt.Sprintf("%s-%s", d.prefix, table): &dynamodb.KeysAndAttributes{
					ConsistentRead: aws.Bool(true),
					Keys:           keys,
				},
			},
		}

		ctx := context.Background()

		p := request.Pagination{
			NewRequest: func() (*request.Request, error) {
				r, _ := d.DB.BatchGetItemRequest(req)
				r.SetContext(ctx)
				return r, nil
			},
		}

		ritems := reflect.ValueOf(items).Elem()
		kind := reflect.TypeOf(items).Elem().Elem()

		for p.Next() {
			page := p.Page().(*dynamodb.BatchGetItemOutput)

			for _, i := range page.Responses[fmt.Sprintf("%s-%s", d.prefix, table)] {
				v := reflect.New(kind)

				if err := dynamique.Unmarshal(i, v.Interface(), d.key); err != nil {
					return errors.WithStack(log.Error(err))
				}

				if d, ok := v.Interface().(Defaulter); ok {
					d.Defaults()
				}

				ritems.Set(reflect.Append(ritems, v.Elem()))
			}
		}

		if err := p.Err(); err != nil {
			return errors.WithStack(log.Error(err))
		}
	}

	if s, ok := reflect.ValueOf(items).Elem().Interface().(Sortable); ok {
		sort.Slice(s, s.Less)
	}

	return log.Success()
}

func (d *Dynamo) GetIndex(table, index string, key map[string]string, items interface{}) error {
	log := log.At("dynamo.getindex").Append("table=%s index=%s key=%v", table, index, key).Start()

	if r := reflect.TypeOf(items); r.Kind() != reflect.Ptr || r.Elem().Kind() != reflect.Slice {
		return errors.WithStack(log.Error(fmt.Errorf("items must be a pointer to a slice")))
	}

	kc := []string{}
	ea := map[string]*string{}
	ev := map[string]*dynamodb.AttributeValue{}

	for k, v := range key {
		ks := strings.Replace(k, "-", "", -1)
		kc = append(kc, fmt.Sprintf("#%s=:%s", ks, ks))
		ea[fmt.Sprintf("#%s", ks)] = aws.String(k)
		ev[fmt.Sprintf(":%s", ks)] = &dynamodb.AttributeValue{S: aws.String(v)}
	}

	res, err := d.DB.Query(&dynamodb.QueryInput{
		ExpressionAttributeNames:  ea,
		ExpressionAttributeValues: ev,
		IndexName:                 aws.String(index),
		KeyConditionExpression:    aws.String(strings.Join(kc, " AND ")),
		TableName:                 aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
	})
	if err != nil {
		return errors.WithStack(log.Error(err))
	}

	ritems := reflect.ValueOf(items).Elem()
	kind := reflect.TypeOf(items).Elem().Elem()

	for _, i := range res.Items {
		v := reflect.New(kind)

		if err := dynamique.Unmarshal(i, v.Interface(), d.key); err != nil {
			return errors.WithStack(log.Error(err))
		}

		if d, ok := v.Interface().(Defaulter); ok {
			d.Defaults()
		}

		ritems.Set(reflect.Append(ritems, v.Elem()))
	}

	if s, ok := reflect.ValueOf(items).Elem().Interface().(Sortable); ok {
		sort.Slice(s, s.Less)
	}

	return log.Success()
}

func (d *Dynamo) List(table string, items interface{}) error {
	log := log.At("dynamo.list").Append("table=%s", table).Start()

	if r := reflect.TypeOf(items); r.Kind() != reflect.Ptr || r.Elem().Kind() != reflect.Slice {
		return errors.WithStack(log.Error(fmt.Errorf("items must be a pointer to a slice")))
	}

	req := &dynamodb.ScanInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
	}

	ctx := context.Background()

	p := request.Pagination{
		NewRequest: func() (*request.Request, error) {
			r, _ := d.DB.ScanRequest(req)
			r.SetContext(ctx)
			return r, nil
		},
	}

	ritems := reflect.ValueOf(items).Elem()
	kind := reflect.TypeOf(items).Elem().Elem()

	for p.Next() {
		page := p.Page().(*dynamodb.ScanOutput)

		for _, i := range page.Items {
			v := reflect.New(kind)

			if err := dynamique.Unmarshal(i, v.Interface(), d.key); err != nil {
				return errors.WithStack(log.Error(err))
			}

			if d, ok := v.Interface().(Defaulter); ok {
				d.Defaults()
			}

			ritems.Set(reflect.Append(ritems, v.Elem()))
		}
	}

	if s, ok := reflect.ValueOf(items).Elem().Interface().(Sortable); ok {
		sort.Slice(s, s.Less)
	}

	if err := p.Err(); err != nil {
		return errors.WithStack(log.Error(err))
	}

	return log.Success()
}

func (d *Dynamo) Put(table string, item interface{}) error {
	log := log.At("dynamo.put").Append("table=%s", table).Start()

	if r := reflect.TypeOf(item); r.Kind() != reflect.Ptr {
		return errors.WithStack(log.Error(fmt.Errorf("item must be a pointer")))
	}

	if d, ok := item.(Defaulter); ok {
		d.Defaults()
	}

	if v, ok := item.(Validator); ok {
		if errs := v.Validate(); len(errs) > 0 {
			msgs := []string{}

			for _, err := range errs {
				msgs = append(msgs, err.Error())
			}

			return errors.New(strings.Join(msgs, ", "))
		}
	}

	attrs, err := dynamique.Marshal(item, d.key)
	if err != nil {
		return errors.WithStack(log.Error(err))
	}

	_, err = d.DB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
		Item:      attrs,
	})
	if err != nil {
		return errors.WithStack(log.Error(err))
	}

	return log.Success()
}

func (d *Dynamo) Query(table string, key map[string]string, opts QueryOptions, items interface{}) error {
	log := log.At("dynamo.query").Append("table=%s key=%v", table, key).Start()

	if r := reflect.TypeOf(items); r.Kind() != reflect.Ptr || r.Elem().Kind() != reflect.Slice {
		return errors.WithStack(log.Error(fmt.Errorf("items must be a pointer to a slice")))
	}

	kc := []string{}
	ea := map[string]*string{}
	ev := map[string]*dynamodb.AttributeValue{}

	for k, v := range key {
		ks := strings.Replace(k, "-", "", -1)
		kc = append(kc, fmt.Sprintf("#%s=:%s", ks, ks))
		ea[fmt.Sprintf("#%s", ks)] = aws.String(k)
		ev[fmt.Sprintf(":%s", ks)] = &dynamodb.AttributeValue{S: aws.String(v)}
	}

	req := &dynamodb.QueryInput{
		ExpressionAttributeNames:  ea,
		ExpressionAttributeValues: ev,
		IndexName:                 opts.Index,
		KeyConditionExpression:    aws.String(strings.Join(kc, " AND ")),
		Limit:                     opts.Limit,
		ScanIndexForward:          opts.Forward,
		TableName:                 aws.String(fmt.Sprintf("%s-%s", d.prefix, table)),
	}

	if len(ea) > 0 {
		req.ExpressionAttributeNames = ea
	}

	if len(ev) > 0 {
		req.ExpressionAttributeValues = ev
	}

	if opts.Start != nil {
		req.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{}

		for k, v := range opts.Start {
			req.ExclusiveStartKey[k] = &dynamodb.AttributeValue{S: aws.String(v)}
		}
	}

	res, err := d.DB.Query(req)
	if err != nil {
		return errors.WithStack(log.Error(err))
	}

	ritems := reflect.ValueOf(items).Elem()
	kind := reflect.TypeOf(items).Elem().Elem()

	for _, i := range res.Items {
		v := reflect.New(kind)

		if err := dynamique.Unmarshal(i, v.Interface(), d.key); err != nil {
			return errors.WithStack(log.Error(err))
		}

		if d, ok := v.Interface().(Defaulter); ok {
			d.Defaults()
		}

		ritems.Set(reflect.Append(ritems, v.Elem()))
	}

	if s, ok := reflect.ValueOf(items).Elem().Interface().(Sortable); ok {
		sort.Slice(s, s.Less)
	}

	return log.Success()
}

func paginate(items []string, max int) [][]string {
	pages := [][]string{}

	for len(items) > max {
		pages = append(pages, items[:max])
		items = items[max:]
	}

	pages = append(pages, items)

	return pages
}
