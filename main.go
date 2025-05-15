package main

import (
	"base_scan/block_getter"
	"base_scan/cache"
	"base_scan/config"
	"base_scan/log"
	"base_scan/parser"
	"base_scan/repository"
	"base_scan/sequencer"
	"base_scan/service"
	"base_scan/service/contract_caller"
	"base_scan/types"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	time.Local = time.UTC

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version information")
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()

	if showVersion {
		fmt.Println(GetVersion())
		os.Exit(0)
	}

	log.Logger.Info(GetVersion().String())
	log.Logger.Info("config", zap.String("file path", configFile))
	loadConfigErr := config.LoadConfigFile(configFile)
	if loadConfigErr != nil {
		log.Logger.Fatal("load config file err", zap.Error(loadConfigErr))
	}

	ethClient, dialEthErr := ethclient.Dial(config.G.Chain.Endpoint)
	if dialEthErr != nil {
		log.Logger.Fatal("Failed to connect to the chain(http): %v", zap.Error(dialEthErr))
	}

	ethClientArchive, dialEthErrArchive := ethclient.Dial(config.G.Chain.EndpointArchive)
	if dialEthErrArchive != nil {
		log.Logger.Fatal("Failed to connect to the chain archive(http): %v", zap.Error(dialEthErrArchive))
	}

	wsEthClient, dialEthWsErr := ethclient.Dial(config.G.Chain.WsEndpoint)
	if dialEthWsErr != nil {
		log.Logger.Fatal("Failed to connect to the chain(ws): %v", zap.Error(dialEthWsErr))
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.G.Redis.Addr,
		Username: config.G.Redis.Username,
		Password: config.G.Redis.Password,
	})
	cache := cache.NewTwoTierCache(redisCli)

	contractCaller := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	pairService := service.NewPairService(cache, contractCaller)
	contractCallerArchive := contract_caller.NewContractCaller(ethClientArchive, config.G.ContractCaller.Retry.GetRetryParams())
	priceService := service.NewPriceService(cache, contractCallerArchive, ethClient, config.G.PriceService.PoolSize)

	blockHandlerSequencer := sequencer.NewBlockSequencer()

	topicRouter := parser.NewTopicRouter()
	kafkaSender := service.NewKafkaSender(config.G.Kafka)

	txDb, txDbErr := gorm.Open(postgres.Open(config.G.DbTx.GetDsn()))
	if txDbErr != nil {
		log.Logger.Fatal("failed to connect to db", zap.Error(txDbErr))
	}
	tokenPairDb, tokenPairDbErr := gorm.Open(postgres.Open(config.G.DbTokenPair.GetDsn()))
	if tokenPairDbErr != nil {
		log.Logger.Fatal("failed to connect to db", zap.Error(tokenPairDbErr))
	}

	tokenRepository := repository.NewTokenRepository(tokenPairDb)
	pairRepository := repository.NewPairRepository(tokenPairDb)
	txRepository := repository.NewTxRepository(txDb)
	dbService := service.NewDBService(tokenRepository, pairRepository, txRepository)
	blockParser := parser.NewBlockParser(
		cache,
		blockHandlerSequencer,
		priceService,
		pairService,
		topicRouter,
		kafkaSender,
		dbService,
		config.G.Kafka.On,
	)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	blockParser.Start(wg)

	blockGetterSequencer := sequencer.NewBlockSequencer()
	blockGetter := block_getter.NewBlockGetter(ethClient, wsEthClient, cache, blockGetterSequencer, config.G.BlockGetter.Retry.GetRetryParams())
	startBlockNumber := blockGetter.GetStartBlockNumber(config.G.BlockGetter.StartBlockNumber)
	if startBlockNumber == 0 {
		log.Logger.Fatal("start block number is zero")
	}

	blockGetterSequencer.Init(startBlockNumber - 1)
	blockHandlerSequencer.Init(startBlockNumber - 1)

	priceService.Start(startBlockNumber)
	blockGetter.Start()
	blockGetter.StartDispatch(startBlockNumber)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Logger.Info("receive signal", zap.String("signal", sig.String()))
		blockGetter.Stop()
	}()

	var blockCtx *types.BlockContext
	for {
		blockCtx = blockGetter.Next()
		if blockCtx == nil {
			log.Logger.Info("no more block to parse")
			blockParser.Stop()
			break
		}
		blockParser.ParseBlockAsync(blockCtx)
	}

	log.Logger.Info("wait all block commited")
	wg.Wait()
	log.Logger.Info("all block commited")
}
