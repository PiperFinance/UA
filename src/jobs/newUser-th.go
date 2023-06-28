package jobs

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/PiperFinance/UA/src/conf"
)

func (u *SyncAddress) mustQueryTrx() {
	c, cancel := context.WithTimeout(context.TODO(), conf.Config.THSaveTimeout)
	defer cancel()
	if err := u.queryTrx(c); err != nil {
		conf.Logger.Error(err)
	}
}

func (u *SyncAddress) queryTrx(c context.Context) error {
	url := conf.Config.TH_URL.JoinPath(conf.Config.THSaveTransactions)
	wg := sync.WaitGroup{}
	wg.Add(len(conf.Config.SupportedChains))
	for _, chain := range conf.Config.SupportedChains {
		go func(chain int64) {
			defer wg.Done()
			_, err := u.cl.Post(url.String(), "application/json", strings.NewReader(
				fmt.Sprintf("{\"chainIds\": [%d],\"userAddresses\": [\"%s\"]}", chain, u.Address.Hash),
			))
			if err != nil {
				conf.Logger.Error(err)
			}
		}(chain)
	}
	wg.Wait()
	return nil
}
