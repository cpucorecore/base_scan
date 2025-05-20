package service

import (
	"base_scan/config"
	"base_scan/log"
	"base_scan/metrics"
	"base_scan/types"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type KafkaSender interface {
	Send(block *types.EthBlock) error
	SendOld(block *types.EthBlockOld) error
}

type kafkaSender struct {
	ID            string
	conf          *config.KafkaConf
	sendTimeout   time.Duration
	asyncProducer sarama.AsyncProducer
}

func NewKafkaSender(conf *config.KafkaConf) KafkaSender {
	client := &kafkaSender{
		conf:        conf,
		sendTimeout: time.Millisecond * time.Duration(conf.SendTimeoutByMs),
	}

	sc := sarama.NewConfig()
	sc.Net.TLS.Enable = false
	sc.Producer.Return.Errors = true
	sc.Producer.RequiredAcks = sarama.WaitForLocal
	sc.Producer.Compression = sarama.CompressionSnappy
	sc.Producer.Flush.Frequency = 100 * time.Millisecond
	sc.Producer.Retry.Max = 10

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, sc)
	if err != nil {
		log.Logger.Fatal("kafka NewAsyncProducer err", zap.Error(err))
	}
	client.asyncProducer = asyncProducer
	client.processErrors()

	return client
}

func (c *kafkaSender) Close() {
	_ = c.asyncProducer.Close()
}

func (c *kafkaSender) processErrors() {
	errCh := c.asyncProducer.Errors()
	go func() {
		for {
			err, ok := <-errCh
			if !ok {
				log.Logger.Info("kafka asyncProducer error @ done", zap.Error(err))
				return
			}
			log.Logger.Info("kafka asyncProducer error", zap.Error(err))
		}
	}()
}

func (c *kafkaSender) Send(block *types.EthBlock) error {
	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %v, %v", err, block)
	}

	now := time.Now()
	c.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: c.conf.Topic,
		Value: sarama.ByteEncoder(data),
	}
	metrics.SendBlockKafkaDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	return nil
}

func (c *kafkaSender) SendOld(block *types.EthBlockOld) error {
	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %v, %v", err, block)
	}

	now := time.Now()
	c.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: c.conf.Topic,
		Value: sarama.ByteEncoder(data),
	}
	metrics.SendBlockKafkaDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	return nil
}
