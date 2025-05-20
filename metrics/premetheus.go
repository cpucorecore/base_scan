package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	defaultMaxAge     = time.Second * 10
	defaultAgeBuckets = uint32(60)
	defaultObjectives = map[float64]float64{
		0.9:  0.01,
		0.99: 0.001,
	}
)

var (
	CurrentHeight = prometheus.NewGauge(prometheus.GaugeOpts{Name: "current_height"})
	NewestHeight  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "newest_height"})

	GetBlockDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_duration_ms",
		Help:       "get_block duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetBlockReceiptsDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_receipts_duration_ms",
		Help:       "get block receipts duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	DbOperationDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "db_operation_duration_ms",
		Help:       "db operation duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	ParseBlockDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "parse_block_duration_ms",
		Help:       "parse block duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	SendBlockKafkaDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "send_block_kafka_duration_ms",
		Help:       "send block kafka duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	CallContractDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_duration_ms",
		Help:       "call contract duration in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	BlockDelayMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "block_delay_ms",
		Help:       "block delay in Milliseconds",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	CallContractForNativeTokenPriceDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "call_contract_price_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetV2PairDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_v2_pair_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetV3PairDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_v3_pair_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})

	GetTokenDurationMs = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_token_duration_ms",
		MaxAge:     defaultMaxAge,
		AgeBuckets: defaultAgeBuckets,
		Objectives: defaultObjectives,
	})
)

func init() {
	prometheus.MustRegister(CurrentHeight)
	prometheus.MustRegister(NewestHeight)

	prometheus.MustRegister(GetBlockDurationMs)
	prometheus.MustRegister(GetBlockReceiptsDurationMs)

	prometheus.MustRegister(DbOperationDurationMs)
	prometheus.MustRegister(ParseBlockDurationMs)
	prometheus.MustRegister(SendBlockKafkaDurationMs)

	prometheus.MustRegister(CallContractDurationMs)
	prometheus.MustRegister(BlockDelayMs)
	prometheus.MustRegister(CallContractForNativeTokenPriceDurationMs)
	prometheus.MustRegister(GetV2PairDurationMs)
	prometheus.MustRegister(GetV3PairDurationMs)
	prometheus.MustRegister(GetTokenDurationMs)
}

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", 9100), nil)
	}()
}
