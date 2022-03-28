package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/crowdigit/ymm/app"
	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/ydl"
)

func main() {
	url := "https://www.youtube.com/watch?v=Ss-ba-g82-0"

	commandProvider := ydl.NewCommandProviderImpl()
	youtubeDl := ydl.NewYoutubeDLImpl(commandProvider)

	metadataDir := filepath.Join(xdg.DataHome, "ymm", "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		log.Fatalf("failed to make metadata directory: %s", err)
	}

	databaseFile, err := xdg.DataFile("ymm/db.sql")
	if err != nil {
		log.Fatalf("failed to create data directory: %s", err)
	}

	db := db.NewDatabaseImpl(db.DatabaseConfig{
		DatabaseFile: databaseFile,
		MetadataDir:  metadataDir,
	})

	app := app.NewApplicationImpl(youtubeDl, db)
	if err := app.DownloadSingle(url); err != nil {
		log.Fatalf("%s", err)
	}
}
