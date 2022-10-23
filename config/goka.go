package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"goka-example/model"

	"github.com/lovoo/goka"
)

var (
	Brokers             = []string{"127.0.0.1:9092"}
	Topic   goka.Stream = "deposits"

	aboveThresholdGroup goka.Group = "above-threshold-group1"
	balanceGroup        goka.Group = "balance-group1"
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
	fmt.Println(wt)
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
	if _, isWallet := value.(*model.Wallet); !isWallet {
		return nil, fmt.Errorf("Codec requires value *Wallet, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a wallet from []byte to it's go representation.
func (jc *walletCodec) Decode(data []byte) (interface{}, error) {
	var w model.Wallet
	err := json.Unmarshal(data, &w)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling wallet: %v", err)
	}
	return &w, nil
}

type walletThresholdCodec struct{}

// Encodes a wallet threshold into []byte
func (jc *walletThresholdCodec) Encode(value interface{}) ([]byte, error) {
	if _, isWalletThreshold := value.(*model.WalletThreshold); !isWalletThreshold {
		return nil, fmt.Errorf("Codec requires value *WalletThreshold, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a wallet threshold from []byte to it's go representation.
func (jc *walletThresholdCodec) Decode(data []byte) (interface{}, error) {
	var wt model.WalletThreshold
	err := json.Unmarshal(data, &wt)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling wallet threshold: %v", err)
	}
	return &wt, nil
}

type TransactionCodec struct{}

// Encodes a wallet threshold into []byte
func (jc *TransactionCodec) Encode(value interface{}) ([]byte, error) {
	if _, isTransaction := value.(*model.Transaction); !isTransaction {
		return nil, fmt.Errorf("Codec requires value *Transaction, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a wallet threshold from []byte to it's go representation.
func (jc *TransactionCodec) Decode(data []byte) (interface{}, error) {
	var t model.Transaction
	err := json.Unmarshal(data, &t)
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
