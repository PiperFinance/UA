package jobs

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
)

// const safeBalUrl = "/v1/tokens/balance/safe"
var ErrSyncedTooEarly = errors.New("sync minimum delay has not reached yet")

type SyncAddress struct {
	// UUID *uuid.UUID `json:"UUID,omitempty"`
	Address    *models.Address
	cl         *http.Client
	balances   map[int64]schemas.TokenMapping
	balRetries int8
	Errs       error
}

// ExecJob - getting user's balance and approves are needed to initialize one time and first time
func (u *SyncAddress) ExecuteAll() error {
	return u.ExecuteAllWithContext(context.TODO())
}

func (u *SyncAddress) ExecuteAllWithContext(c context.Context) error {
	if u.Address.LastSync.Valid && u.Address.LastSync.Time.Add(conf.Config.PSSyncMinTimeDelay).Compare(time.Now()) == 1 {
		return ErrSyncedTooEarly
	}
	if u.cl != nil {
		return fmt.Errorf("already executed")
	}
	u.cl = &http.Client{}

	// STUB - sync user's bal
	u.balRetries = 2
	if err := u.bal(c); err != nil {
		return err
	}
	if err := u.saveBal(c); err != nil {
		return err
	}
	// STUB - sync user's nft
	// STUB - sync user's th
	// STUB - sync user's approve

	if tx := conf.DB.Save(&u.Address); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u *SyncAddress) bal(c context.Context) error {
	// TODO set block number in PS

	// STUB - get user bal
	// STUB - retry on err
	// STUB - log err
	// STUB - save new balances
	// STUB - save last synced date
	for u.balRetries > 0 {
		url := conf.Config.PortfolioScannerURL.JoinPath(conf.Config.PSSafeBalUrl)
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
			u.balRetries--
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

func (u *SyncAddress) saveBal(c context.Context) error {
	if u.balRetries == 0 {
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
