package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	CurrentHeight = prometheus.NewGauge(prometheus.GaugeOpts{Name: "current_height"})
	NewestHeight  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "newest_height"})

	GetBlockDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_duration",
		Help:       "get_block duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetBlockReceiptsDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_receipts_duration",
		Help:       "get block receipts duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	DbOperationDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "db_operation_duration",
		Help:       "db operation duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	ParseBlockDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "parse_block_duration",
		Help:       "parse block duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	SendBlockKafkaDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "send_block_kafka_duration",
		Help:       "send block kafka duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	CallContractDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_duration",
		Help:       "call contract duration in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 100,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	BlockDelay = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "block_delay",
		Help:       "block delay in seconds",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	CallContractForNativeTokenPrice = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_price",
		MaxAge:     time.Minute,
		AgeBuckets: 20,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetV2PairDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_v2_pair_duration",
		MaxAge:     time.Minute * 10,
		AgeBuckets: 100,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetV3PairDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_v3_pair_duration",
		MaxAge:     time.Minute * 10,
		AgeBuckets: 100,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetTokenDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_token_duration",
		MaxAge:     time.Minute * 10,
		AgeBuckets: 100,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})
)

func init() {
	prometheus.MustRegister(CurrentHeight)
	prometheus.MustRegister(NewestHeight)

	prometheus.MustRegister(GetBlockDuration)
	prometheus.MustRegister(GetBlockReceiptsDuration)

	prometheus.MustRegister(DbOperationDuration)
	prometheus.MustRegister(ParseBlockDuration)
	prometheus.MustRegister(SendBlockKafkaDuration)

	prometheus.MustRegister(CallContractDuration)
	prometheus.MustRegister(BlockDelay)
	prometheus.MustRegister(CallContractForNativeTokenPrice)
	prometheus.MustRegister(GetV2PairDuration)
	prometheus.MustRegister(GetV3PairDuration)
	prometheus.MustRegister(GetTokenDuration)
}

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", 9100), nil)
	}()
}
