package database

import (
	"fmt"

	junodb "github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/database/postgresql"
	juno "github.com/forbole/juno/v4/types"

	"github.com/hjcore/gjuno/x/authz"
)

type Database interface {
	junodb.Database

	authz.Database
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ Database = &Db{}
)

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store Desmos-related data.
type Db struct {
	*postgresql.Database
}

// Builder allows to create a new Db instance implementing the database.Builder type
func Builder(ctx *junodb.Context) (junodb.Database, error) {
	database, err := postgresql.Builder(ctx)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid database type")
	}

	return &Db{
		Database: psqlDb,
	}, nil
}

// Cast casts the given database to be a *Db
func Cast(database junodb.Database) Database {
	gotabitDb, ok := (database).(Database)
	if !ok {
		panic(fmt.Errorf("database is not a GJuno database instance"))
	}
	return gotabitDb
}

// --------------------------------------------------------------------------------------------------------------------

// SaveTx overrides postgresql.Database to perform a no-op
func (db *Db) SaveTx(_ *juno.Tx) error {
	return nil
}

// HasValidator overrides postgresql.Database to perform a no-op
func (db *Db) HasValidator(_ string) (bool, error) {
	return true, nil
}

// SaveValidators overrides postgresql.Database to perform a no-op
func (db *Db) SaveValidators(_ []*juno.Validator) error {
	return nil
}

// SaveCommitSignatures overrides postgresql.Database to perform a no-op
func (db *Db) SaveCommitSignatures(_ []*juno.CommitSig) error {
	return nil
}
