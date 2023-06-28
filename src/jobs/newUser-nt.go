package jobs

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/PiperFinance/UA/src/conf"
)

func (u *SyncAddress) mustQueryNFTs() {
	c, cancel := context.WithTimeout(context.TODO(), conf.Config.NTSaveTimeout)
	defer cancel()
	if err := u.queryNFTs(c); err != nil {
		conf.Logger.Error(err)
	}
}

func (u *SyncAddress) queryNFTs(c context.Context) error {
	url := conf.Config.NT_URL.JoinPath(conf.Config.NTSaveNFTs)
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
