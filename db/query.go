package db

import "github.com/uptrace/bun"

type QueryBuilder interface {
	Query() bun.Query
}

type InsertUploaderQueryBuilder struct {
}

func (b InsertUploaderQueryBuilder) Query() bun.Query {
	panic("not implemented")
}

func NewInsertUploaderQuery(db *bun.DB, uploader Uploader) QueryBuilder {
	panic("not implemented")
}

func NewGetUploaderQuery(db *bun.DB, id string) QueryBuilder {
	panic("not implemented")
}
