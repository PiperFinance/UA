package tasks

import (
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"

	"github.com/PiperFinance/UA/src/bg"
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
)

var ErrSyncedTooEarly = errors.New("sync minimum delay has not reached yet")

// ExecJob - getting user's balance and approves are needed to initialize one time and first time
func EnqueueSyncAdd(Address *models.Address) error {
	add, err := Address.ETHAddress()
	if err != nil {
		return err
	}
	user := schemas.SyncUser{Hash: add}
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}
	conf.QueueClient.Enqueue(
		asynq.NewTask(bg.SyncTokenBalTaskKey, payload),
		asynq.Queue(conf.UASyncBalQ),
		asynq.Timeout(conf.Config.PSV1TokenSyncTimeout))
	conf.QueueClient.Enqueue(
		asynq.NewTask(bg.SyncPairBalTaskKey, payload),
		asynq.Queue(conf.UASyncBalQ),
		asynq.Timeout(conf.Config.PSV1PairSyncTimeout))
	conf.QueueClient.Enqueue(
		asynq.NewTask(bg.SyncNTTaskKey, payload),
		asynq.Queue(conf.UASyncNTQ),
		asynq.Timeout(conf.Config.NTSaveTimeout))
	conf.QueueClient.Enqueue(
		asynq.NewTask(bg.SyncTHTaskKey, payload),
		asynq.Queue(conf.UASyncTHQ),
		asynq.Timeout(conf.Config.THSaveTimeout))

	if tx := conf.DB.Save(&Address); tx.Error != nil {
		return tx.Error
	}
	return nil
}
