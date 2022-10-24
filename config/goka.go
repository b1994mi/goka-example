package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/b1994mi/goka-example/model"

	"github.com/golang/protobuf/proto"
	"github.com/lovoo/goka"
)

var (
	Brokers             = []string{"127.0.0.1:9092"}
	Topic   goka.Stream = "deposits"

	aboveThresholdGroup goka.Group = "above-threshold-group"
	balanceGroup        goka.Group = "balance-group"

	thresholdAmount float64 = 10000
)

func aboveThresholdProcess(ctx goka.Context, msg interface{}) {
	var wt *model.WalletThreshold
	if val := ctx.Value(); val != nil {
		wt = val.(*model.WalletThreshold)
	} else {
		wt = new(model.WalletThreshold)
	}

	m := msg.(*model.Transaction)

	wt.Transactions = append(wt.Transactions, m)

	// check 2 minute rolling-period for threshold
	lastTwoMin := time.Now().Add(-time.Minute * 2)
	var trxAmtInLastTwoMin float64
	for _, v := range wt.Transactions {
		if v.Time.AsTime().Before(lastTwoMin) {
			continue
		}

		trxAmtInLastTwoMin += v.Amount
	}
	fmt.Printf("trx in last two min: %v\n", trxAmtInLastTwoMin)
	if trxAmtInLastTwoMin > thresholdAmount {
		wt.IsAboveThreshold = true
	}

	ctx.SetValue(wt)
	fmt.Printf("[aboveThresholdProcess] key: %s, msg: %v\n", ctx.Key(), msg)
}

func balanceProcess(ctx goka.Context, msg interface{}) {
	var w *model.Wallet
	if val := ctx.Value(); val != nil {
		w = val.(*model.Wallet)
	} else {
		w = new(model.Wallet)
	}

	m := msg.(*model.Transaction)

	w.Balance += m.Amount
	ctx.SetValue(w)
	fmt.Printf("[balanceProcess] key: %s, msg: %v, w.Balance: %v\n", ctx.Key(), msg, w.Balance)
}

func InitGokaProcessor() {
	tmc := goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	tm, err := goka.NewTopicManager(Brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	defer tm.Close()
	err = tm.EnsureStreamExists(string(Topic), 8)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", Topic, err)
	}

	aboveThresholdProcessor, err := goka.NewProcessor(
		Brokers,
		goka.DefineGroup(
			aboveThresholdGroup,
			goka.Input(Topic, new(TransactionCodec), aboveThresholdProcess),
			goka.Persist(new(walletThresholdCodec)),
		),
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	go aboveThresholdProcessor.Run(context.Background())

	balanceProcessor, err := goka.NewProcessor(
		Brokers,
		goka.DefineGroup(
			balanceGroup,
			goka.Input(Topic, new(TransactionCodec), balanceProcess),
			goka.Persist(new(walletCodec)),
		),
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	go balanceProcessor.Run(context.Background())
}

type walletCodec struct{}

// Encodes a wallet into []byte
func (jc *walletCodec) Encode(value interface{}) ([]byte, error) {
	v, ok := value.(*model.Wallet)
	if !ok {
		return nil, fmt.Errorf("Codec requires value *Wallet, got %T", value)
	}
	return proto.Marshal(v)
}

// Decodes a wallet from []byte to it's go representation.
func (jc *walletCodec) Decode(data []byte) (interface{}, error) {
	var w model.Wallet
	err := proto.Unmarshal(data, &w)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling wallet: %v", err)
	}
	return &w, nil
}

type walletThresholdCodec struct{}

// Encodes a wallet threshold into []byte
func (jc *walletThresholdCodec) Encode(value interface{}) ([]byte, error) {
	v, ok := value.(*model.WalletThreshold)
	if !ok {
		return nil, fmt.Errorf("Codec requires value *WalletThreshold, got %T", value)
	}
	return proto.Marshal(v)
}

// Decodes a wallet threshold from []byte to it's go representation.
func (jc *walletThresholdCodec) Decode(data []byte) (interface{}, error) {
	var wt model.WalletThreshold
	err := proto.Unmarshal(data, &wt)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling wallet threshold: %v", err)
	}
	return &wt, nil
}

type TransactionCodec struct{}

// Encodes a wallet threshold into []byte
func (jc *TransactionCodec) Encode(value interface{}) ([]byte, error) {
	v, ok := value.(*model.Transaction)
	if !ok {
		return nil, fmt.Errorf("Codec requires value *Transaction, got %T", value)
	}
	return proto.Marshal(v)
}

// Decodes a wallet threshold from []byte to it's go representation.
func (jc *TransactionCodec) Decode(data []byte) (interface{}, error) {
	var t model.Transaction
	err := proto.Unmarshal(data, &t)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling Transaction: %v", err)
	}
	return &t, nil
}

func InitAboveThresholdView() *goka.View {
	gv, err := goka.NewView(Brokers, goka.GroupTable(aboveThresholdGroup), new(walletThresholdCodec))
	if err != nil {
		panic(err)
	}

	go gv.Run(context.Background())

	return gv
}

func InitBalanceView() *goka.View {
	gv, err := goka.NewView(Brokers, goka.GroupTable(balanceGroup), new(walletCodec))
	if err != nil {
		panic(err)
	}

	go gv.Run(context.Background())

	return gv
}
