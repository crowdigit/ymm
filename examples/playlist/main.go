package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/crowdigit/ymm/app"
	"github.com/crowdigit/ymm/command"
	"github.com/crowdigit/ymm/db"
	"github.com/crowdigit/ymm/jq"
	"github.com/crowdigit/ymm/loudness"
	"github.com/crowdigit/ymm/ydl"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
)

func main() {
	loggerP, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}
	defer loggerP.Sync()
	logger := loggerP.Sugar()

	url := "https://www.youtube.com/playlist?list=PLcRmrvAR4VMI059jnMfqB4-6eicKuSt66"

	commandProvider := command.NewCommandProviderImpl()
	youtubeDl := ydl.NewYoutubeDLImpl(logger, commandProvider)
	loudnessScanner := loudness.NewLoudnessScanner(logger, commandProvider)
	jq := jq.NewJq(logger, commandProvider)

	metadataDir := filepath.Join(xdg.DataHome, "ymm", "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		log.Fatalf("failed to make metadata directory: %s", err)
	}

	databaseFile, err := xdg.DataFile("ymm/db.sql")
	if err != nil {
		log.Fatalf("failed to create data directory: %s", err)
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, databaseFile)
	if err != nil {
		log.Fatalf("failed to open Sqlite DB: %s", err)
	}
	defer sqldb.Close()

	db, err := db.NewDatabaseImpl(db.DatabaseConfig{MetadataDir: metadataDir}, sqldb)
	if err != nil {
		log.Fatalf("failed to initialize DB: %s", err)
	}

	config := app.ApplicationConfig{
		DownloadRootDir: filepath.Join(xdg.DataHome, "ymm", "music"),
	}

	app := app.NewApplicationImpl(logger, youtubeDl, loudnessScanner, jq, db, config)
	logger.Info("initialized application")
	if err := app.DownloadPlaylist(url); err != nil {
		log.Fatalf("%s", err)
	}
}
