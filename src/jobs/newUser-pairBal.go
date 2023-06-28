package jobs

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/schemas"
)

func (u *SyncAddress) pairBal(c context.Context) error {
	// TODO set block number in PS

	for u.pairBalRetries > 0 {
		url := conf.Config.PortfolioScannerURL.JoinPath(conf.Config.PSSafePairBalUrl)
		q := url.Query()
		for _, chain := range conf.Config.SupportedChains {
			q.Add("chainId", fmt.Sprint(chain))
		}
		q.Add("wallet", u.Address.Hash)
		url.RawQuery = q.Encode()
		r, err := u.cl.Get(url.String())
		if err != nil {
			return err
		}
		if r.StatusCode >= 300 {
			time.Sleep(5 * time.Second)
			u.pairBalRetries--
			continue
		}
		defer r.Body.Close()
		userBal := make(map[int64]schemas.TokenMapping)
		if err := json.NewDecoder(r.Body).Decode(&userBal); err != nil {
			return err
		}
		u.balances = userBal
		break
	}
	return nil
}

func (u *SyncAddress) savePairBal(c context.Context) error {
	if u.pairBalRetries == 0 {
		// NOTE - this means getting balance was unsuccessful
		return nil
	}
	for _, chain := range conf.Config.SupportedChains {
		url := conf.Config.BlockScannerURL.JoinPath(fmt.Sprintf(conf.Config.BSSetBalURL, chain))
		if len(u.balances[chain]) < 1 {
			continue
		}
		r := make([]schemas.UserBalance, len(u.balances[chain]))
		x := 0
		for tokenId, userBal := range u.balances[chain] {
			add, _ := u.Address.ETHAddress()
			r[x] = schemas.UserBalance{
				TokenStr: userBal.Detail.Address.String(),
				UserStr:  u.Address.Hash,
				User:     add,
				Token:    userBal.Detail.Address,
				TokenId:  tokenId,
				Balance:  userBal.BalanceNoDecimalStr,
			}
			x++
		}
		body, _ := json.Marshal(r)
		if resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(body)); err != nil {
			return err
		} else if resp.StatusCode >= 300 {
			return fmt.Errorf("bs responded with %d, url %s resp %+v", resp.StatusCode, url.String(), resp)
		} else {
			conf.Logger.Debugw("BSSetBal", resp)
		}
	}
	u.Address.LastSync = sql.NullTime{Valid: true, Time: time.Now()}
	// conf.DB.Save(&u.Address)
	// conf.DB.Model(&models.Address{}).Where("hash = ?", u.Address.Hash).Update("LastSync", "hello")
	return nil
}

// func (u *SyncAddress) approves(c context.Context) error {
// 	return nil
// }
