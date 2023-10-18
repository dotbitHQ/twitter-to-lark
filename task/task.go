package task

import (
	"context"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"sync"
	"time"
	"twitter-to-lark/cache"
)

var log = logger.NewLogger("task", logger.LevelDebug)

type Task struct {
	Ctx      context.Context
	Wg       *sync.WaitGroup
	Rc       *cache.RedisCache
	MaxRetry int
}

// smt_status,tx_status: (2,1)->(3,3)
func (t *Task) Run() {

	tickerCheckAddressBalance := time.NewTicker(time.Second * 5)

	//tickerCheckAddressBalanceMultisignaddr := time.NewTicker(time.Minute * 10)
	//var dataSamplePayadd []BalanceSample
	//var dataSampleMultisignaddr []BalanceSample

	t.Wg.Add(1)
	go func() {
		for {
			select {

			case <-tickerCheckAddressBalance.C:
				err := t.doTwitter()
				if err != nil {
					log.Error("doCheckBalancePayaddr err:", err.Error())
				}

			case <-t.Ctx.Done():
				log.Info("task monitor done")
				t.Wg.Done()
				return
			}
		}
	}()
}
