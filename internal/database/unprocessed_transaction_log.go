package database

import (
	"encoding/json"
	"time"

	"example_project/internal/structs"

	"github.com/jmoiron/sqlx"
)

type UnprocessedTransactionLogStruct struct {
	ID         int        `db:"id"`
	TrxId      string     `db:"trx_id"`
	Payload    []byte     `db:"payload"`
	Remarks    string     `db:"remarks"`
	Status     string     `db:"status"`
	RetryCount int        `db:"retry_count"`
	UpdatedAt  *time.Time `db:"updated_at"`
	CreatedAt  *time.Time `db:"created_at"`
}

type UnprocessedTransactionLogTable DB

func (db *UnprocessedTransactionLogTable) Create(payload []byte, trx_id string, remarks string, createdAt time.Time, status string) bool {
	query := `Insert into unprocessed_transaction_log (trx_id, payload, remarks, created_at, status) Values ($1,$2,$3,$4,$5)`
	_, err := db.Exec(query, trx_id, payload, remarks, createdAt, status)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	return true
}

func (db *UnprocessedTransactionLogTable) GetOne(maxRetry int, status string) (int, *sqlx.Tx, *structs.TransactionCompletePayload, error) {
	tx, err := db.Beginx()
	if err != nil {
		logger.Error("failed to begin a db transaction, details: %v", err.Error())
		return 0, nil, nil, err
	}
	query := `SELECT * FROM unprocessed_transaction_log WHERE retry_count < $1 AND status=$2 ORDER BY created_at DESC LIMIT 1 FOR UPDATE SKIP LOCKED;`

	unprocessedTransactionLog := &UnprocessedTransactionLogStruct{}
	err = tx.Get(unprocessedTransactionLog, query, maxRetry, status)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			logger.Error("failed to rollback transaction, details: %v", err.Error())
		}
		return 0, nil, nil, err
	}

	logger.Debug("unprocessed transaction log selected for reprocessing: %v", unprocessedTransactionLog)
	trxPayload := &structs.TransactionCompletePayload{}

	err = json.Unmarshal(unprocessedTransactionLog.Payload, trxPayload)
	if err != nil {
		logger.Error("error happened during un-marshalling payload from database, details: %v", err.Error())
		if err = tx.Rollback(); err != nil {
			logger.Error("failed to rollback transaction, details: %v", err.Error())
		}
		return 0, nil, nil, err
	}
	return unprocessedTransactionLog.ID, tx, trxPayload, nil
}

func (db *UnprocessedTransactionLogTable) UpdateRetryCount(tx *sqlx.Tx, id int) error {
	query := `UPDATE unprocessed_transaction_log SET retry_count = retry_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := tx.Exec(query, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (db *UnprocessedTransactionLogTable) UpdateLogStatus(tx *sqlx.Tx, id int, status string) error {
	query := `UPDATE unprocessed_transaction_log SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := tx.Exec(query, status, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
