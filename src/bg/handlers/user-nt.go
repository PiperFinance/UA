package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hibiken/asynq"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
)

func SyncNTFsTaskHandler(ctx context.Context, task *asynq.Task) error {
	user := schemas.SyncUser{}
	if err := json.Unmarshal(task.Payload(), &user); err != nil {
		return err
	}
	return queryNFTs(ctx, user.Hash)
}

func SyncNTFsScheduleTaskHandler(ctx context.Context, task *asynq.Task) error {
	var addresses []*models.Address
	if res := conf.DB.Model(&models.Address{}).Find(&addresses); res.Error != nil {
		return res.Error
	} else {
		for _, add := range addresses {
			_add, err := add.ETHAddress()
			if err != nil {
				return err
			}
			if err := queryNFTs(ctx, _add); err != nil {
				return err
			}
		}
	}
	_ = task
	return nil
}

func queryNFTs(c context.Context, Address common.Address) error {
	url := conf.Config.NT_URL.JoinPath(conf.Config.NTSaveNFTs)
	wg := sync.WaitGroup{}
	wg.Add(len(conf.Config.SupportedChains))
	cl := &http.Client{}
	for _, chain := range conf.Config.SupportedChains {
		go func(chain int64) {
			defer wg.Done()
			_, err := cl.Post(url.String(), "application/json", strings.NewReader(
				fmt.Sprintf("{\"chainIds\": [%d],\"userAddresses\": [\"%s\"]}", chain, Address.String()),
			))
			if err != nil {
				conf.Logger.Error(err)
			}
		}(chain)
	}
	wg.Wait()
	_ = c
	return nil
}
