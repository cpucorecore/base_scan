package config

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"os"
	"time"
)

type LogConf struct {
	Async                      bool `json:"async"`
	AsyncBufferSizeByByte      int  `json:"async_buffer_size_by_byte"`
	AsyncFlushIntervalBySecond int  `json:"async_flush_interval_by_second"`
}

type ChainConf struct {
	Endpoint        string `json:"endpoint"`
	EndpointArchive string `json:"endpoint_archive"`
	WsEndpoint      string `json:"ws_endpoint"`
}

type RedisConf struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BlockGetterConf struct {
	PoolSize         int       `json:"pool_size"`
	QueueSize        int       `json:"queue_size"`
	StartBlockNumber uint64    `json:"start_block_number"`
	Retry            RetryConf `json:"retry"`
}

type BlockHandlerConf struct {
	PoolSize  int `json:"pool_size"`
	QueueSize int `json:"queue_size"`
}

type RetryConf struct {
	Attempts  uint `json:"attempts"`
	DelayMs   int  `json:"delay_ms"`
	TimeoutMs int  `json:"timeout_ms"`
}

func (rc *RetryConf) GetRetryParams() *RetryParams {
	return &RetryParams{
		Attempts: retry.Attempts(rc.Attempts),
		Delay:    retry.Delay(time.Duration(rc.DelayMs) * time.Millisecond),
		Timeout:  time.Duration(rc.TimeoutMs) * time.Millisecond,
	}
}

type RetryParams struct {
	Attempts retry.Option
	Delay    retry.Option
	Timeout  time.Duration
}

type PriceServiceConf struct {
	PoolSize int `json:"pool_size"`
}

type KafkaConf struct {
	On                bool
	Brokers           []string
	Topic             string
	SendTimeoutByMs   int
	MaxRetry          int
	RetryIntervalByMs int
}

type ContractCallerConf struct {
	Retry *RetryConf
}

type DbConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

func (dc *DbConf) GetDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dc.Host, dc.Username, dc.Password, dc.Db, dc.Port)
}

type Config struct {
	Log             *LogConf            `json:"log"`
	Chain           *ChainConf          `json:"chain"`
	Redis           *RedisConf          `json:"redis"`
	BlockGetter     *BlockGetterConf    `json:"block_getter"`
	BlockHandler    *BlockHandlerConf   `json:"block_handler"`
	EnableSequencer bool                `json:"enable_sequencer"`
	PriceService    *PriceServiceConf   `json:"price_service"`
	Kafka           *KafkaConf          `json:"kafka"`
	ContractCaller  *ContractCallerConf `json:"contract_caller"`
	DbTx            *DbConf             `json:"db_tx"`
	DbTokenPair     *DbConf             `json:"db_token_pair"`
}

var (
	defaultConfig = Config{
		Log: &LogConf{
			Async:                      false,
			AsyncBufferSizeByByte:      1000000,
			AsyncFlushIntervalBySecond: 1,
		},
		Chain: &ChainConf{
			Endpoint:        "https://base-rpc.publicnode.com",
			EndpointArchive: "https://base-rpc.publicnode.com",
			WsEndpoint:      "wss://base-rpc.publicnode.com",
		},
		Redis: &RedisConf{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		},
		BlockGetter: &BlockGetterConf{
			PoolSize:         1,
			QueueSize:        1,
			StartBlockNumber: 48000000,
			Retry: RetryConf{
				Attempts:  10,
				DelayMs:   100,
				TimeoutMs: 5000,
			},
		},
		BlockHandler: &BlockHandlerConf{
			PoolSize:  1,
			QueueSize: 1,
		},
		EnableSequencer: true,
		PriceService: &PriceServiceConf{
			PoolSize: 1,
		},
		Kafka: &KafkaConf{
			On:                false,
			Brokers:           []string{"localhost:9092"},
			Topic:             "block",
			SendTimeoutByMs:   5000,
			MaxRetry:          10,
			RetryIntervalByMs: 100,
		},
		ContractCaller: &ContractCallerConf{
			Retry: &RetryConf{
				Attempts:  10,
				DelayMs:   100,
				TimeoutMs: 3000,
			},
		},
		DbTx: &DbConf{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Db:       "test",
		},
		DbTokenPair: &DbConf{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Db:       "test",
		},
	}

	G = defaultConfig
)

func LoadConfigFile(configFilePath string) error {
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&G); err != nil {
		return err
	}

	return nil
}
