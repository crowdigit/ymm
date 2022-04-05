package db

import "github.com/uptrace/bun"

func NewInsertUploaderQuery(db *bun.DB, uploader Uploader) *bun.InsertQuery {
	return db.NewInsert().Model(&uploader)
}

func NewSelectUploaderQuery(db *bun.DB, id string) *bun.SelectQuery {
	return db.NewSelect().Where("id = ?", id)
}

func NewInsertDownloadQuery(db *bun.DB, download Download) *bun.InsertQuery {
	return db.NewInsert().Model(&download)
}

func NewSelectDownloadQuery(db *bun.DB, id string) *bun.SelectQuery {
	return db.NewSelect().Where("id = ?", id)
}
