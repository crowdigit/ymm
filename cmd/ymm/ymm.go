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

var rootCmd = &cobra.Command{
	Use:   "ymm",
	Short: "Download MP3 files from Youtube and manage files",
}

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "Download a MP3 file from URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ci := exec.NewCommandProvider()

		var cmdYtMetadata internal.ExecConfig
		if err := viper.UnmarshalKey("command.metadata.youtube", &cmdYtMetadata); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		var cmdJq internal.ExecConfig
		if err := viper.UnmarshalKey("command.metadata.json", &cmdJq); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		var cmdYtDownload internal.ExecConfig
		if err := viper.UnmarshalKey("command.download.youtube", &cmdYtDownload); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return internal.DownloadSingle(ci, cmdYtMetadata, cmdJq, cmdYtDownload, args[0])
	},
}

func init() {
	configPath, err := xdg.ConfigFile(filepath.Join("ymm", "config.toml"))
	if err != nil {
		panic(err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Dir(configPath))

	viper.SetDefault("command.metadata.youtube.path", "yt-dlp")
	viper.SetDefault("command.metadata.youtube.args", []string{"--dump-json", "<url>"})
	viper.SetDefault("command.metadata.json.path", "jaq")
	viper.SetDefault("command.metadata.json.args", []string{"--slurp", "."})
	viper.SetDefault("command.download.youtube.path", "yt-dlp")
	viper.SetDefault(
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
		if err := viper.SafeWriteConfig(); err != nil {
			panic(fmt.Errorf("failed to write config file: %w", err))
		}
	}

	rootCmd.AddCommand(singleCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
