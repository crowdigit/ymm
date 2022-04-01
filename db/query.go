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

func NewInsertUploaderQuery(uploader Uploader) QueryBuilder {
	panic("not implemented")
}
