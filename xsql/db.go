package xsql

import (
	"context"
	"database/sql"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ExecQuerier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
type DBI interface {
	ExecQuerier
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Config struct {
	DSN          string        // write data source name.
	ReadDSN      []string      // read data source name.
	Active       int           // pool
	Idle         int           // pool
	IdleTimeout  time.Duration // connect max life time.
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
}

func NewMySQL(c *Config) (*DB, error) {
	m, err := connect(c.DSN, c.Active, c.Idle, c.IdleTimeout)
	if err != nil {
		return nil, err
	}
	var rs []*sql.DB
	if len(c.ReadDSN) == 0 {
		rs = append(rs, m)
	}
	for _, v := range c.ReadDSN {
		r, err := connect(v, c.Active, c.Idle, c.IdleTimeout)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	db := &DB{
		master: m,
		slaves: rs,
		config: c,
	}

	return db, nil
}

func connect(dsn string, active, idle int, idleTimeout time.Duration) (*sql.DB, error) {
	m, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	m.SetMaxOpenConns(active)
	m.SetMaxIdleConns(idle)
	m.SetConnMaxIdleTime(idleTimeout)
	return m, nil
}

type DB struct {
	master *sql.DB
	slaves []*sql.DB
	idx    int64
	config *Config
}

func (db *DB) Master() *sql.DB {
	return db.master
}

func (db *DB) slave() *sql.DB {
	v := atomic.AddInt64(&db.idx, 1)
	return db.slaves[int(v)%len(db.slaves)]
}

func (db *DB) PingContext(ctx context.Context) error {
	if err := db.master.PingContext(ctx); err != nil {
		return err
	}
	for _, v := range db.slaves {
		if err := v.PingContext(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.master.ExecContext(ctx, query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.slave().QueryContext(ctx, query, args...)
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.master.BeginTx(ctx, opts)
}
