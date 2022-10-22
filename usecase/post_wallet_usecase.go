package usecase

import (
	"fmt"
	"time"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

func (uc *Usecase) PostWalletUsecase() (interface{}, error) {
	var topic goka.Stream = "user-click"

	emitter, err := goka.NewEmitter([]string{"127.0.0.1:9092"}, topic, new(codec.String))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	key := fmt.Sprintf("user")
	value := fmt.Sprintf("%s", time.Now())
	emitter.EmitSync(key, value)

	return map[string]interface{}{"acknowledge": true}, nil
}
