package kafka

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	ccid "github.com/kperanovic/epam-systems/internal/cid"
	logger "github.com/kperanovic/epam-systems/internal/logger"
	"go.uber.org/zap"
)

type Producer struct {
	log      *zap.Logger
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string, topic string, log *zap.Logger) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		log:      log,
		producer: producer,
	}, nil
}

// NewMockProducer will create a new kafka producer
// with `mocks.SyncProducer` as the underlying sarama producer.
func NewMockProducer(producer *mocks.SyncProducer, log *zap.Logger) *Producer {
	return &Producer{
		producer: producer,
		log:      log,
	}
}

// Close will close the producer and log an error if it happened.
func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		p.log.Error("error closing producer", zap.Error(err))
	}
}

// SendMessage sends proto message to a given topic
func (p *Producer) SendMessage(ctx context.Context, topic string, partitionKey string, msgName string, msg interface{}) error {
	ctx, cid := ccid.FromContextOrNew(ctx)
	log := p.getLogger(ctx, cid)

	headers := p.createHeaders(msgName, cid)

	toSend, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	m := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(partitionKey),
		Value:   sarama.ByteEncoder(toSend),
		Headers: p.createRecordHeaders(headers),
	}

	partition, offset, err := p.producer.SendMessage(m)
	if err != nil {
		return err
	}

	log.Info("message sent", zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}

// createHeaders will create the kafka headers.
func (p *Producer) createHeaders(msgName, cid string) map[string]string {
	return map[string]string{
		MessageNameHeader: msgName,
		CIDHeader:         cid,
	}
}

// getLogger will try to load the logger from the context.
// If the logger doesn't exists it will load the instance logger.
// The function will then add the correlation id filed to it
// and return it
func (p *Producer) getLogger(ctx context.Context, cid string) *zap.Logger {
	log := logger.FromContext(ctx)
	if log == nil {
		log = p.log
	}
	return log.With(zap.String("cid", cid))
}

// createRecordHeaders converts the headers map
// to a sarama.RecordHeader array which is required
// by the sarama.Message struct
func (p *Producer) createRecordHeaders(headers map[string]string) []sarama.RecordHeader {
	var rh []sarama.RecordHeader
	for k, v := range headers {
		rh = append(rh, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	return rh
}
