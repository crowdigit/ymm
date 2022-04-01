package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/crowdigit/ymm/app"
	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/ydl"
	"go.uber.org/zap"
)

func main() {
	loggerP, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}
	defer loggerP.Sync()
	logger := loggerP.Sugar()

	url := "https://www.youtube.com/watch?v=Ss-ba-g82-0"

	commandProvider := ydl.NewCommandProviderImpl()
	youtubeDl := ydl.NewYoutubeDLImpl(logger, commandProvider)

	metadataDir := filepath.Join(xdg.DataHome, "ymm", "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		log.Fatalf("failed to make metadata directory: %s", err)
	}

	databaseFile, err := xdg.DataFile("ymm/db.sql")
	if err != nil {
		log.Fatalf("failed to create data directory: %s", err)
	}

	db, err := db.NewDatabaseImpl(db.DatabaseConfig{
		DatabaseFile: databaseFile,
		MetadataDir:  metadataDir,
	})
	if err != nil {
		log.Fatalf("failed to initialize DB: %s", err)
	}
	defer db.Close()

	app := app.NewApplicationImpl(logger, youtubeDl, db)
	logger.Info("initialized application")
	if err := app.DownloadSingle(url); err != nil {
		log.Fatalf("%s", err)
	}
}
