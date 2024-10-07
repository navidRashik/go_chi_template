package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"example_project/internal/structs"
)

type MerchantMapTable DB
type MerchantMapStruct struct {
	MerchantWalletNumber string       `db:"merchant_wallet_number"`
	IsActive             bool         `db:"is_active"`
	UpdatedAt            sql.NullTime `db:"updated_at"`
	CreatedAt            time.Time    `db:"created_at"`
}

type MerchantRepo interface {
	GetMerchantMap(string) (*MerchantMapStruct, error)
	UpdateMerchantStatus(merchantStatusUpdatePayload structs.MerchantStatusUpdatePayload) error
}

func NewMerchantRepo(db *sqlx.DB) MerchantRepo {
	return &DB{db}
}

func (db *DB) UpdateMerchantStatus(merchantStatusUpdatePayload structs.MerchantStatusUpdatePayload) error {
	// Begin transaction
	tx, err := db.Beginx()
	if err != nil {
		// logger.Error(err.Error())
		return err
	}

	// Get all the existing active merchants from the list provided tobe inserted
	existingNumbers := []string{}
	existsQuery, args, err := sqlx.In(`SELECT merchant_wallet_number FROM merchant_map WHERE merchant_wallet_number IN (?) AND is_active = true`, merchantStatusUpdatePayload.ActiveMerchant)
	if err != nil {
		tx.Rollback()
		// logger.Error(err.Error())
		return err
	}
	existsQuery = db.Rebind(existsQuery)
	err = tx.Select(&existingNumbers, existsQuery, args...)
	if err != nil {
		tx.Rollback()
		// logger.Error(err.Error())
		return err
	}

	// fmt.Printf("%+v", res)

	// Insert the new merchants
	tobeInsertedMerchants := []MerchantMapStruct{}
	for _, merchant := range merchantStatusUpdatePayload.ActiveMerchant {
		if !contains(existingNumbers, merchant) {
			tobeInsertedMerchants = append(tobeInsertedMerchants, MerchantMapStruct{
				MerchantWalletNumber: merchant,
				IsActive:             true,
				CreatedAt:            time.Now(),
			})
		}
	}
	if len(tobeInsertedMerchants) > 0 {
		_, err := tx.NamedExec(`INSERT INTO merchant_map (merchant_wallet_number, is_active, created_at) VALUES (:merchant_wallet_number, :is_active, :created_at);`, tobeInsertedMerchants)
		if err != nil {
			tx.Rollback()
			// logger.Error(err.Error())
			return err
		}
	}

	// set inactive merchants
	if len(merchantStatusUpdatePayload.InactiveMerchant) > 0 {
		query, args, err := sqlx.In(`UPDATE merchant_map SET is_active=?, updated_at=?  WHERE merchant_wallet_number IN (?) and is_active=?;`, false, time.Now(), merchantStatusUpdatePayload.InactiveMerchant, true)
		if err != nil {
			tx.Rollback()
			// logger.Error(err.Error())
			return err
		}

		query = db.Rebind(query)
		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			// logger.Error(err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		// logger.Error(err.Error())
		return err
	}

	return nil
}

func (db *DB) GetMerchantMap(wallet string) (*MerchantMapStruct, error) {
	query := `SELECT * FROM merchant_map WHERE merchant_wallet_number = ?`
	merchantData := new(MerchantMapStruct)
	err := db.QueryRowx(query, wallet).StructScan(&merchantData)
	if err != nil {
		return nil, err
	}
	return merchantData, nil
}

func (db *MerchantMapTable) IsMerchantExist(wallet string) (bool, error) {
	query := `SELECT COUNT(id) FROM merchant_map WHERE merchant_wallet_number = $1 AND is_active = true`
	var merchantExist int
	err := db.QueryRowx(query, wallet).Scan(&merchantExist)
	if err != nil {
		return false, err
	}
	return merchantExist > 0, nil
}

// Helper function to check if a string exists in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
