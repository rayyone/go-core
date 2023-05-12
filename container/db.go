package corecontainer

import (
	"context"
	"github.com/rayyone/go-core/ryerr"
	"gorm.io/gorm"
)

type Database struct {
	db                *gorm.DB
	dbTransaction     *gorm.DB
	transactionOpened bool
}

func (d *Database) GetTx() *gorm.DB {
	if !d.transactionOpened {
		return d.db
	}
	return d.dbTransaction
}

func (d *Database) BeginTransaction() {
	d.dbTransaction = d.db.Begin()
	d.transactionOpened = true
}

func (d *Database) SetContext(ctx context.Context) {
	d.db = d.db.WithContext(ctx)
}

func (d *Database) Commit() error {
	if !d.transactionOpened {
		return ryerr.New("TX has been committed or rolled back")
	}
	err := d.dbTransaction.Commit().Error
	if err != nil {
		err = ryerr.Newf("Commit Error. Error: %v", err)
		_ = d.Rollback()
		return err
	}
	d.Clear()
	return nil
}

func (d *Database) Rollback() error {
	if !d.transactionOpened {
		// TX has been committed or rolled back
		return nil
	}
	err := d.dbTransaction.Rollback().Error
	if err != nil {
		err = ryerr.Newf("Rollback Error. Error: %v", err)
		return err
	}
	d.Clear()
	return nil
}

func (d *Database) Clear() {
	d.dbTransaction = nil
	d.transactionOpened = false
}

func NewCoreDBManager(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}
