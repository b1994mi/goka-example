package usecase

import (
	"fmt"
	"time"

	"goka-example/config"
	"goka-example/model"
	"goka-example/request"

	"github.com/lovoo/goka"
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
			Time:   time.Now(),
		},
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"acknowledge": true}, nil
}
