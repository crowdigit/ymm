package db

import "github.com/uptrace/bun"

func NewInsertUploaderQuery(db *bun.DB, uploader Uploader) *bun.InsertQuery {
	return nil
}

func NewGetUploaderQuery(db *bun.DB, id string) *bun.SelectQuery {
	return db.NewSelect().TableExpr("uploaders").Where("id = ?", id)
}
