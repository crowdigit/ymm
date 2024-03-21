package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/crowdigit/exec"
	"github.com/crowdigit/ymm/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ymm struct {
	cp     exec.CommandProvider
	config internal.Config
}

type Command func(*cobra.Command, []string) error

func (a ymm) downloadSingle() *cobra.Command {
	return &cobra.Command{
		Use:   "single",
		Short: "Download a MP3 file from URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.DownloadSingle(a.cp, a.config, args[0])
		},
	}
}

var rootCmd = &cobra.Command{
	Use:   "ymm",
	Short: "Download MP3 files from Youtube and manage files",
}

func initConfig() (internal.Config, error) {
	v := viper.New()

	configPath, err := xdg.ConfigFile(filepath.Join("ymm", "config.toml"))
	if err != nil {
		return internal.Config{}, fmt.Errorf("failed to get config file location: %w", err)
	}

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(filepath.Dir(configPath))

	v.SetDefault("command.metadata.youtube.path", "yt-dlp")
	v.SetDefault("command.metadata.youtube.args", []string{"--dump-json", "<url>"})
	v.SetDefault("command.metadata.json.path", "jaq")
	v.SetDefault("command.metadata.json.args", []string{"--slurp", "."})
	v.SetDefault("command.download.youtube.path", "yt-dlp")
	v.SetDefault(
		"command.download.youtube.args",
		[]string{
			// "--cookies",
			// "<cookies>",
			"--format",
			"<format>",
			"--extract-audio",
			"--audio-format",
			"mp3",
			"--audio-quality",
			"0",
			"<url>",
		},
	)
	v.SetDefault("command.replaygain.path", "loudness")
	v.SetDefault("command.replaygain.args", []string{"tag", "--track", "<path>"})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return internal.Config{}, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := v.SafeWriteConfig(); err != nil {
			return internal.Config{}, fmt.Errorf("failed to write config file: %w", err)
		}
	}

	config := internal.Config{}
	if err := v.Unmarshal(&config); err != nil {
		return internal.Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

func main() {
	config, err := initConfig()
	if err != nil {
		panic(err)
	}

	ymm := ymm{
		cp:     exec.NewCommandProvider(),
		config: config,
	}

	rootCmd.AddCommand(ymm.downloadSingle())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
