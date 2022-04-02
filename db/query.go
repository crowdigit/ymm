package db

import "github.com/uptrace/bun"

func NewInsertUploaderQuery(db *bun.DB, uploader Uploader) bun.Query {
	return nil
}

func NewGetUploaderQuery(db *bun.DB, id string) bun.Query {
	return db.NewSelect().TableExpr("uploaders").Where("id = ?", id)
}
