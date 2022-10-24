package usecase

import (
	"fmt"

	"goka-example/config"
	"goka-example/model"
	"goka-example/request"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (uc *Usecase) PostWalletUsecase(req request.PostWallet) (interface{}, error) {
	emitter, err := goka.NewEmitter(config.Brokers, config.Topic, new(config.TransactionCodec))
	if err != nil {
		return nil, err
	}
	defer emitter.Finish()

	err = emitter.EmitSync(
		fmt.Sprintf("%d", req.WalletID),
		&model.Transaction{
			Amount: req.Amount,
			Time:   timestamppb.Now(),
		},
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"acknowledge": true}, nil
}
