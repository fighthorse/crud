package crud

import (
	"context"
	"database/sql"

	"github.com/fighthorse/crud/crud/user"
	"github.com/fighthorse/crud/xsql"
)

type Client struct {
	config *xsql.Config
	db     *xsql.DB
	Master *ClientM
	User   *UserClient
}

type ClientM struct {
	User *UserClient
}

func (c *Client) init() {
	c.User = &UserClient{eq: c.db, config: c.config}
	c.Master = &ClientM{
		User: &UserClient{eq: c.db.Master(), config: c.config},
	}
}

type Tx struct {
	config *xsql.Config
	tx     *sql.Tx
	User   *UserClient
}

func (tx *Tx) init() {
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
