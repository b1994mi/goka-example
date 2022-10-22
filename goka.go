package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers             = []string{"127.0.0.1:9092"}
	topic   goka.Stream = "deposits"

	aboveThresholdGroup goka.Group = "above-threshold-group"
	balanceGroup        goka.Group = "balance-group"
)

type user struct {
	// number of clicks the user has performed.
	Clicks int
}

func aboveThresholdProcess(ctx goka.Context, msg interface{}) {
	var u *user
	if val := ctx.Value(); val != nil {
		u = val.(*user)
	} else {
		u = new(user)
	}

	u.Clicks++
	ctx.SetValue(u)
	fmt.Printf("[proc] key: %s clicks: %d, msg: %v\n", ctx.Key(), u.Clicks, msg)
}

func balanceProcess(ctx goka.Context, msg interface{}) {
	var u *user
	if val := ctx.Value(); val != nil {
		u = val.(*user)
	} else {
		u = new(user)
	}

	u.Clicks++
	ctx.SetValue(u)
	fmt.Printf("Hello there, this is the second consumer\n")
	fmt.Printf("[proc2] key2: %s clicks2: %d, msg2: %v\n", ctx.Key(), u.Clicks, msg)
}

func initGokaProcessor() {
	tmc := goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	defer tm.Close()
	err = tm.EnsureStreamExists(string(topic), 8)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", topic, err)
	}

	aboveThresholdProcessor, err := goka.NewProcessor(
		brokers,
		goka.DefineGroup(
			aboveThresholdGroup,
			goka.Input(topic, new(codec.String), aboveThresholdProcess),
			goka.Persist(new(userCodec)),
		),
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	go aboveThresholdProcessor.Run(context.Background())

	balanceProcessor, err := goka.NewProcessor(
		brokers,
		goka.DefineGroup(
			balanceGroup,
			goka.Input(topic, new(codec.String), balanceProcess),
			goka.Persist(new(userCodec)),
		),
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	go balanceProcessor.Run(context.Background())
}

type userCodec struct{}

// Encodes a user into []byte
func (jc *userCodec) Encode(value interface{}) ([]byte, error) {
	if _, isUser := value.(*user); !isUser {
		return nil, fmt.Errorf("Codec requires value *user, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a user from []byte to it's go representation.
func (jc *userCodec) Decode(data []byte) (interface{}, error) {
	var (
		c   user
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling user: %v", err)
	}
	return &c, nil
}

func initGokaView() *goka.View {
	gv, err := goka.NewView(brokers, goka.GroupTable(aboveThresholdGroup), new(userCodec))
	if err != nil {
		panic(err)
	}

	go gv.Run(context.Background())

	return gv
}

func initGokaView2() *goka.View {
	gv, err := goka.NewView(brokers, goka.GroupTable(balanceGroup), new(userCodec))
	if err != nil {
		panic(err)
	}

	go gv.Run(context.Background())

	return gv
}
