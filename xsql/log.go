package xsql

import (
	"context"
	"database/sql"
	"log"
)

func Debug(db *sql.DB) *DebugDB {
	return &DebugDB{db}
}

type DebugDB struct {
	*sql.DB
}

type DebugTx struct {
	*sql.Tx
}

func (d *DebugDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Debug Exec: %s args:%+v", query, args)
	return d.DB.ExecContext(ctx, query, args...)

}
func (d *DebugDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Debug Query: %s args:%+v", query, args)
	return d.DB.QueryContext(ctx, query, args...)
}

func (d *DebugDB) Begin() (*DebugTx, error) {
	tx, err := d.DB.Begin()
	return &DebugTx{tx}, err
}

func (d *DebugTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Debug Tx Exec: %s args:%+v", query, args)
	return d.Tx.ExecContext(ctx, query, args...)

}
func (d *DebugTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Debug TX Query: %s args:%+v", query, args)
	return d.Tx.QueryContext(ctx, query, args...)
}
