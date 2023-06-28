package jobs

import (
	"context"
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
	Address         *models.Address
	cl              *http.Client
	balances        map[int64]schemas.TokenMapping
	pairBalRetries  int8
	tokenBalRetries int8
	Errs            error
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
	// STUB - sync user's th
	conf.TransactionUpdater = append(conf.TransactionUpdater, u.mustQueryTrx)
	// STUB - sync user's nft
	conf.NFTUpdater = append(conf.NFTUpdater, u.mustQueryNFTs)

	// STUB - sync user's bal
	u.tokenBalRetries = 3
	if err := u.tokenBal(c); err != nil {
		return err
	}
	if err := u.saveTokenBal(c); err != nil {
		return err
	}
	u.pairBalRetries = 3
	if err := u.pairBal(c); err != nil {
		return err
	}
	if err := u.savePairBal(c); err != nil {
		return err
	}

	// STUB - sync user's approve

	if tx := conf.DB.Save(&u.Address); tx.Error != nil {
		return tx.Error
	}
	return nil
}
