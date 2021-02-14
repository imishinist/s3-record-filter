package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type FilterConfig struct {
	RecordAfter Time   `envconfig:"RECORD_AFTER"`
	QueueURL    string `envconfig:"QUEUE_URL"`
}

func createSNSMessage(records events.S3Event) (string, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(records); err != nil {
		return "", xerrors.Errorf("create sns message error: %w", err)
	}

	snsMessage := new(bytes.Buffer)
	if err := json.NewEncoder(snsMessage).Encode(events.SNSEntity{
		Message: buf.String(),
	}); err != nil {
		return "", xerrors.Errorf("create sns message error: %w", err)
	}
	return snsMessage.String(), nil
}

func printJson(out io.Writer, data interface{}) {
	if err := json.NewEncoder(out).Encode(data); err != nil {
		log.Println(err)
	}
}

func handler(ctx context.Context, event events.SQSEvent) error {
	var f FilterConfig
	if err := envconfig.Process("filter", &f); err != nil {
		log.Fatal(err)
	}
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)

	s3Records := make([]events.S3EventRecord, 0, 10)
	for _, r := range event.Records {
		body := r.Body
		r.Body = ""

		var record events.SNSEntity
		if err := json.Unmarshal([]byte(body), &record); err != nil {
			return fmt.Errorf("unmarshal sns entity error: %w", err)
		}

		var records events.S3Event
		if err := json.Unmarshal([]byte(record.Message), &records); err != nil {
			return fmt.Errorf("unmarshal s3 event error: %w", err)
		}

		if f.RecordAfter.IsZero() {
			fmt.Println("`RECORD_AFTER` is nil, so append all")
			s3Records = append(s3Records, records.Records...)
			continue
		}

		// filter s3 record
		for _, r := range records.Records {
			if r.EventTime.After(f.RecordAfter.ToTime()) {
				continue
			}
			s3Records = append(s3Records, r)
		}
	}
	if len(s3Records) == 0 {
		return nil
	}

	message, err := createSNSMessage(events.S3Event{Records: s3Records})
	if err != nil {
		return err
	}
	if _, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(f.QueueURL),
		MessageBody: aws.String(message),
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
