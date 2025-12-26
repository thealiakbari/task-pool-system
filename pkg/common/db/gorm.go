package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrTransactionStarted = errors.New("transaction has already been started")
	ErrDBType             = errors.New("wrong type of DB interface")
	ErrDBTimeOut          = errors.New("can't set transaction timeout")
)

const (
	ApmDBConnsStat           = "db.connections_status"
	ApmQueryIndexUsagePrefix = "db.query_index_usage."
)

type txKey struct{}

type DB interface {
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Last(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB
	FirstOrInit(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Update(column string, value interface{}) (tx *gorm.DB)
	Updates(values interface{}) (tx *gorm.DB)
	UpdateColumn(column string, value interface{}) (tx *gorm.DB)
	UpdateColumns(values interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Count(count *int64) (tx *gorm.DB)
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) (tx *gorm.DB)
	Pluck(column string, dest interface{}) (tx *gorm.DB)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Connection(fc func(tx *gorm.DB) error) (err error)
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	SavePoint(name string) *gorm.DB
	RollbackTo(name string) *gorm.DB
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	WithContext(ctx context.Context) *gorm.DB
	Model(value interface{}) (tx *gorm.DB)
	Table(name string, args ...interface{}) (tx *gorm.DB)
	Unscoped() (tx *gorm.DB)
}

func GormConnection(ctx context.Context, db *gorm.DB) DB {
	tx, ok := gormTxFromContext(ctx)
	if ok {
		return tx
	}
	return db.WithContext(ctx)
}

func gormTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

func BeginTx(ctx context.Context, db DB) (*gorm.DB, context.Context, error) {
	tx, ok := gormTxFromContext(ctx)
	if ok {
		return tx, ctx, ErrTransactionStarted
	}

	dbt, ok := db.(*gorm.DB)
	if !ok {
		return nil, ctx, ErrDBType
	}
	ctx, _ = context.WithTimeout(ctx, transactionTimeOut)

	tx = dbt.Begin().WithContext(ctx)
	go func() {
		<-ctx.Done()
		tx.Rollback()
		fmt.Println("The context has been canceled and transaction timeout")
	}()

	return tx, withTx(ctx, tx), nil
}

func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
