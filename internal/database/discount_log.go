package database

import "database/sql"

type DiscountLogTable DB

type DiscountLog struct {
	TrxId                string `db:"trx_id"`
	MerchantWalletNumber string `db:"merchant_wallet_number"`
	DiscountAggregatorId string `db:"discount_aggregator_status"`
	DiscountAmount       string `db:"discount_amount"`
	FtStatus             string `db:"ft_status"`
	FtBatchId            string `db:"ft_batch_id"`
	FtTransactionId      string `db:"ft_transaction_id"`
}

func (db *DiscountLogTable) Create(trxId string, discountAggregatorId string, merchantWalletNumber string, discountAmount string) bool {
	query := `Insert into discount_log (trx_id, discount_aggregator_id, merchant_wallet_number, discount_amount) Values ($1,$2,$3,$4)`
	_, err := db.Exec(query, trxId, discountAggregatorId, merchantWalletNumber, discountAmount)
	if err != nil {
		logger.Error("DB_EXECUTION_ERROR %+v\n ", err.Error())
		return false
	}
	return true
}
func (db *DiscountLogTable) Update(trxId string, ftStatus string, ftBatchId *string, ftTransactionId *string, vendorStatus *string) bool {
	query := `Update discount_log set ft_status=$1, ft_batch_id=$2, ft_transaction_id=$3 , vendor_status=$5 where trx_id=$4`
	_, err := db.Exec(query, ftStatus, ftBatchId, ftTransactionId, trxId, vendorStatus)
	if err != nil {
		logger.Error("DB_EXECUTION_ERROR %+v\n", err.Error())
		return false
	}
	return true
}

func (db *DiscountLogTable) List(limit *int, offset *int) (*sql.Rows, bool) {
	query := `Select * from discount_log order by id LIMIT $1 OFFSET $2`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		logger.Error("DB_EXECUTION_ERROR %+v\n", err.Error())
		return nil, false
	}
	return rows, true
}
