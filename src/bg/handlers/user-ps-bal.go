package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hibiken/asynq"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/schemas"
)

func TokenBalTaskHandler(ctx context.Context, task *asynq.Task) error {
	user := schemas.SyncUser{}
	if err := json.Unmarshal(task.Payload(), &user); err != nil {
		return err
	}
	if err := TokenBal(ctx, conf.Config.PSSafeTokenBalUrl, user.Hash); err != nil {
		return err
	}
	// TOdO
	// Address.LastSync = sql.NullTime{Valid: true, Time: time.Now()}
	// return conf.DB.Save(&Address).Error
	return nil
}

func PairBalTaskHandler(ctx context.Context, task *asynq.Task) error {
	user := schemas.SyncUser{}
	if err := json.Unmarshal(task.Payload(), &user); err != nil {
		return err
	}
	if err := TokenBal(ctx, conf.Config.PSSafePairBalUrl, user.Hash); err != nil {
		return err
	}
	// TOdO
	// Address.LastSync = sql.NullTime{Valid: true, Time: time.Now()}
	// return conf.DB.Save(&Address).Error
	return nil
}

func TokenBal(c context.Context, getBalUrl string, Address common.Address) error {
	// STUB - save last synced date
	url := conf.Config.PortfolioScannerURL.JoinPath(getBalUrl)
	q := url.Query()
	cl := &http.Client{}
	for _, chain := range conf.Config.SupportedChains {
		q.Add("chainId", fmt.Sprint(chain))
	}
	q.Add("wallet", Address.String())
	url.RawQuery = q.Encode()
	r, err := cl.Get(url.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	userBal := make(map[int64]schemas.TokenMapping)
	if r.StatusCode >= 300 {
		b, _ := io.ReadAll(r.Body)
		return fmt.Errorf("bad response : [%d] %s", r.StatusCode, string(b))
	}

	if err := json.NewDecoder(r.Body).Decode(&userBal); err != nil {
		return err
	}
	return saveTokenBal(c, Address, userBal)
}

func saveTokenBal(c context.Context, Address common.Address, balances map[int64]schemas.TokenMapping) error {
	for _, chain := range conf.Config.SupportedChains {
		url := conf.Config.BlockScannerURL.JoinPath(fmt.Sprintf(conf.Config.BSSetBalURL, chain))
		if len(balances[chain]) < 1 {
			continue
		}
		r := make([]schemas.UserBalance, len(balances[chain]))
		x := 0
		for tokenId, userBal := range balances[chain] {

			r[x] = schemas.UserBalance{
				TokenStr: userBal.Detail.Address.String(),
				UserStr:  Address.String(),
				User:     Address,
				Token:    userBal.Detail.Address,
				TokenId:  tokenId,
				Balance:  userBal.BalanceNoDecimalStr,
			}
			x++
		}
		body, _ := json.Marshal(r)
		// TODO use context
		if resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(body)); err != nil {
			return err
		} else if resp.StatusCode >= 300 {
			body, _ = io.ReadAll(resp.Body)
			return fmt.Errorf("bs responded with %d, url %s resp %s", resp.StatusCode, url.String(), string(body))
		} else {
			conf.Logger.Debugw("BSSetBal", resp)
		}
	}
	return nil
}
