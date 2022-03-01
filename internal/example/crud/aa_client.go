package crud

import (
	"context"
	"database/sql"

	"example/crud/alltypetable"
	"example/crud/user"

	"github.com/hongshengjie/crud/xsql"
)

type Client struct {
	config       *xsql.Config
	db           *xsql.DB
	Master       *ClientM
	AllTypeTable *AllTypeTableClient
	User         *UserClient
}

type ClientM struct {
	AllTypeTable *AllTypeTableClient
	User         *UserClient
}

func (c *Client) init() {
	c.AllTypeTable = &AllTypeTableClient{eq: c.db, config: c.config}
	c.User = &UserClient{eq: c.db, config: c.config}
	c.Master = &ClientM{
		AllTypeTable: &AllTypeTableClient{eq: c.db.Master(), config: c.config},
		User:         &UserClient{eq: c.db.Master(), config: c.config},
	}
}

type Tx struct {
	config       *xsql.Config
	tx           *sql.Tx
	AllTypeTable *AllTypeTableClient
	User         *UserClient
}

func (tx *Tx) init() {
	tx.AllTypeTable = &AllTypeTableClient{eq: tx.tx, config: tx.config}
	tx.User = &UserClient{eq: tx.tx, config: tx.config}
}

func NewClient(config *xsql.Config) (*Client, error) {
	db, err := xsql.NewMySQL(config)
	if err != nil {
		return nil, err
	}
	c := &Client{config: config, db: db}
	c.init()
	return c, nil
}

func (c *Client) Begin(ctx context.Context) (*Tx, error) {
	return c.BeginTx(ctx, nil)
}

func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := c.db.Master().BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	t := &Tx{tx: tx, config: c.config}
	t.init()
	return t, nil
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

type AllTypeTableClient struct {
	eq     xsql.ExecQuerier
	config *xsql.Config
}

func (c *AllTypeTableClient) Find() *alltypetable.SelectBuilder {
	return alltypetable.Find(c.eq).Timeout(c.config.QueryTimeout)
}

func (c *AllTypeTableClient) Create() *alltypetable.InsertBuilder {
	return alltypetable.Create(c.eq).Timeout(c.config.ExecTimeout)
}

func (c *AllTypeTableClient) Update() *alltypetable.UpdateBuilder {
	return alltypetable.Update(c.eq).Timeout(c.config.ExecTimeout)
}

func (c *AllTypeTableClient) Delete() *alltypetable.DeleteBuilder {
	return alltypetable.Delete(c.eq).Timeout(c.config.ExecTimeout)
}

type UserClient struct {
	eq     xsql.ExecQuerier
	config *xsql.Config
}

func (c *UserClient) Find() *user.SelectBuilder {
	return user.Find(c.eq).Timeout(c.config.QueryTimeout)
}

func (c *UserClient) Create() *user.InsertBuilder {
	return user.Create(c.eq).Timeout(c.config.ExecTimeout)
}

func (c *UserClient) Update() *user.UpdateBuilder {
	return user.Update(c.eq).Timeout(c.config.ExecTimeout)
}

func (c *UserClient) Delete() *user.DeleteBuilder {
	return user.Delete(c.eq).Timeout(c.config.ExecTimeout)
}
