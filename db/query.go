package db

import "github.com/uptrace/bun"

func NewInsertUploaderQuery(db *bun.DB, uploader Uploader) *bun.InsertQuery {
	return db.NewInsert().Model(&uploader)
}

func NewGetUploaderQuery(db *bun.DB, id string) *bun.SelectQuery {
	return db.NewSelect().Where("id = ?", id)
}
