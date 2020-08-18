package queue

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type Queue struct {
	sqs *sqs.SQS
	url string
}

type WaitFunc func(map[string]string) (bool, error)

func New(url string) *Queue {
	return &Queue{
		sqs: sqs.New(session.New()),
		url: url,
	}
}

func (q *Queue) Delete(m *sqs.Message) error {
	_, err := q.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.url),
		ReceiptHandle: m.ReceiptHandle,
	})
	return errors.WithStack(err)
}

func (q *Queue) Dequeue(wait WaitFunc) (map[string]string, error) {
	for {
		res, err := q.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames:      []*string{aws.String("ApproximateReceiveCount")},
			QueueUrl:            aws.String(q.url),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(10),
			VisibilityTimeout:   aws.Int64(10),
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if len(res.Messages) < 1 {
			continue
		}

		message := res.Messages[0]

		if cs := message.Attributes["ApproximateReceiveCount"]; cs != nil {
			c, err := strconv.Atoi(*cs)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if c > 720 {
				q.Delete(message)
				continue
			}
		}

		job, err := parseMessage(message)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		w, err := wait(job)
		if err != nil {
			fmt.Printf("wait error: %s\n", err)
			q.Delete(message)
		}

		if w {
			continue
		}

		if err := q.Delete(message); err != nil {
			return nil, errors.WithStack(err)
		}

		return job, nil
	}
}

func (q *Queue) Enqueue(id, group string, job map[string]string) error {
	data, err := json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = q.sqs.SendMessage(&sqs.SendMessageInput{
		QueueUrl:               aws.String(q.url),
		MessageBody:            aws.String(string(data)),
		MessageDeduplicationId: aws.String(id),
		MessageGroupId:         aws.String(fmt.Sprintf("%s:%s", job["type"], group)),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func parseMessage(m *sqs.Message) (map[string]string, error) {
	if m.Body == nil {
		return nil, errors.WithStack(fmt.Errorf("message has no body"))
	}

	var job map[string]string

	if err := json.Unmarshal([]byte(*m.Body), &job); err != nil {
		return nil, errors.WithStack(err)
	}

	return job, nil
}
