// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Wrapper around sql.Conn.
type Connection struct {
	con   *sql.Conn
	ctx   context.Context
	hooks []ConnectionHook
}

// Runs all connection hooks.
// Returns all successful hooks.
func (c *Connection) enter() ([]ConnectionHook, error) {
	var hooks []ConnectionHook
	for _, hook := range c.hooks {
		if err := hook.Enter(c); err != nil {
			return hooks, err
		}
		hooks = append(hooks, hook)
	}
	return hooks, nil
}

func (c *Connection) Exec(query string, args ...any) (sql.Result, error) {
	return c.con.ExecContext(c.ctx, query, args...)
}

func (c *Connection) Query(query string, args ...any) (*sql.Rows, error) {
	return c.con.QueryContext(c.ctx, query, args...)
}

func (c *Connection) QueryRow(query string, args ...any) *sql.Row {
	return c.con.QueryRowContext(c.ctx, query, args...)
}

func (c *Connection) Begin() (*sql.Tx, error) {
	return c.con.BeginTx(c.ctx, nil)
}

func (c *Connection) Close() error {
	for i := len(c.hooks) - 1; i >= 0; i-- {
		if err := c.hooks[i].Exit(c); err != nil {
			return fmt.Errorf("could not run exit hooks: %w", err)
		}
	}
	return c.con.Close()
}

// The caller is expected to close the Connection after use.
func NewConnection(
	db *sql.DB,
	ctx context.Context,
	hooks ...ConnectionHook,
) (*Connection, error) {
	con, err := db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get database connection: %w", err)
	}

	c := Connection{
		con:   con,
		ctx:   ctx,
		hooks: hooks,
	}

	success, err := c.enter()
	if err != nil {
		defer con.Close()
		for _, hook := range success {
			defer func(h ConnectionHook) {
				_ = h.Exit(&c)
			}(hook)
		}
		return nil, fmt.Errorf("could not run connection hooks: %w", err)
	}
	return &c, nil
}

type ConnectionHook struct {
	// Called after the connection is created.
	Enter func(con *Connection) error

	// Called before the connection is closed.
	Exit func(con *Connection) error
}

// Does nothing.
func noop(_ *Connection) error {
	return nil
}

func DefaultConnectionHook() ConnectionHook {
	return ConnectionHook{
		Enter: noop,
		Exit:  noop,
	}
}
