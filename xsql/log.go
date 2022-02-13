package xsql

import (
	"context"
	"database/sql"
	"log"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

func Debug(db DBI) *DebugDB {
	d := &DebugDB{dbt: db, log: log.Default()}
	return d
}

type DebugDB struct {
	log Logger
	dbt DBI
}

type DebugTx struct {
	log Logger
	eq  ExecQuerier
}

func (d *DebugDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.log.Printf("Debug Exec: %s args:%+v", query, args)
	return d.dbt.ExecContext(ctx, query, args...)

}
func (d *DebugDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.log.Printf("Debug Query: %s args:%+v", query, args)
	return d.dbt.QueryContext(ctx, query, args...)
}

func (d *DebugDB) Begin() (*DebugTx, error) {
	tx, err := d.dbt.Begin()
	return &DebugTx{eq: tx, log: d.log}, err
}

func (d *DebugDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*DebugTx, error) {
	tx, err := d.dbt.BeginTx(ctx, opts)
	return &DebugTx{eq: tx, log: d.log}, err
}

func (d *DebugTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.log.Printf("Debug Tx Exec: %s args:%+v", query, args)
	return d.eq.ExecContext(ctx, query, args...)

}
func (d *DebugTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.log.Printf("Debug TX Query: %s args:%+v", query, args)
	return d.eq.QueryContext(ctx, query, args...)
}
